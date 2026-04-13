package cmd

import (
	"errors"
	"time"

	"github.com/avast/retry-go"
	"github.com/majd/ipatool/v2/internal/core"
	"github.com/spf13/cobra"
)

// nolint:wrapcheck
func getVersionMetadataCmd() *cobra.Command {
	var (
		appID             int64
		bundleID          string
		externalVersionID string
	)

	cmd := &cobra.Command{
		Use:   "get-version-metadata",
		Short: "Retrieves the metadata for a specific version of an app",
		RunE: func(cmd *cobra.Command, args []string) error {
			if appID == 0 && bundleID == "" {
				return errors.New("either the app ID or the bundle identifier must be specified")
			}

			var lastErr error
			var acc core.Account

			return retry.Do(func() error {
				infoResult, err := dependencies.Core.AccountInfo()
				if err != nil {
					return err
				}

				acc = infoResult.Account

				if errors.Is(lastErr, core.ErrPasswordTokenExpired) {
					loginResult, err := dependencies.Core.Login(core.LoginInput{Email: acc.Email, Password: acc.Password})
					if err != nil {
						return err
					}

					acc = loginResult.Account
				}

				app := core.App{ID: appID}
				if bundleID != "" {
					lookupResult, err := dependencies.Core.Lookup(core.LookupInput{Account: acc, BundleID: bundleID})
					if err != nil {
						return err
					}

					app = lookupResult.App
				}

				out, err := dependencies.Core.GetVersionMetadata(core.GetVersionMetadataInput{
					Account:   acc,
					App:       app,
					VersionID: externalVersionID,
				})
				if err != nil {
					return err
				}

				dependencies.Logger.Log().
					Str("externalVersionID", externalVersionID).
					Str("displayVersion", out.DisplayVersion).
					Time("releaseDate", out.ReleaseDate).
					Bool("success", true).
					Send()

				return nil
			},
				retry.LastErrorOnly(true),
				retry.DelayType(retry.FixedDelay),
				retry.Delay(time.Millisecond),
				retry.Attempts(2),
				retry.RetryIf(func(err error) bool {
					lastErr = err

					return errors.Is(err, core.ErrPasswordTokenExpired)
				}),
			)
		},
	}

	cmd.Flags().Int64VarP(&appID, "app-id", "i", 0, "ID of the target iOS app (required)")
	cmd.Flags().StringVarP(&bundleID, "bundle-identifier", "b", "", "The bundle identifier of the target iOS app (overrides the app ID)")
	cmd.Flags().StringVar(&externalVersionID, "external-version-id", "", "External version identifier of the target iOS app (required)")

	_ = cmd.MarkFlagRequired("external-version-id")

	return cmd
}
