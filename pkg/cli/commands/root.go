package commands

import (
	"github.com/LazyBachelor/LazyPM/internal/service"

	"github.com/spf13/cobra"
)

var svc *service.Services

var rootCmd = &cobra.Command{
	Use:   "pm",
	Short: "Project Management CLI",
	Long:  `Project Management CLI for managing issues and tasks.`,
}

func Execute(services *service.Services) error {
	svc = services
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = false
	rootCmd.AddGroup(&cobra.Group{ID: "other", Title: "Helping Commands"})
	rootCmd.SetCompletionCommandGroupID("other")
	rootCmd.SetHelpCommandGroupID("other")
}
