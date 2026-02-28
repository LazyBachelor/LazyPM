package main

import (
	"context"

	issues "github.com/LazyBachelor/LazyPM/internal/commands/issues"
	survey "github.com/LazyBachelor/LazyPM/internal/commands/survey"
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/pkg/cli"
)

func main() {
	if err := cli.New(issues.RootCmd).Run(context.Background(), models.BaseConfig); err != nil {
		return
	}
}

func init() {
	issues.RootCmd.AddCommand(issues.GetCmd)
	issues.RootCmd.AddCommand(issues.ListCmd)
	issues.RootCmd.AddCommand(issues.CloseCmd)
	issues.RootCmd.AddCommand(issues.CreateCmd)
	issues.RootCmd.AddCommand(issues.DeleteCmd)
	issues.RootCmd.AddCommand(issues.UpdateCmd)
	issues.RootCmd.AddCommand(issues.CommentCmd)
	issues.RootCmd.AddCommand(issues.CommentsCmd)
	issues.RootCmd.AddCommand(survey.StatusCmd)
}
