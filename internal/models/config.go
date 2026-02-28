package models

type Config struct {
	AutoInit              bool
	RootCmd               string
	WebAddress            string
	BeadsDBPath           string
	IssuePrefix           string
	StatisticsStoragePath string
}

var BaseConfig = Config{
	AutoInit:              false,
	RootCmd:               "pm",
	IssuePrefix:           "pm",
	WebAddress:            ":8080",
	BeadsDBPath:           "./.pm/db.db",
	StatisticsStoragePath: "./.pm/stats.json",
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
