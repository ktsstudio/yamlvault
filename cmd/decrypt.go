package cmd

import (
	"github.com/ktsstudio/yamlvault/pkg/processor"
	"github.com/ktsstudio/yamlvault/pkg/vaulter"
	"github.com/ktsstudio/yamlvault/pkg/yamldoc"
	"github.com/spf13/cobra"
)

var conf struct {
	InputFile  string
	OutputFile string
}

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use: "decrypt",
	RunE: func(cmd *cobra.Command, args []string) error {
		doc, err := yamldoc.New(conf.InputFile)
		if err != nil {
			return err
		}

		vaulter, err := vaulter.New()
		if err != nil {
			return err
		}

		processor, err := processor.New(vaulter)
		if err != nil {
			return err
		}

		if err := processor.Process(doc); err != nil {
			return err
		}

		if err := doc.Save(conf.OutputFile); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	decryptCmd.PersistentFlags().StringVarP(&conf.InputFile, "input", "i", "", "Input file")
	decryptCmd.PersistentFlags().StringVarP(&conf.OutputFile, "output", "o", "", "Output file")
	_ = decryptCmd.MarkPersistentFlagRequired("input")
	_ = decryptCmd.MarkPersistentFlagRequired("output")
}
