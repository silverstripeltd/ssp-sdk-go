package ssp

import (
	"os"
	"os/exec"
	"testing"
)

func TestNewIniConfig(t *testing.T) {
	c := NewIniConfig("testdata/.dashboard.env")
	checkConf(t, c)
}

func TestNewEnvConfig(t *testing.T) {
	if os.Getenv("EXEC_THIS_WITH_ENV") == "1" {
		c := NewEnvConfig()
		checkConf(t, c)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestNewEnvConfig")
	cmd.Env = append(os.Environ(), []string{
		"DASHBOARD_EMAIL=admin",
		"DASHBOARD_TOKEN=token",
		"DASHBOARD_URL=http://localhost",
		"EXEC_THIS_WITH_ENV=1",
	}...)
	stdout, err := cmd.Output()
	if err != nil {
		t.Fatalf("Inner cmd exec failed: %s", stdout)
	}
}

func TestNewHomeConfig(t *testing.T) {
	if os.Getenv("EXEC_THIS_WITH_ENV") == "1" {
		c := NewDefaultConfig()
		checkConf(t, c)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestNewHomeConfig")
	cmd.Env = append(os.Environ(), []string{
		"HOME=testdata",
		"EXEC_THIS_WITH_ENV=1",
	}...)
	stdout, err := cmd.Output()
	if err != nil {
		t.Fatalf("Inner cmd exec failed: %s", stdout)
	}

}

func TestNewDefaultConfig(t *testing.T) {
	if os.Getenv("EXEC_THIS_WITH_ENV") == "1" {
		c := NewDefaultConfig()
		if c.Email != "override" {
			t.Error("Email not loaded properly")
		}
		if c.BaseURL != "http://localhost" {
			t.Error("BaseURL not loaded properly")
		}
		if c.Token != "token" {
			t.Error("Token not loaded properly")
		}
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestNewDefaultConfig")
	cmd.Env = append(os.Environ(), []string{
		"DASHBOARD_EMAIL=override",
		"HOME=testdata",
		"EXEC_THIS_WITH_ENV=1",
	}...)
	stdout, err := cmd.Output()
	if err != nil {
		t.Errorf("Inner cmd exec failed: %s", stdout)
	}

}

func TestClient(t *testing.T) {
	tr := NewIniConfig("testdata/.dashboard.env").Client().Transport
	bat := tr.(*BasicAuthTransport)
	if bat.Username != "admin" {
		t.Error("Email not used in the transport")
	}
	if bat.Password != "token" {
		t.Error("Token not used in the transport")
	}
}

func checkConf(t *testing.T, c *Config) {
	if c.Email != "admin" {
		t.Error("Email not loaded properly")
	}
	if c.BaseURL != "http://localhost" {
		t.Error("BaseURL not loaded properly")
	}
	if c.Token != "token" {
		t.Error("Token not loaded properly")
	}
}
