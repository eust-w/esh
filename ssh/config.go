package ssh

import (
	"time"
)

type Config struct {
	User     string
	Host     string
	Port     int
	Password string

	StickySession          bool
	DisableAgentForwarding bool
	HandshakeTimeout       time.Duration
	KeepAliveInterval      time.Duration
	Timeout                time.Duration
}

func (c *Config) WithUser(user string) *Config {
	if user == "" {
		user = "root"
	}
	c.User = user
	return c
}

func (c *Config) WithHost(host string) *Config {
	if host == "" {
		host = "localhost"
	}
	c.Host = host
	return c
}

func (c *Config) WithPassword(password string) *Config {
	c.Password = password
	return c
}

func (c *Config) SetKeys(keyfiles []string) *Config {
	if keyfiles == nil {
		return c
	}
	t := make([]string, len(keyfiles))
	copy(t, keyfiles)
	return c
}
