package commands

import (
	"beadstest/internal/service"

	"github.com/spf13/cobra"
)

var svc *service.Service

var rootCmd = &cobra.Command{
	Use:   "pm",
	Short: "Project Management CLI",
	Long:  `Project Management CLI for managing issues and tasks.`,
}

func Execute(beadsService *service.Service) error {
	svc = beadsService
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(createCmd)
}
