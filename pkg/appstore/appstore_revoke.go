package appstore

import (
	"fmt"
)

func (t *appstore) Revoke() error {
	err := t.credentialStore.Remove("account")
	if err != nil {
		return fmt.Errorf("failed to remove account from keychain: %w", err)
	}

	return nil
}
