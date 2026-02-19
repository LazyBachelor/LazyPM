package service

type Config struct {
	RootCmd               string
	WebAddress            string
	BeadsDBPath           string
	IssuePrefix           string
	StatisticsStoragePath string
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
