package superuser

import (
	"github.com/rilldata/rill/cli/pkg/cmdutil"
	adminv1 "github.com/rilldata/rill/proto/gen/rill/admin/v1"
	"github.com/spf13/cobra"
)

func ListCmd(ch *cmdutil.Helper) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Args:  cobra.NoArgs,
		Short: "List superusers",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cfg := ch.Config

			client, err := cmdutil.Client(cfg)
			if err != nil {
				return err
			}
			defer client.Close()

			res, err := client.ListSuperusers(ctx, &adminv1.ListSuperusersRequest{})
			if err != nil {
				return err
			}

			err = cmdutil.PrintUsers(ch.Printer, res.Users)
			if err != nil {
				return err
			}

			return nil
		},
	}
	return listCmd
}
