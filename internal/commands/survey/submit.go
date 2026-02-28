package survey

import "github.com/spf13/cobra"

var SubmitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit your survey responses",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Submitting responses and metrics...")
		return nil
	},
}
