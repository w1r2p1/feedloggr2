package feedloggr2

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

type FeedConfig struct {
	Title string
	URL   string
}

type Config struct {
	Verbose    bool
	Database   string
	OutputPath string
	Feeds      []*FeedConfig
}

func (c *Config) String() string {
	var b bytes.Buffer
	err := toml.NewEncoder(&b).Encode(c)
	if err != nil {
		return err.Error()
	}
	return b.String()
}

func NewConfig() *Config {
	c := &Config{
		Verbose:    true,
		Database:   ".feedloggr2.db",
		OutputPath: "./feeds",
		Feeds: []*FeedConfig{
			&FeedConfig{"Example", "https://example.com/rss"},
		},
	}

	return c
}

func LoadConfig(path string) (*Config, error) {
	c := &Config{}
	_, err := toml.DecodeFile(path, &c)
	return c, err
}