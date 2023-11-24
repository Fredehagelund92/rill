package queries_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	// "github.com/golang/protobuf/proto"
	// garuntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	// "github.com/grpc-ecosystem/grpc-gateway/utilities"
	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/queries"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"

	_ "github.com/rilldata/rill/runtime/drivers/duckdb"
)

func TestMetricsViewsAggregation(t *testing.T) {
	rt, instanceID := testruntime.NewInstanceForProject(t, "ad_bids")

	ctrl, err := rt.Controller(context.Background(), instanceID)
	require.NoError(t, err)
	r, err := ctrl.Get(context.Background(), &runtimev1.ResourceName{Kind: runtime.ResourceKindMetricsView, Name: "ad_bids_metrics"}, false)
	require.NoError(t, err)
	mv := r.GetMetricsView().Spec

	limit := int64(10)
	q := &queries.MetricsViewAggregation{
		MetricsViewName: "ad_bids_metrics",
		Dimensions: []*runtimev1.MetricsViewAggregationDimension{
			{
				Name: "pub",
			},

			{
				Name:      "timestamp",
				TimeGrain: runtimev1.TimeGrain_TIME_GRAIN_MONTH,
			},
		},
		Measures: []*runtimev1.MetricsViewAggregationMeasure{
			{
				Name: "measure_1",
			},
		},
		MetricsView: mv,
		Sort: []*runtimev1.MetricsViewAggregationSort{
			{
				Name: "pub",
			},
			{
				Name: "timestamp",
			},
		},

		Limit: &limit,
	}
	err = q.Resolve(context.Background(), rt, instanceID, 0)
	require.NoError(t, err)
	require.NotEmpty(t, q.Result)
	for i, row := range q.Result.Data {
		for _, f := range row.Fields {
			fmt.Printf("%v ", f.AsInterface())
		}
		fmt.Printf(" %d \n", i)

	}
	rows := q.Result.Data

	i := 0
	require.Equal(t, "Facebook,2022-01-01", fieldsToString(rows[i], "pub", "timestamp"))
	i++
	require.Equal(t, "Facebook,2022-02-01", fieldsToString(rows[i], "pub", "timestamp"))
	i++
	require.Equal(t, "Facebook,2022-03-01", fieldsToString(rows[i], "pub", "timestamp"))
	i++
	require.Equal(t, "Google,2022-01-01", fieldsToString(rows[i], "pub", "timestamp"))
	i++
	require.Equal(t, "Google,2022-02-01", fieldsToString(rows[i], "pub", "timestamp"))
	i++
	require.Equal(t, "Google,2022-03-01", fieldsToString(rows[i], "pub", "timestamp"))
	i++
	require.Equal(t, "Microsoft,2022-01-01", fieldsToString(rows[i], "pub", "timestamp"))
	i++
	require.Equal(t, "Microsoft,2022-02-01", fieldsToString(rows[i], "pub", "timestamp"))
	i++
	require.Equal(t, "Microsoft,2022-03-01", fieldsToString(rows[i], "pub", "timestamp"))
	i++
	require.Equal(t, "Yahoo,2022-01-01", fieldsToString(rows[i], "pub", "timestamp"))
}

func TestMetricsViewsAggregation_pivot(t *testing.T) {
	rt, instanceID := testruntime.NewInstanceForProject(t, "ad_bids")

	ctrl, err := rt.Controller(context.Background(), instanceID)
	require.NoError(t, err)
	r, err := ctrl.Get(context.Background(), &runtimev1.ResourceName{Kind: runtime.ResourceKindMetricsView, Name: "ad_bids_metrics"}, false)
	require.NoError(t, err)
	mv := r.GetMetricsView().Spec

	limit := int64(10)
	q := &queries.MetricsViewAggregation{
		MetricsViewName: "ad_bids_metrics",
		Dimensions: []*runtimev1.MetricsViewAggregationDimension{
			{
				Name: "pub",
			},

			{
				Name:      "timestamp",
				TimeGrain: runtimev1.TimeGrain_TIME_GRAIN_MONTH,
			},
		},
		Measures: []*runtimev1.MetricsViewAggregationMeasure{
			{
				Name: "measure_1",
			},
		},
		MetricsView: mv,
		Sort: []*runtimev1.MetricsViewAggregationSort{
			{
				Name: "pub",
			},
		},
		PivotOn: []string{
			"timestamp",
		},
		Limit: &limit,
	}
	err = q.Resolve(context.Background(), rt, instanceID, 0)
	require.NoError(t, err)
	require.NotEmpty(t, q.Result)
	for i, row := range q.Result.Data {
		for _, f := range row.Fields {
			fmt.Printf("%v ", f.AsInterface())
		}
		fmt.Printf(" %d \n", i)

	}
	rows := q.Result.Data

	require.Equal(t, 4, len(q.Result.Schema.Fields))
	require.Equal(t, "pub", q.Result.Schema.Fields[0].Name)
	require.Equal(t, "2022-01-01_measure_1", q.Result.Schema.Fields[1].Name)
	require.Equal(t, "2022-02-01_measure_1", q.Result.Schema.Fields[2].Name)
	require.Equal(t, "2022-03-01_measure_1", q.Result.Schema.Fields[3].Name)

	i := 0
	require.Equal(t, "Facebook", fieldsToString(rows[i], "pub"))
	i++
	require.Equal(t, "Google", fieldsToString(rows[i], "pub"))
	i++
	require.Equal(t, "Microsoft", fieldsToString(rows[i], "pub"))
	i++
	require.Equal(t, "Yahoo", fieldsToString(rows[i], "pub"))
	i++
	require.Equal(t, "", fieldsToString(rows[i], "pub"))
}

func TestMetricsViewsAggregation_pivot_2_measures(t *testing.T) {
	rt, instanceID := testruntime.NewInstanceForProject(t, "ad_bids")

	ctrl, err := rt.Controller(context.Background(), instanceID)
	require.NoError(t, err)
	r, err := ctrl.Get(context.Background(), &runtimev1.ResourceName{Kind: runtime.ResourceKindMetricsView, Name: "ad_bids_metrics"}, false)
	require.NoError(t, err)
	mv := r.GetMetricsView().Spec

	limit := int64(10)
	q := &queries.MetricsViewAggregation{
		MetricsViewName: "ad_bids_metrics",
		Dimensions: []*runtimev1.MetricsViewAggregationDimension{
			{
				Name: "pub",
			},

			{
				Name:      "timestamp",
				TimeGrain: runtimev1.TimeGrain_TIME_GRAIN_MONTH,
			},
		},
		Measures: []*runtimev1.MetricsViewAggregationMeasure{
			{
				Name: "measure_1",
			},
			{
				Name: "measure_0",
			},
		},
		MetricsView: mv,
		Sort: []*runtimev1.MetricsViewAggregationSort{
			{
				Name: "pub",
			},
		},
		PivotOn: []string{
			"timestamp",
		},
		Limit: &limit,
	}
	err = q.Resolve(context.Background(), rt, instanceID, 0)
	require.NoError(t, err)
	require.NotEmpty(t, q.Result)
	for i, row := range q.Result.Data {
		for _, f := range row.Fields {
			fmt.Printf("%v ", f.AsInterface())
		}
		fmt.Printf(" %d \n", i)

	}
	rows := q.Result.Data

	require.Equal(t, q.Result.Schema.Fields[0].Name, "pub")
	require.Equal(t, q.Result.Schema.Fields[1].Name, "2022-01-01_measure_1")
	require.Equal(t, q.Result.Schema.Fields[2].Name, "2022-01-01_measure_0")

	require.Equal(t, q.Result.Schema.Fields[3].Name, "2022-02-01_measure_1")
	require.Equal(t, q.Result.Schema.Fields[4].Name, "2022-02-01_measure_0")

	require.Equal(t, q.Result.Schema.Fields[5].Name, "2022-03-01_measure_1")
	require.Equal(t, q.Result.Schema.Fields[6].Name, "2022-03-01_measure_0")

	i := 0
	require.Equal(t, "Facebook", fieldsToString(rows[i], "pub"))
	i++
	require.Equal(t, "Google", fieldsToString(rows[i], "pub"))
	i++
	require.Equal(t, "Microsoft", fieldsToString(rows[i], "pub"))
	i++
	require.Equal(t, "Yahoo", fieldsToString(rows[i], "pub"))
	i++
	require.Equal(t, "", fieldsToString(rows[i], "pub"))
}

func TestMetricsViewsAggregation_pivot_2_measures_and_filter(t *testing.T) {
	rt, instanceID := testruntime.NewInstanceForProject(t, "ad_bids")

	ctrl, err := rt.Controller(context.Background(), instanceID)
	require.NoError(t, err)
	r, err := ctrl.Get(context.Background(), &runtimev1.ResourceName{Kind: runtime.ResourceKindMetricsView, Name: "ad_bids_metrics"}, false)
	require.NoError(t, err)
	mv := r.GetMetricsView().Spec

	limit := int64(10)
	q := &queries.MetricsViewAggregation{
		MetricsViewName: "ad_bids_metrics",
		Dimensions: []*runtimev1.MetricsViewAggregationDimension{
			{
				Name: "pub",
			},

			{
				Name:      "timestamp",
				TimeGrain: runtimev1.TimeGrain_TIME_GRAIN_MONTH,
			},
		},
		Measures: []*runtimev1.MetricsViewAggregationMeasure{
			{
				Name: "measure_1",
			},
			{
				Name: "measure_0",
			},
		},
		MetricsView: mv,
		Sort: []*runtimev1.MetricsViewAggregationSort{
			{
				Name: "pub",
			},
		},
		PivotOn: []string{
			"timestamp",
		},
		Filter: &runtimev1.MetricsViewFilter{
			Include: []*runtimev1.MetricsViewFilter_Cond{
				{
					Name: "pub",
					In:   []*structpb.Value{structpb.NewStringValue("Google")},
				},
			},
		},
		Limit: &limit,
	}
	err = q.Resolve(context.Background(), rt, instanceID, 0)
	require.NoError(t, err)
	require.NotEmpty(t, q.Result)
	for i, row := range q.Result.Data {
		for _, f := range row.Fields {
			fmt.Printf("%v ", f.AsInterface())
		}
		fmt.Printf(" %d \n", i)

	}
	rows := q.Result.Data

	require.Equal(t, q.Result.Schema.Fields[0].Name, "pub")
	require.Equal(t, q.Result.Schema.Fields[1].Name, "2022-01-01_measure_1")
	require.Equal(t, q.Result.Schema.Fields[2].Name, "2022-01-01_measure_0")

	require.Equal(t, q.Result.Schema.Fields[3].Name, "2022-02-01_measure_1")
	require.Equal(t, q.Result.Schema.Fields[4].Name, "2022-02-01_measure_0")

	require.Equal(t, q.Result.Schema.Fields[5].Name, "2022-03-01_measure_1")
	require.Equal(t, q.Result.Schema.Fields[6].Name, "2022-03-01_measure_0")

	require.Equal(t, 1, len(rows))
	i := 0
	require.Equal(t, "Google", fieldsToString(rows[i], "pub"))
}

func TestMetricsViewsAggregation_pivot_dim_and_measure(t *testing.T) {
	rt, instanceID := testruntime.NewInstanceForProject(t, "ad_bids")

	ctrl, err := rt.Controller(context.Background(), instanceID)
	require.NoError(t, err)
	r, err := ctrl.Get(context.Background(), &runtimev1.ResourceName{Kind: runtime.ResourceKindMetricsView, Name: "ad_bids_metrics"}, false)
	require.NoError(t, err)
	mv := r.GetMetricsView().Spec

	limit := int64(10)
	q := &queries.MetricsViewAggregation{
		MetricsViewName: "ad_bids_metrics",
		Dimensions: []*runtimev1.MetricsViewAggregationDimension{
			{
				Name: "pub",
			},
			{
				Name: "dom",
			},
			{
				Name:      "timestamp",
				TimeGrain: runtimev1.TimeGrain_TIME_GRAIN_MONTH,
			},
		},
		Measures: []*runtimev1.MetricsViewAggregationMeasure{
			{
				Name: "measure_1",
			},
		},
		Filter: &runtimev1.MetricsViewFilter{
			Include: []*runtimev1.MetricsViewFilter_Cond{
				{
					Name: "pub",
					In:   []*structpb.Value{structpb.NewStringValue("Google")},
				},
			},
		},
		MetricsView: mv,
		Sort: []*runtimev1.MetricsViewAggregationSort{
			{
				Name: "dom",
			},
		},
		PivotOn: []string{
			"timestamp",
			"pub",
		},
		Limit: &limit,
	}
	err = q.Resolve(context.Background(), rt, instanceID, 0)
	require.NoError(t, err)
	require.NotEmpty(t, q.Result)
	for _, s := range q.Result.Schema.Fields {
		fmt.Printf("%v ", s.Name)
	}
	for i, row := range q.Result.Data {
		for _, f := range row.Fields {
			fmt.Printf("%v ", f.AsInterface())
		}
		fmt.Printf(" %d \n", i)

	}
	rows := q.Result.Data

	require.Equal(t, q.Result.Schema.Fields[0].Name, "dom")
	require.Equal(t, q.Result.Schema.Fields[1].Name, "2022-01-01_Google_measure_1")
	require.Equal(t, q.Result.Schema.Fields[2].Name, "2022-02-01_Google_measure_1")
	require.Equal(t, q.Result.Schema.Fields[3].Name, "2022-03-01_Google_measure_1")

	i := 0
	require.Equal(t, "google.com", fieldsToString(rows[i], "dom"))
}

func Ignore_TestMetricsViewsAggregation_Druid(t *testing.T) {
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}

	conn, err := grpc.Dial(":49009", dialOpts...)
	if err != nil {
		require.NoError(t, err)
	}
	defer conn.Close()

	client := runtimev1.NewQueryServiceClient(conn)
	req := &runtimev1.MetricsViewAggregationRequest{
		InstanceId:  "default",
		MetricsView: "test_data_test",
		Dimensions: []*runtimev1.MetricsViewAggregationDimension{
			{
				Name: "publisher",
			},
			{
				Name:      "__time",
				TimeGrain: runtimev1.TimeGrain_TIME_GRAIN_MONTH,
			},
		},
		Measures: []*runtimev1.MetricsViewAggregationMeasure{
			{
				Name: "bp",
			},
		},
		Sort: []*runtimev1.MetricsViewAggregationSort{
			{
				Name: "publisher",
			},
			{
				Name: "__time",
			},
		},
	}

	resp, err := client.MetricsViewAggregation(context.Background(), req)
	if err != nil {
		require.NoError(t, err)
	}
	rows := resp.Data

	for _, s := range resp.Schema.Fields {
		fmt.Printf("%v ", s.Name)
	}
	fmt.Println()
	for i, row := range resp.Data {
		for _, s := range resp.Schema.Fields {
			fmt.Printf("%v ", row.Fields[s.Name].AsInterface())
		}
		fmt.Printf(" %d \n", i)

	}
	i := 0
	require.Equal(t, ",2022-01-01T00:00:00Z", fieldsToString(rows[i], "publisher", "__time"))
	i++
	require.Equal(t, ",2022-02-01T00:00:00Z", fieldsToString(rows[i], "publisher", "__time"))
	i++
	require.Equal(t, ",2022-03-01T00:00:00Z", fieldsToString(rows[i], "publisher", "__time"))
	i++
	require.Equal(t, "Facebook,2022-01-01T00:00:00Z", fieldsToString(rows[i], "publisher", "__time"))
	i++
	require.Equal(t, "Facebook,2022-02-01T00:00:00Z", fieldsToString(rows[i], "publisher", "__time"))
	i++
	require.Equal(t, "Facebook,2022-03-01T00:00:00Z", fieldsToString(rows[i], "publisher", "__time"))
	i++
	require.Equal(t, "Google,2022-01-01T00:00:00Z", fieldsToString(rows[i], "publisher", "__time"))
	i++
	require.Equal(t, "Google,2022-02-01T00:00:00Z", fieldsToString(rows[i], "publisher", "__time"))
	i++
	require.Equal(t, "Google,2022-03-01T00:00:00Z", fieldsToString(rows[i], "publisher", "__time"))
	i++
	require.Equal(t, "Microsoft,2022-01-01T00:00:00Z", fieldsToString(rows[i], "publisher", "__time"))
	i++
	require.Equal(t, "Microsoft,2022-02-01T00:00:00Z", fieldsToString(rows[i], "publisher", "__time"))
	i++
	require.Equal(t, "Microsoft,2022-03-01T00:00:00Z", fieldsToString(rows[i], "publisher", "__time"))
	i++
	require.Equal(t, "Yahoo,2022-01-01T00:00:00Z", fieldsToString(rows[i], "publisher", "__time"))
}

func Ignore_TestMetricsViewsAggregation_Druid_pivot(t *testing.T) {
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}

	conn, err := grpc.Dial(":49009", dialOpts...)
	if err != nil {
		require.NoError(t, err)
	}
	defer conn.Close()

	client := runtimev1.NewQueryServiceClient(conn)
	req := &runtimev1.MetricsViewAggregationRequest{
		InstanceId:  "default",
		MetricsView: "test_data_test",
		Dimensions: []*runtimev1.MetricsViewAggregationDimension{
			{
				Name: "publisher",
			},
			{
				Name:      "__time",
				TimeGrain: runtimev1.TimeGrain_TIME_GRAIN_MONTH,
			},
		},
		Measures: []*runtimev1.MetricsViewAggregationMeasure{
			{
				Name: "bp",
			},
		},
		Sort: []*runtimev1.MetricsViewAggregationSort{
			{
				Name: "publisher",
			},
			{
				Name: "__time",
			},
		},
		PivotOn: []string{
			"__time",
		},
	}

	resp, err := client.MetricsViewAggregation(context.Background(), req)
	if err != nil {
		require.NoError(t, err)
	}
	rows := resp.Data

	for _, s := range resp.Schema.Fields {
		fmt.Printf("%v ", s.Name)
	}
	fmt.Println()
	for i, row := range resp.Data {
		for _, s := range resp.Schema.Fields {
			fmt.Printf("%v ", row.Fields[s.Name].AsInterface())
		}
		fmt.Printf(" %d \n", i)

	}
	require.Equal(t, 4, len(resp.Schema.Fields))
	require.Equal(t, "publisher", resp.Schema.Fields[0].Name)
	require.Equal(t, "2022-01-01_bp", resp.Schema.Fields[1].Name)
	require.Equal(t, "2022-02-01_bp", resp.Schema.Fields[2].Name)
	require.Equal(t, "2022-03-01_bp", resp.Schema.Fields[3].Name)

	i := 0
	require.Equal(t, "", fieldsToString(rows[i], "publisher"))
	i++
	require.Equal(t, "Facebook", fieldsToString(rows[i], "publisher"))
	i++
	require.Equal(t, "Google", fieldsToString(rows[i], "publisher"))
	i++
	require.Equal(t, "Microsoft", fieldsToString(rows[i], "publisher"))
	i++
	require.Equal(t, "Yahoo", fieldsToString(rows[i], "publisher"))
}

func fieldsToString(row *structpb.Struct, args ...string) string {
	s := make([]string, 0, len(args))
	for _, arg := range args {
		s = append(s, row.Fields[arg].GetStringValue())
	}
	return strings.Join(s, ",")
}
