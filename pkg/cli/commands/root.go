package commands

import (
	"github.com/LazyBachelor/LazyPM/internal/service"

	"github.com/spf13/cobra"
)

var svc *service.Services

var rootCmd = &cobra.Command{
	Short: "Project Management CLI",
	Long:  `Project Management CLI for managing issues and tasks.`,
}

func Execute(services *service.Services) error {
	svc = services
	rootCmd.Use = svc.Config.RootCmd
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = false
	rootCmd.AddGroup(&cobra.Group{ID: "help", Title: "Helping Commands"})
	rootCmd.SetCompletionCommandGroupID("help")
	rootCmd.SetHelpCommandGroupID("help")
}
