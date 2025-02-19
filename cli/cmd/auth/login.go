package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/rilldata/rill/cli/pkg/browser"
	"github.com/rilldata/rill/cli/pkg/cmdutil"
	"github.com/rilldata/rill/cli/pkg/deviceauth"
	"github.com/rilldata/rill/cli/pkg/dotrill"
	adminv1 "github.com/rilldata/rill/proto/gen/rill/admin/v1"
	"github.com/spf13/cobra"
)

// LoginCmd is the command for logging into a Rill account.
func LoginCmd(ch *cmdutil.Helper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with the Rill API",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ch.Config
			ctx := cmd.Context()

			// updating this as its not required to logout first and login again
			if cfg.AdminTokenDefault != "" {
				err := Logout(ctx, ch)
				if err != nil {
					return err
				}
			}

			// Login user
			err := Login(ctx, ch, "")
			if err != nil {
				return err
			}

			// Set default org after login
			err = SelectOrgFlow(ctx, ch)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func Login(ctx context.Context, ch *cmdutil.Helper, redirectURL string) error {
	// In production, the REST and gRPC endpoints are the same, but in development, they're served on different ports.
	// We plan to move to connect.build for gRPC, which will allow us to serve both on the same port in development as well.
	// Until we make that change, this is a convenient hack for local development (assumes gRPC on port 9090 and REST on port 8080).
	cfg := ch.Config
	authURL := cfg.AdminURL
	if strings.Contains(authURL, "http://localhost:9090") {
		authURL = "http://localhost:8080"
	}

	authenticator, err := deviceauth.New(authURL)
	if err != nil {
		return err
	}

	deviceVerification, err := authenticator.VerifyDevice(ctx, redirectURL)
	if err != nil {
		return err
	}

	ch.Printer.PrintBold("\nConfirmation Code: ")
	ch.Printer.PrintlnSuccess(deviceVerification.UserCode)

	ch.Printer.PrintBold(fmt.Sprintf("\nOpen this URL in your browser to confirm the login: %s\n\n", deviceVerification.VerificationCompleteURL))

	_ = browser.Open(deviceVerification.VerificationCompleteURL)

	res1, err := authenticator.GetAccessTokenForDevice(ctx, deviceVerification)
	if err != nil {
		return err
	}

	err = dotrill.SetAccessToken(res1.AccessToken)
	if err != nil {
		return err
	}
	// set the default token to the one we just got
	cfg.AdminTokenDefault = res1.AccessToken
	ch.Printer.PrintBold("Successfully logged in. Welcome to Rill!\n")
	return nil
}

func SelectOrgFlow(ctx context.Context, ch *cmdutil.Helper) error {
	cfg := ch.Config
	client, err := cmdutil.Client(cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	res, err := client.ListOrganizations(context.Background(), &adminv1.ListOrganizationsRequest{})
	if err != nil {
		return err
	}

	if len(res.Organizations) == 0 {
		ch.Printer.PrintlnWarn("You are not part of an org. Run `rill org create` or `rill deploy` to create one.")
		return nil
	}

	var orgNames []string
	for _, org := range res.Organizations {
		orgNames = append(orgNames, org.Name)
	}

	defaultOrg := orgNames[0]
	if len(orgNames) > 1 {
		defaultOrg = cmdutil.SelectPrompt("Select default org (to change later, run `rill org switch`).", orgNames, defaultOrg)
	}

	err = dotrill.SetDefaultOrg(defaultOrg)
	if err != nil {
		return err
	}

	cfg.Org = defaultOrg

	ch.Printer.Print(fmt.Sprintf("Set default organization to %q. Change using `rill org switch`.\n", defaultOrg))
	return nil
}
