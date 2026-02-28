package main

import (
	"context"

	issuesCmd "github.com/LazyBachelor/LazyPM/internal/commands/issues"
	surveyCmd "github.com/LazyBachelor/LazyPM/internal/commands/survey"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/cli"
)

func main() {
	if err := cli.NewCli(issuesCmd.RootCmd).Run(context.Background(), models.BaseConfig); err != nil {
		return
	}
}

func init() {
	issuesCmd.RootCmd.AddCommand(issuesCmd.GetCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.ListCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.CloseCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.CreateCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.DeleteCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.UpdateCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.CommentCmd)
	issuesCmd.RootCmd.AddCommand(issuesCmd.CommentsCmd)
	issuesCmd.RootCmd.AddCommand(surveyCmd.StatusCmd)
}
