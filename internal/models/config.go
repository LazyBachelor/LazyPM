package models

import "os"

type Config struct {
	DbUri                 string
	AutoInit              bool
	RootCmd               string
	WebAddress            string
	BeadsDBPath           string
	IssuePrefix           string
	StatisticsStoragePath string
	ActionLogger          func(string)
}

var BaseConfig = Config{
	DbUri:                 "",
	AutoInit:              false,
	RootCmd:               "pm",
	IssuePrefix:           "pm",
	WebAddress:            ":8080",
	BeadsDBPath:           "./.pm/db.db",
	StatisticsStoragePath: "./.pm/stats.json",
}

func (c Config) LoadFromEnv() Config {
	c.DbUri = os.Getenv("DB_URI")
	return c
}

func (c Config) WithDbUri(uri string) Config {
	c.DbUri = uri
	return c
}

func (c Config) WithAutoInit(autoInit bool) Config {
	c.AutoInit = autoInit
	return c
}

func (c Config) WithRootCmd(rootCmd string) Config {
	c.RootCmd = rootCmd
	return c
}

func (c Config) WithWebAddress(webAddress string) Config {
	c.WebAddress = webAddress
	return c
}

func (c Config) WithBeadsDBPath(beadsDBPath string) Config {
	c.BeadsDBPath = beadsDBPath
	return c
}

func (c Config) WithIssuePrefix(issuePrefix string) Config {
	c.IssuePrefix = issuePrefix
	return c
}

func (c Config) WithStatisticsStoragePath(path string) Config {
	c.StatisticsStoragePath = path
	return c
}

func (c Config) WithActionLogger(logger func(string)) Config {
	c.ActionLogger = logger
	return c
}
