package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/99designs/keyring"
	cookiejar "github.com/juju/persistent-cookiejar"
	"github.com/majd/ipatool/v2/pkg/appstore"
	"github.com/majd/ipatool/v2/pkg/credentialstore"
	"github.com/majd/ipatool/v2/pkg/http"
	"github.com/majd/ipatool/v2/pkg/log"
	"github.com/majd/ipatool/v2/pkg/util"
	"github.com/majd/ipatool/v2/pkg/util/machine"
	"github.com/majd/ipatool/v2/pkg/util/operatingsystem"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var dependencies = Dependencies{}
var keychainPassphrase string
var keychainPassphraseHostSecret string

type Dependencies struct {
	Logger          log.Logger
	OS              operatingsystem.OperatingSystem
	Machine         machine.Machine
	CookieJar       http.CookieJar
	CredentialStore credentialstore.Store
	AppStore        appstore.AppStore
}

// newLogger returns a new logger instance.
func newLogger(format OutputFormat, verbose bool) log.Logger {
	var writer io.Writer

	switch format {
	case OutputFormatJSON:
		writer = log.NewMaskingWriter(zerolog.SyncWriter(os.Stdout))
	case OutputFormatText:
		writer = log.NewWriter()
	}

	return log.NewLogger(log.Args{
		Verbose: verbose,
		Writer:  writer,
	},
	)
}

// newCookieJar returns a new cookie jar instance.
func newCookieJar(machine machine.Machine) http.CookieJar {
	return util.Must(cookiejar.New(&cookiejar.Options{
		Filename: filepath.Join(machine.HomeDirectory(), ConfigDirectoryName, CookieJarFileName),
	}))
}

func hostSecretValue(secretName string) (string, error) {
	if secretName == "" {
		return "", nil
	}

	value, ok := os.LookupEnv(secretName)
	if !ok {
		return "", fmt.Errorf("host secret %q is not set in the environment", secretName)
	}

	if value == "" {
		return "", fmt.Errorf("host secret %q is empty", secretName)
	}

	return value, nil
}

// newCredentialStore returns a new credential store instance.
func newCredentialStore(machine machine.Machine, logger log.Logger, interactive bool) (credentialstore.Store, error) {
	hostSecretPassphrase := util.Must(hostSecretValue(keychainPassphraseHostSecret))

	if keychainPassphrase != "" && hostSecretPassphrase != "" {
		return nil, errors.New("both --keychain-passphrase and --keychain-passphrase-host-secret were provided; use only one input channel")
	}

	ring, err := keyring.Open(keyring.Config{
		AllowedBackends: []keyring.BackendType{
			keyring.KeychainBackend,
			keyring.SecretServiceBackend,
			keyring.FileBackend,
		},
		ServiceName: KeychainServiceName,
		FileDir:     filepath.Join(machine.HomeDirectory(), ConfigDirectoryName),
		FilePasswordFunc: func(s string) (string, error) {
			if hostSecretPassphrase != "" {
				return hostSecretPassphrase, nil
			}

			if keychainPassphrase != "" {
				return keychainPassphrase, nil
			}

			if !interactive {
				return "", errors.New("keychain passphrase is required in non-interactive mode; provide --keychain-passphrase or --keychain-passphrase-host-secret")
			}

			path := strings.Split(s, " unlock ")[1]
			logger.Log().Msgf("enter passphrase to unlock %s (this is separate from your Apple ID password): ", path)
			bytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return "", fmt.Errorf("failed to read password: %w", err)
			}

			password := string(bytes)
			password = strings.Trim(password, "\n")
			password = strings.Trim(password, "\r")

			return password, nil
		},
	})
	if err != nil {
		return nil, err
	}

	return credentialstore.NewNative(credentialstore.NativeArgs{Keyring: ring}), nil
}

// initWithCommand initializes the dependencies of the command.
func initWithCommand(cmd *cobra.Command) {
	verbose := cmd.Flag("verbose").Value.String() == "true"
	interactive, _ := cmd.Context().Value("interactive").(bool)
	format := util.Must(OutputFormatFromString(cmd.Flag("format").Value.String()))

	dependencies.Logger = newLogger(format, verbose)
	dependencies.OS = operatingsystem.New()
	dependencies.Machine = machine.New(machine.Args{OS: dependencies.OS})
	dependencies.CookieJar = newCookieJar(dependencies.Machine)
	dependencies.CredentialStore = util.Must(newCredentialStore(dependencies.Machine, dependencies.Logger, interactive))
	dependencies.AppStore = appstore.NewAppStore(appstore.Args{
		CookieJar:       dependencies.CookieJar,
		OperatingSystem: dependencies.OS,
		CredentialStore: dependencies.CredentialStore,
		Machine:         dependencies.Machine,
	})

	util.Must("", createConfigDirectory(dependencies.OS, dependencies.Machine))
}

// createConfigDirectory creates the configuration directory for the CLI tool, if needed.
func createConfigDirectory(os operatingsystem.OperatingSystem, machine machine.Machine) error {
	configDirectoryPath := filepath.Join(machine.HomeDirectory(), ConfigDirectoryName)
	_, err := os.Stat(configDirectoryPath)

	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(configDirectoryPath, 0700)
		if err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("could not read metadata: %w", err)
	}

	return nil
}
