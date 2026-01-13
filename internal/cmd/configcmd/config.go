package configcmd

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/piekstra/newrelic-cli/internal/cmd/root"
	"github.com/piekstra/newrelic-cli/internal/config"
	"github.com/piekstra/newrelic-cli/internal/confirm"
	"github.com/piekstra/newrelic-cli/internal/validate"
	"github.com/piekstra/newrelic-cli/internal/view"
)

// Register adds the config commands to the root command
func Register(rootCmd *cobra.Command, opts *root.Options) {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configure newrelic-cli credentials",
	}

	configCmd.AddCommand(newSetAPIKeyCmd(opts))
	configCmd.AddCommand(newDeleteAPIKeyCmd(opts))
	configCmd.AddCommand(newSetAccountIDCmd(opts))
	configCmd.AddCommand(newDeleteAccountIDCmd(opts))
	configCmd.AddCommand(newSetRegionCmd(opts))
	configCmd.AddCommand(newShowCmd(opts))
	configCmd.AddCommand(newFixPermissionsCmd(opts))

	rootCmd.AddCommand(configCmd)
}

func newSetAPIKeyCmd(opts *root.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "set-api-key [key]",
		Short: "Set the New Relic API key",
		Long: `Set the New Relic API key for authentication.

On macOS: Key is stored securely in the system Keychain.
On Linux: Key is stored in ~/.config/newrelic-cli/credentials (file permissions 0600).

If no key is provided as an argument, you will be prompted to enter it.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetAPIKey(opts, args)
		},
	}
}

func runSetAPIKey(opts *root.Options, args []string) error {
	v := opts.View()

	if !config.IsSecureStorage() {
		v.Warning("Warning: On Linux, your API key will be stored in a config file")
		v.Println("         (~/.config/newrelic-cli/credentials) with restricted permissions (0600).")
		v.Println("         This is less secure than macOS Keychain storage.")
		v.Println("")
	}

	var apiKey string
	if len(args) > 0 {
		apiKey = args[0]
	} else {
		fmt.Fprint(opts.Stdout, "Enter New Relic API key: ")
		reader := bufio.NewReader(opts.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		apiKey = strings.TrimSpace(input)
	}

	// Validate API key
	warning, err := validate.APIKey(apiKey)
	if err != nil {
		return err
	}
	if warning != "" {
		v.Warning("Warning: " + warning)
	}

	if err := config.SetAPIKey(apiKey); err != nil {
		return fmt.Errorf("failed to store API key: %w", err)
	}

	if config.IsSecureStorage() {
		v.Success("API key stored securely in Keychain")
	} else {
		v.Success("API key stored in ~/.config/newrelic-cli/credentials")
	}
	return nil
}

// deleteAPIKeyOptions holds options for the delete-api-key command
type deleteAPIKeyOptions struct {
	*root.Options
	force bool
}

func newDeleteAPIKeyCmd(opts *root.Options) *cobra.Command {
	deleteOpts := &deleteAPIKeyOptions{Options: opts}

	cmd := &cobra.Command{
		Use:   "delete-api-key",
		Short: "Delete the stored New Relic API key",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeleteAPIKey(deleteOpts)
		},
	}

	cmd.Flags().BoolVarP(&deleteOpts.force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

func runDeleteAPIKey(opts *deleteAPIKeyOptions) error {
	v := opts.View()

	if !opts.force {
		p := &confirm.Prompter{
			In:  opts.Stdin,
			Out: opts.Stderr,
		}
		if !p.Confirm("Delete stored API key?") {
			v.Warning("Operation canceled")
			return nil
		}
	}

	if err := config.DeleteAPIKey(); err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	if config.IsSecureStorage() {
		v.Success("API key deleted from Keychain")
	} else {
		v.Success("API key deleted from config file")
	}
	return nil
}

func newSetAccountIDCmd(opts *root.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "set-account-id <account-id>",
		Short: "Set the New Relic account ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetAccountID(opts, args[0])
		},
	}
}

func runSetAccountID(opts *root.Options, accountID string) error {
	v := opts.View()

	// Validate account ID
	if err := validate.AccountID(accountID); err != nil {
		return err
	}

	if err := config.SetAccountID(accountID); err != nil {
		return fmt.Errorf("failed to store account ID: %w", err)
	}

	if config.IsSecureStorage() {
		v.Success("Account ID stored securely in Keychain")
	} else {
		v.Success("Account ID stored in config file")
	}
	return nil
}

// deleteAccountIDOptions holds options for the delete-account-id command
type deleteAccountIDOptions struct {
	*root.Options
	force bool
}

func newDeleteAccountIDCmd(opts *root.Options) *cobra.Command {
	deleteOpts := &deleteAccountIDOptions{Options: opts}

	cmd := &cobra.Command{
		Use:   "delete-account-id",
		Short: "Delete the stored New Relic account ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeleteAccountID(deleteOpts)
		},
	}

	cmd.Flags().BoolVarP(&deleteOpts.force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

func runDeleteAccountID(opts *deleteAccountIDOptions) error {
	v := opts.View()

	if !opts.force {
		p := &confirm.Prompter{
			In:  opts.Stdin,
			Out: opts.Stderr,
		}
		if !p.Confirm("Delete stored account ID?") {
			v.Warning("Operation canceled")
			return nil
		}
	}

	if err := config.DeleteAccountID(); err != nil {
		return fmt.Errorf("failed to delete account ID: %w", err)
	}

	if config.IsSecureStorage() {
		v.Success("Account ID deleted from Keychain")
	} else {
		v.Success("Account ID deleted from config file")
	}
	return nil
}

func newSetRegionCmd(opts *root.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "set-region <region>",
		Short: "Set the New Relic region (US or EU)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetRegion(opts, args[0])
		},
	}
}

func runSetRegion(opts *root.Options, region string) error {
	v := opts.View()

	region = strings.ToUpper(region)

	// Validate region
	if err := validate.Region(region); err != nil {
		return err
	}

	if err := config.SetRegion(region); err != nil {
		return fmt.Errorf("failed to store region: %w", err)
	}

	v.Success("Region set to %s", region)
	return nil
}

func newShowCmd(opts *root.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current configuration status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShow(opts)
		},
	}
}

// ConfigStatus represents configuration status for JSON output
// NOTE: API key value is intentionally NOT included for security
type ConfigStatus struct {
	APIKeyConfigured bool   `json:"api_key_configured"`
	APIKeySource     string `json:"api_key_source,omitempty"`
	AccountID        string `json:"account_id,omitempty"`
	AccountIDSource  string `json:"account_id_source,omitempty"`
	Region           string `json:"region"`
	RegionSource     string `json:"region_source"`
	StorageType      string `json:"storage_type"`
}

func runShow(opts *root.Options) error {
	v := opts.View()
	status := config.GetCredentialStatus()

	// Check for permission warnings (Linux only)
	if warning := config.CheckPermissions(); warning != "" {
		v.Warning(warning)
		v.Println("Run 'newrelic-cli config fix-permissions' to correct this")
		v.Println("")
	}

	// Build configuration status
	configStatus := ConfigStatus{
		Region:      config.GetRegion(),
		StorageType: "config_file",
	}

	if config.IsSecureStorage() {
		configStatus.StorageType = "keychain"
	}

	// API Key
	var apiKeyMasked string
	if apiKey, err := config.GetAPIKey(); err == nil {
		configStatus.APIKeyConfigured = true
		if status["api_key_env"] {
			configStatus.APIKeySource = "environment"
		} else {
			configStatus.APIKeySource = "stored"
		}
		// Mask API key for display (first 8 + last 4)
		if len(apiKey) > 12 {
			apiKeyMasked = apiKey[:8] + strings.Repeat("*", len(apiKey)-12) + apiKey[len(apiKey)-4:]
		} else {
			apiKeyMasked = strings.Repeat("*", len(apiKey))
		}
	}

	// Account ID
	if accountID, err := config.GetAccountID(); err == nil {
		configStatus.AccountID = accountID
		if status["account_id_env"] {
			configStatus.AccountIDSource = "environment"
		} else {
			configStatus.AccountIDSource = "stored"
		}
	}

	// Region source
	if status["region_stored"] {
		configStatus.RegionSource = "stored"
	} else if status["region_env"] {
		configStatus.RegionSource = "environment"
	} else {
		configStatus.RegionSource = "default"
	}

	// JSON output - never include API key value
	if v.Format == view.FormatJSON {
		return v.JSON(configStatus)
	}

	// Table/Plain output
	v.Println("Configuration Status:")
	v.Println("")

	// API Key
	if configStatus.APIKeyConfigured {
		v.Print("  API Key:    %s (%s)\n", apiKeyMasked, configStatus.APIKeySource)
	} else {
		v.Println("  API Key:    Not configured")
	}

	// Account ID
	if configStatus.AccountID != "" {
		v.Print("  Account ID: %s (%s)\n", configStatus.AccountID, configStatus.AccountIDSource)
	} else {
		v.Println("  Account ID: Not configured")
	}

	// Region
	v.Print("  Region:     %s (%s)\n", configStatus.Region, configStatus.RegionSource)

	v.Println("")

	// Storage type
	if config.IsSecureStorage() {
		v.Println("Storage: macOS Keychain (secure)")
	} else {
		v.Println("Storage: Config file (~/.config/newrelic-cli/credentials)")
	}

	return nil
}

func newFixPermissionsCmd(opts *root.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "fix-permissions",
		Short: "Fix config file permissions to 0600 (Linux only)",
		Long: `Fix the permissions on the credentials file to ensure they are secure.

On Linux, the credentials file should have permissions 0600 (owner read/write only).
On macOS, this command has no effect as credentials are stored in the Keychain.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFixPermissions(opts)
		},
	}
}

func runFixPermissions(opts *root.Options) error {
	v := opts.View()

	if config.IsSecureStorage() {
		v.Println("On macOS, credentials are stored in the Keychain - no file permissions to fix")
		return nil
	}

	if err := config.FixPermissions(); err != nil {
		return fmt.Errorf("failed to fix permissions: %w", err)
	}

	v.Success("Permissions fixed to 0600")
	return nil
}
