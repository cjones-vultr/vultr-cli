// Package account provides the account functionality for the CLI
package account

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	accountLong    = `Retrieve information about your account.`
	accountExample = `
	# Full example
	vultr-cli account
	`
)

// Interface for account
type AccountInterface interface {
	Get() (*govultr.Account, error)
	validate(cmd *cobra.Command, args []string)
}

// Options for account
type Options struct {
	Base *cli.Base
}

// NewAccountOptions returns Options struct
func NewAccountOptions(base *cli.Base) *Options {
	return &Options{Base: base}
}

// NewCmdAccount creates a cobra command for Account
func NewCmdAccount(base *cli.Base) *cobra.Command {
	o := NewAccountOptions(base)

	cmd := &cobra.Command{
		Use:     "account",
		Short:   "get account information",
		Long:    accountLong,
		Example: accountExample,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !o.Base.HasAuth {
				return errors.New(utils.APIKeyError)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.validate(cmd, args)
			account, err := o.Get()

			o.Base.Printer.Display(&AccountPrinter{Account: account}, err)
		},
	}

	return cmd
}

func (o *Options) validate(cmd *cobra.Command, args []string) {
	o.Base.Printer.Output = viper.GetString("output")
}

// Get account information
func (o *Options) Get() (*govultr.Account, error) {
	account, _, err := o.Base.Client.Account.Get(context.Background())
	if err != nil {
		fmt.Printf("Error getting account information : %v\n", err)
		os.Exit(1)
	}

	return account, nil
}
