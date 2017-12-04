package ssp

import (
	"fmt"
	"github.com/go-ini/ini"
	"net/http"
	"os"
	"reflect"
)

// Config contains the details of the Dashboard endpoint needed to send API requests from this SDK.
//
// SDK needs to be configured with the following variables:
// * DASHBOARD_EMAIL: set to the email of the account you have generated the token from
// * DASHBOARD_TOKEN: access token you've just obtained
// * DASHBOARD_URL: address of the dashboard you are connecting to.
type Config struct {
	Email   string `ini:"DASHBOARD_EMAIL" env:"DASHBOARD_EMAIL"`
	Token   string `ini:"DASHBOARD_TOKEN" env:"DASHBOARD_TOKEN"`
	BaseURL string `ini:"DASHBOARD_URL" env:"DASHBOARD_URL"`
}

// NewDefaultConfig loads base configuration from $HOME/.dashboard.env, but also allows overriding
// from environment variables.
func NewDefaultConfig() *Config {
	creds := NewHomeConfig()
	overrideFromEnv(creds)
	return creds
}

// NewEnvConfig loads configuration from the environment variables.
func NewEnvConfig() *Config {
	c := new(Config)
	overrideFromEnv(c)
	return c
}

// NewHomeConfig loads configuration from $HOME/.dashboard.env.
func NewHomeConfig() *Config {
	path := fmt.Sprintf("%s/.dashboard.env", os.Getenv("HOME"))
	return NewIniConfig(path)
}

// NewIniConfig loads configuration from an arbitrary path.
func NewIniConfig(path string) *Config {
	c := new(Config)
	ini.MapTo(c, path)
	return c
}

func (c *Config) Client() *http.Client {
	t := &BasicAuthTransport{
		Username: c.Email,
		Password: c.Token,
	}

	return t.Client()
}

func overrideFromEnv(c *Config) {
	t := reflect.TypeOf(c).Elem()
	v := reflect.ValueOf(c).Elem()
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("env")
		if os.Getenv(tag) != "" {
			v.Field(i).SetString(os.Getenv(tag))
		}
	}
}
