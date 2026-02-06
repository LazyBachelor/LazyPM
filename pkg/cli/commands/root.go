package commands

import (
	"context"
	"os"

	"github.com/LazyBachelor/LazyPM/internal/service"
	"github.com/charmbracelet/fang"

	"github.com/spf13/cobra"
)

var svc *service.Services

var rootCmd = &cobra.Command{
	Short: "Project Management CLI",
	Long:  `Project Management CLI for managing issues and tasks.`,
}

func Execute(services *service.Services) error {
	SetServices(services)
	return fang.Execute(context.Background(), rootCmd,
		fang.WithColorSchemeFunc(fang.AnsiColorScheme),
		fang.WithNotifySignal(os.Interrupt, os.Kill))
}

func SetServices(services *service.Services) {
	svc = services
	rootCmd.Use = svc.Config.RootCmd
}

func ExecuteWithArgs(args []string) error {
	rootCmd.SetArgs(args)
	return fang.Execute(context.Background(), rootCmd,
		fang.WithColorSchemeFunc(fang.AnsiColorScheme))
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = false
	rootCmd.AddGroup(&cobra.Group{ID: "help", Title: "Helping Commands"})
	rootCmd.SetCompletionCommandGroupID("help")
	rootCmd.SetHelpCommandGroupID("help")
}
