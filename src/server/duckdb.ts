// @ts-nocheck
import fs from "fs";
import duckdb from 'duckdb';
import { default as glob } from 'glob';

import { guidGenerator } from "../util/guid.js";

interface DB {
	all: Function;
	exec: Function;
	run: Function
}

export function connect() : DB {
	return new duckdb.Database(':memory:');
}

const db:DB = connect();

let onCallback;
let offCallback;

/** utilize these for setting the "running" and "not running" state in the frontend */
export function registerDBRunCallbacks(onCall:Function, offCall:Function) {
	onCallback = onCall;
	offCallback = offCall;
}

function dbAll(db:DB, query:string) {
	if (onCallback) {
		onCallback();
	}
	return new Promise((resolve, reject) => {
		try {
			db.all(query, (err, res) => {
				if (err !== null) {
					reject(err);
				} else {
					if (offCallback) offCallback();
					resolve(res);
				}
			});
		} catch (err) {
			reject(err);
		}
	});
};

export function dbRun(query:string) { 
	return new Promise((resolve, reject) => {
		db.run(query, (err) => {
				if (err !== null) reject(false);
				resolve(true);
			}
		)
	})
}

export async function validQuery(db:DB, query:string): Promise<{value: boolean, message?: string}> {
	return new Promise((resolve) => {
		db.run(query, (err) => {
			if (err !== null) {
				resolve({
					value: false,
					message: err.message
				});
			} else {
				resolve({ value: true});
			}
		});
	});
}

export function hasCreateStatement(query:string) {
	return query.toLowerCase().startsWith('create')
		? `Query has a CREATE statement. 
	Let us handle that for you!
	Just use SELECT and we'll do the rest.
	`
		: false;
}

export function containsMultipleQueries(query:string) {
	return query.split().filter((character) => character == ';').length > 1
		? 'Keep it to a single query please!'
		: false;
}

export function validateQuery(query:string, ...validators:Function[]) {
	return validators.map((validator) => validator(query)).filter((validation) => validation);
}

function wrapQueryAsTemporaryView(query:string) {
	return `CREATE OR REPLACE TEMPORARY VIEW tmp AS (
	${query.replace(';', '')}
);`;
}

export async function checkQuery(query:string) : Promise<void> {
	const output = {};
	const isValid = await validQuery(db, query);
	if (!(isValid.value)) {
		throw Error(isValid.message);
	}
	const validation = validateQuery(query, hasCreateStatement, containsMultipleQueries);
	if (validation.length) {
		throw Error(validation[0])
	}
}

export async function wrapQueryAsView(query:string) {
	return new Promise((resolve, reject) => {
		db.run(wrapQueryAsTemporaryView(query), (err) => {
			if (err !== null) reject(err);
			resolve(true);
		})
	})
}

export async function createPreview(query:string) {
    // FIXME: sort out the type here
	let preview:any;
    try {
		try {
			// get the preview.
			preview = await dbAll(db, 'SELECT * from tmp LIMIT 25;');
		} catch (err) {
			throw Error(err);
		}
	} catch (err) {
		throw Error(err)
	}
    return preview;
}

export async function createSourceProfile(parquetFile:string) {
	return await dbAll(db, `select * from parquet_schema('${parquetFile}');`) as any[];
}

export async function parquetToDBTypes(parquetFile:string) {
	const guid = guidGenerator().replace(/-/g, '_');
    await dbAll(db, `
	CREATE TEMP TABLE tbl_${guid} AS (
        SELECT * from '${parquetFile}' LIMIT 1
    );
	`);
	const tableDef = await dbAll(db, `PRAGMA table_info(tbl_${guid});`)
	await dbAll(db, `DROP TABLE tbl_${guid};`);
    return tableDef;
}

export async function getCardinality(parquetFile:string) {
	const [cardinality] =  await dbAll(db, `select count(*) as count FROM '${parquetFile}';`);
	return cardinality.count;
}

export async function getFirstN(table, n=1) {
	return  dbAll(db, `SELECT * from ${table} LIMIT ${n};`);
}

export function extractParquetFilesFromQuery(query:string) {
	let re = /'[^']*\.parquet'/g;
	const matches = query.match(re).map(match => match.replace(/'/g, ''));
	return matches;
}

export async function createSourceProfileFromQuery(query:string) {
	// capture output from parquet query.
	const matches = extractParquetFilesFromQuery(query);
	const tables = (matches === null) ? [] : await Promise.all(matches.map(async (strippedMatch) => {
		//let strippedMatch = match.replace(/'/g, '');
		let match = `'${strippedMatch}'`;
		const info = await createSourceProfile(strippedMatch);
		const head = await getFirstN(match);
		const cardinality = await getCardinality(strippedMatch);
		const sizeInBytes = await getDestinationSize(strippedMatch);
		return {
			profile: info.filter(i => i.name !== 'duckdb_schema'),
			head, 
			cardinality,
			table: strippedMatch,
			sizeInBytes,
			path: strippedMatch,
			name: strippedMatch.split('/').slice(-1)[0]
		}
	}))
	return tables;
}

export async function getDestinationSize(path:string) {
	if (fs.existsSync(path)) {
		const size = await dbAll(db, `SELECT total_compressed_size from parquet_metadata('${path}')`) as any[];
		return size.reduce((acc:number, v:object) => acc + v.total_compressed_size, 0)
	}
	return undefined;
}

export async function calculateDestinationCardinality(query:string) {
	const [outputSize] = await dbAll(db, 'SELECT count(*) AS cardinality from tmp;') as any[];
	return outputSize.cardinality;
}

export async function createDestinationProfile(query:string) {
	const info = await dbAll(db, `PRAGMA table_info(tmp);`);
	return info;
}

export async function exportToParquet(query:string, output:string) {
	// generate export just in case.
	if (!fs.existsSync('./export')) {
		fs.mkdirSync('./export');
	}
	const exportQuery = `COPY (${query.replace(';', '')}) TO '${output}' (FORMAT 'parquet')`;
	return dbRun(exportQuery);
}

export async function getParquetFilesInRoot() {
	return new Promise((resolve, reject) => {
		glob.glob('./**/*.parquet', {ignore: ['./node_modules/', './.svelte-kit/', './build/', './src/', './tsc-tmp']},
			(err, output) => {
				if (err!==null) reject(err);
				resolve(output);
			}
		)
	});
}
/**
 * getSummary
 * number: five number summary + mean
 * date: max, min, total time between the two
 * categorical: cardinality
 */

//  export function toDistributionSummary(column:string) {
// 	return [
// 		`min(${column}) as min_${column}`,
// 		`approx_quantile(${column}, 0.25) as q25_${column}`,
// 		`approx_quantile(${column}, 0.5)  as q50_${column}`,
// 		`approx_quantile(${column}, 0.75) as q75_${column}`,
// 		`max(${column}) as max_${column}`,
// 		`avg(${column}) as mean_${column}`,
// 		`stddev_pop(${column}) as sd_${column}`,
// 	]
// }

// // FIXME: deprecate and remove all code paths
// export async function getDistributionSummary(parquetFilePath:string, column:string) {
// 	const [point] = await dbAll(db, `
// SELECT 
// 	min(${column}) as min, 
// 	approx_quantile(${column}, 0.25) as q25, 
// 	approx_quantile(${column}, 0.5)  as q50,
// 	approx_quantile(${column}, 0.75) as q75,
// 	max(${column}) as max,
// 	avg(${column}) as mean,
// 	stddev_pop(${column}) as sd
// 	FROM '${parquetFilePath}';`);
// 	return point;
// }










export function toDistributionSummary(column) {
	return [
		`min(${column}) as ${column}_min`,
		`reservoir_quantile(${column}, 0.25) as ${column}_q25`,
		`reservoir_quantile(${column}, 0.5)  as ${column}_q50`,
		`reservoir_quantile(${column}, 0.75) as ${column}_q75`,
		`max(${column}) as ${column}_max`,
		`avg(${column}) as ${column}_mean`,
		`stddev_pop(${column}) as ${column}_sd`,
	]
}

function topK(parquetFile, column) {
	return `SELECT ${column} as value, count(*) AS count from '${parquetFile}'
GROUP BY ${column}
ORDER BY count desc
LIMIT 50;`
}

async function getTopKAndCardinality(parquetFilePath, column) {
	const topKValues = await dbAll(db, topK(parquetFilePath, column));
	const [cardinality] = await dbAll(db, `SELECT approx_count_distinct(${column}) as count from '${parquetFilePath}';`);
	return {
		column,
		topK: topKValues,
		cardinality: cardinality.count
	}
}

export async function getDistributionSummaries(parquetFilePath, fields) {
	const numericSelects = fields.map(n => toDistributionSummary(n.name).join(',\n  ')).join(',\n  ');
	const [summaries] = await dbAll(db, `SELECT \n${numericSelects}\n FROM '${parquetFilePath}';`);
	return fields.map(n => n.name).reduce((acc, field) => {
		acc[field] = {
			min: summaries[`${field}_min`],
			q25: summaries[`${field}_q25`],
			q50: summaries[`${field}_q50`],
			q75: summaries[`${field}_q75`],
			max: summaries[`${field}_max`],
			mean: summaries[`${field}_mean`],
			sd: summaries[`${field}_sd`],
		};
		return acc;
	}, {});
}

export async function getCategoricalSummaries(parquetFilePath, fields) {
	const summaries = await Promise.all(fields.map((s) => {
		return getTopKAndCardinality(parquetFilePath, s.name);
	}));
	return summaries.reduce((acc, fieldSummary) => {
		acc[fieldSummary.column] = {
			topK: fieldSummary.topK,
			cardinality: fieldSummary.cardinality
		}
		return acc;
	}, {});
}

export async function getTimestampSummaries(parquetFilePath:string, fields:any) {
	const queries = fields.map(field => {
		return `
		max(${field.name}) - min(${field.name}) AS ${field.name}_interval,
		min(${field.name}) as ${field.name}_min,
		max(${field.name}) as ${field.name}_max
		`
	}).join(',\n  ');
	const [results] = await dbAll(db, `
		SELECT
			${queries}
		FROM '${parquetFilePath}';
	`);
	return fields.reduce((acc, field) => {
		acc[field.name] = { 
			min: results[`${field.name}_min`],
			max: results[`${field.name}_max`],
			interval: results[`${field.name}_interval`]
		 };
		return acc;
	}, {});
}

export async function getNullCounts(parquetFilePath:string, fields:any) {
	const [nullities] = await dbAll(db, `
		SELECT
		${fields.map(field => {
			return `COUNT(CASE WHEN ${field.name} IS NULL THEN 1 ELSE NULL END) as ${field.name}`
		}).join(',\n')}
		FROM '${parquetFilePath}';
	`);
	return nullities;
}