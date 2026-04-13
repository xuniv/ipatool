package cmd

import (
	"errors"
	"time"

	"github.com/avast/retry-go"
	"github.com/majd/ipatool/v2/internal/core"
	"github.com/spf13/cobra"
)

// nolint:wrapcheck
func purchaseCmd() *cobra.Command {
	var bundleID string

	cmd := &cobra.Command{
		Use:   "purchase",
		Short: "Obtain a license for the app from the App Store",
		RunE: func(cmd *cobra.Command, args []string) error {
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

				lookupResult, err := dependencies.Core.Lookup(core.LookupInput{Account: acc, BundleID: bundleID})
				if err != nil {
					return err
				}

				err = dependencies.Core.Purchase(core.PurchaseInput{Account: acc, App: lookupResult.App})
				if err != nil {
					return err
				}

				dependencies.Logger.Log().Bool("success", true).Send()

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

	cmd.Flags().StringVarP(&bundleID, "bundle-identifier", "b", "", "Bundle identifier of the target iOS app (required)")
	_ = cmd.MarkFlagRequired("bundle-identifier")

	return cmd
}
