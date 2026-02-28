package utils

import "github.com/spf13/cobra"

// CompletionFunc returns a function that provides shell completion for the given options.
func CompletionFunc(options []string) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return options, cobra.ShellCompDirectiveDefault
	}
}
