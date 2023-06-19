package cmd

import (
	"fmt"
	"log"

	app "passKeeper/internal/cmd/app"
	cc "passKeeper/internal/cmd/tui/new/creditcard"
	f "passKeeper/internal/cmd/tui/new/file"
	kv "passKeeper/internal/cmd/tui/new/kv"
	txt "passKeeper/internal/cmd/tui/new/txt"
	conf "passKeeper/internal/cmd/tui/setup"
	sec "passKeeper/internal/models/secret"

	"github.com/spf13/cobra"
)

var (
	username, password string
)
var (
	rootCmd = &cobra.Command{
		Use:   "passKeeper",
		Short: "A tool for managing secrets.",
		Long:  "passKeeper is an advanced tool that facilitates the secure handling and management of secrets.",
	}
	newCmd = &cobra.Command{
		Use:   "new",
		Short: "Generate a new secret.",
		Long:  "Generate a new secret of a specific type, options include key-value pair (kv), credit card details (cc), text (txt), or file.",
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(root *cobra.Command) error {
	err := root.Execute()
	if err != nil {
		return err
	}
	return nil
}

func NewRootCommand() *cobra.Command {
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(describeCmd)
	rootCmd.AddCommand(dumpCmd)
	newCmd.AddCommand(newTextCmd)
	newCmd.AddCommand(newKVCmd)
	newCmd.AddCommand(newCCCmd)
	newCmd.AddCommand(newFileCmd)

	return rootCmd
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to an existing passKeeper account.",
	Long:  "Initiate login process for a user with an existing passKeeper account. Enter your login credentials when prompted.",
	RunE: func(cmd *cobra.Command, args []string) error { // Replace Run with RunE
		login := true
		if err := conf.SetupTui(login); err != nil {
			return fmt.Errorf("could not start passKeeper: %s", err)
		}
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Delete local passKeeper configuration.",
	Long:  "Clear all locally stored passKeeper configuration. This ensures that your sensitive data is safe and secure after usage.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := app.ClearLocalData(); err != nil {
			return fmt.Errorf("remove local secret and files attempt has failed: %s", err)
		}
		return nil
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up initial passKeeper configuration.",
	Long:  "Establish initial configurations for passKeeper. This includes setting up the username and password. If these parameters are not provided, a user interface will guide the setup process.",
	RunE: func(cmd *cobra.Command, args []string) error {
		login := false
		if username == "" || password == "" {
			if err := conf.SetupTui(login); err != nil {
				return fmt.Errorf("could not start passKeeper: %s", err)
			}
		} else {
			if err := app.SetUsername(username); err != nil {
				return err
			}
			if err := app.SetKey(app.AppName, password); err != nil {
				return err
			}
		}
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a secret.",
	Long:  "Remove a secret stored in passKeeper by its unique identifier. The secret will be permanently deleted from the system.",
	Run: func(cmd *cobra.Command, args []string) {
		app := app.GetApplication()

		for _, v := range args {
			app.DeleteSecret(v)

		}

	},
}
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Export binary secret data.",
	Long:  "Extract and export the binary data of a secret by its unique identifier on disk. This is useful for backing up or transferring secret information.",
	Run: func(cmd *cobra.Command, args []string) {
		app := app.GetApplication()

		for _, v := range args {
			app.DumpSecret(v)

		}

	},
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Modify a secret.",
	Long:  "Edit the contents of a secret stored in passKeeper by its unique identifier. Depending on the type of the secret (key-value pair, text, or credit card), the corresponding user interface will be invoked for modification.",
	RunE: func(cmd *cobra.Command, args []string) error {
		app := app.GetApplication()

		if len(args) > 1 || len(args) == 0 {
			return fmt.Errorf("wrong number of arguments. expected only one id")
		}

		secret, err := app.GetSecret(args[0])
		if err != nil {
			return fmt.Errorf("cannot get secret")
		}

		decodedSecret, err := sec.GetDecodedSecrets([]sec.Secret{*secret})
		if err != nil {
			return fmt.Errorf("cannot decode secret")
		}

		switch v := decodedSecret[0].Value.(type) {
		case *sec.KeyValue:
			if err := kv.EditKVTui(*v, secret.Metadata, secret.ID); err != nil {
				return fmt.Errorf("could not start passKeeper: %s", err)
			}
		case *sec.Text:
			if err := txt.EditTextTui(*v, secret.Metadata, secret.ID); err != nil {
				return fmt.Errorf("could not start passKeeper: %s", err)
			}
		case *sec.CreditCard:
			if err := cc.EditCCTui(*v, secret.Metadata, secret.ID); err != nil {
				return fmt.Errorf("could not start passKeeper: %s", err)
			}
		case *sec.ByteSlice:
		default:
			return nil
		}
		return nil
	},
}

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Display a secret's details.",
	Long:  "Provide comprehensive details of a secret stored in passKeeper by its unique identifier. This includes the metadata, value, and other associated information.",
	Run: func(cmd *cobra.Command, args []string) {
		app := app.GetApplication()

		if len(args) > 1 || len(args) == 0 {
			log.Printf("%s", "Wrong number of arguments. Expected only one id.")
			return
		} else {

			secret, err := app.GetSecret(args[0])
			if err != nil {
				log.Printf("%s", "Cannot get secret")
				return
			}

			decodedSecret, err := sec.GetDecodedSecrets([]sec.Secret{*secret})
			if err != nil {
				log.Printf("%s", "Cannot decode secret")
				return
			}
			fmt.Printf("Secret Id: %d \nSecret metadata: %s\n", secret.ID, secret.Metadata)
			fmt.Printf("Secret value:\n%s", decodedSecret[0].ValueToString())

		}

	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available secrets.",
	Long:  "Display a list of all secrets currently stored in passKeeper. This includes the secret's identifier, value, and any associated metadata.",
	Run: func(cmd *cobra.Command, args []string) {

		appl := app.GetApplication()
		list := appl.ListSecrets()
		decodedSecrets, err := sec.GetDecodedSecrets(*list)
		if err != nil {
			log.Printf("%s", err.Error())
		}

		app.List(*appl, decodedSecrets)

	},
}

var newTextCmd = &cobra.Command{
	Use:   "txt",
	Short: "Create a new text secret.",
	Long:  "Generate a new secret of the 'text' type. The secret contain plain text data.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := txt.NewTextTui(); err != nil {
			return fmt.Errorf("could not start passKeeper: %s", err)
		}
		return nil
	},
}

var newFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Create a new file secret.",
	Long:  "Generate a new secret of the 'file' type. The secret can contain binary data.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := f.FileTui(); err != nil {
			return fmt.Errorf("could not start passKeeper: %s", err)
		}
		return nil
	},
}

var newKVCmd = &cobra.Command{
	Use:   "kv",
	Short: "Create a new key-value secret.",
	Long:  "Generate a new secret of the 'key-value' type. The secret can contain a key-value pair.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := kv.NewKVTui(); err != nil {
			return fmt.Errorf("could not start passKeeper: %s", err)
		}
		return nil
	},
}

var newCCCmd = &cobra.Command{
	Use:   "cc",
	Short: "Create a new credit card secret.",
	Long:  "Generate a new secret of the 'credit card' type. The secret can contain credit card information.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cc.NewCCTui(); err != nil {
			return fmt.Errorf("could not start passKeeper: %s", err)
		}
		return nil
	},
}
