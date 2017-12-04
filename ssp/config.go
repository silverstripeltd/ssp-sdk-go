package ssp

import (
	"fmt"
	"github.com/go-ini/ini"
	"net/http"
	"os"
	"reflect"
)

type Config struct {
	Email   string `ini:"DASHBOARD_EMAIL" env:"DASHBOARD_EMAIL"`
	Token   string `ini:"DASHBOARD_TOKEN" env:"DASHBOARD_TOKEN"`
	BaseURL string `ini:"DASHBOARD_URL" env:"DASHBOARD_URL"`
}

func NewDefaultConfig() *Config {
	creds := NewHomeConfig()
	overrideFromEnv(creds)
	return creds
}

func NewEnvConfig() *Config {
	c := new(Config)
	overrideFromEnv(c)
	return c
}

func NewHomeConfig() *Config {
	path := fmt.Sprintf("%s/.dashboard.env", os.Getenv("HOME"))
	return NewIniConfig(path)
}

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
