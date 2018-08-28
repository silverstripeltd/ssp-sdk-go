package ssp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/blang/semver"
	"github.com/google/jsonapi"
)

type Usage string

const (
	UsageProduction  = "Production"
	UsageUAT         = "UAT"
	UsageTest        = "Test"
	UsageUnspecified = "Unspecified"
)

type Environment struct {
	ID                     string `jsonapi:"primary,environments"`
	Name                   string `jsonapi:"attr,name"`
	BuildSHA               string `jsonapi:"attr,build_sha"`
	MaintenanceDay         time.Weekday
	MaintenanceUnspecified bool
	MaintenanceTime        time.Time
	MaintenanceDuration    time.Duration
	MaintenanceTz          *time.Location
	// TODO nothing beyond X.Y.Z is supported. What about 2.3? Or composer constraints?
	DesiredManifestSha          semver.Version
	CurrentManifestSha          semver.Version
	Stack                       *Stack       `jsonapi:"relation,stack"`
	BaseEnvironment             *Environment `jsonapi:"relation,base_environment"`
	Usage                       Usage
	OriginalUsage               string `jsonapi:"attr,usage"`
	OriginalMaintenanceDay      string `jsonapi:"attr,maintenance_day"`
	OriginalMaintenanceTime     string `jsonapi:"attr,maintenance_time"`
	OriginalMaintenanceDuration string `jsonapi:"attr,maintenance_duration"`
	OriginalMaintenanceTz       string `jsonapi:"attr,maintenance_tz"`
	OriginalDesiredManifestSha  string `jsonapi:"attr,desired_manifest_sha"`
	OriginalCurrentManifestSha  string `jsonapi:"attr,current_manifest_sha"`
}

type UpdateInstanceType struct {
	InstanceType string `json:"instanceType"`
}

func (a *Client) UpdateInstanceType(sID string, eID string, updateData *UpdateInstanceType) error {
	req, err := json.Marshal(updateData)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("naut/project/%s/environment/%s/sspattributes", sID, eID)
	r, err := a.post(url, bytes.NewReader(req))
	if err != nil {
		return err
	}
	defer r.Close()

	return nil
}

func (a *Client) GetEnvironment(sID string, eID string) (*Environment, error) {
	url := fmt.Sprintf("naut/project/%s/environment/%s", sID, eID)
	r, err := a.get(url)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	env := new(Environment)
	err = jsonapi.UnmarshalPayload(r, env)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshaling environment: '%s'", err)
	}

	var days = map[string]time.Weekday{
		"Sunday":    time.Sunday,
		"Monday":    time.Monday,
		"Tuesday":   time.Tuesday,
		"Wednesday": time.Wednesday,
		"Thursday":  time.Thursday,
		"Friday":    time.Friday,
		"Saturday":  time.Saturday,
	}
	weekday, ok := days[env.OriginalMaintenanceDay]
	if ok {
		env.MaintenanceDay = weekday
		env.MaintenanceUnspecified = false
	} else {
		env.MaintenanceUnspecified = true
	}

	env.MaintenanceTime, err = parseSSTime(env.OriginalMaintenanceTime)
	if err != nil {
		env.MaintenanceUnspecified = true
	}

	dur, err := parseSSTime(env.OriginalMaintenanceDuration)
	if err != nil {
		env.MaintenanceUnspecified = true
	}
	env.MaintenanceDuration = time.Duration(dur.Hour()) * time.Hour
	env.MaintenanceDuration += time.Duration(dur.Minute()) * time.Minute
	env.MaintenanceDuration += time.Duration(dur.Second()) * time.Second

	env.MaintenanceTz, err = time.LoadLocation(env.OriginalMaintenanceTz)
	if err != nil {
		env.MaintenanceUnspecified = true
	}

	var usages = map[string]Usage{
		"Production":  UsageProduction,
		"UAT":         UsageUAT,
		"Test":        UsageTest,
		"Unspecified": UsageUnspecified,
	}
	usage, ok := usages[env.OriginalUsage]
	if ok {
		env.Usage = usage
	}

	env.CurrentManifestSha, _ = semver.Make(env.OriginalCurrentManifestSha)
	env.DesiredManifestSha, _ = semver.Make(env.OriginalDesiredManifestSha)

	return env, nil
}
