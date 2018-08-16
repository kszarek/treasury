package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/AirHelp/treasury/client"
	"github.com/spf13/cobra"
)

const (
	templateSuccessMsg                     = "File with secrets successfully generated"
	templateErrorMissingSourceFile         = "Missing source file path"
	templateErrorMissingDestinationFile    = "Missing destination file path"
	templateCommandSourceFileArgument      = "src"
	templateCommandDestinationFileArgument = "dst"
	templateCommandPermissionFileArgument  = "perms"
	templateCommandAppendArgument          = "append"
)

var (
	templateCmd = &cobra.Command{
		Use:   "template --src TEMPLATE_FILE --dst DESTINATION_FILE",
		Short: "Generates a file with secrets from given template",
		Long:  `Generates a file with secrets from given template`,
		RunE:  template,
	}
	append []string
)

func init() {
	RootCmd.AddCommand(templateCmd)
	templateCmd.PersistentFlags().String(templateCommandSourceFileArgument, "", "template file path")
	templateCmd.PersistentFlags().String(templateCommandDestinationFileArgument, "", "destination file path")
	templateCmd.PersistentFlags().Int(templateCommandPermissionFileArgument, 0, "destination file permission, e.g.: 0644")
	templateCmd.PersistentFlags().StringArrayVar(&append, templateCommandAppendArgument, []string{}, "variable suffix, e.g: --append \"DATABASE_URL:?pool=10\"")
}

func template(cmd *cobra.Command, args []string) error {
	sourceFilePath, err := cmd.Flags().GetString(templateCommandSourceFileArgument)
	if err != nil {
		return err
	}
	if sourceFilePath == "" {
		return errors.New(templateErrorMissingSourceFile)
	}

	destinationFilePath, err := cmd.Flags().GetString(templateCommandDestinationFileArgument)
	if err != nil {
		return err
	}
	if destinationFilePath == "" {
		return errors.New(templateErrorMissingDestinationFile)
	}

	perms, err := cmd.Flags().GetInt(templateCommandPermissionFileArgument)
	if err != nil {
		return err
	}

	treasury, err := client.New(&client.Options{
		Region:       s3Region,
		S3BucketName: treasuryS3,
		Append:       append,
	})
	if err != nil {
		return err
	}
	if err := treasury.Template(sourceFilePath, destinationFilePath, os.FileMode(perms)); err != nil {
		return err
	}
	fmt.Println(templateSuccessMsg)
	return nil
}
