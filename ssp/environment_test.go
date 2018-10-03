package ssp

import (
	"net/http"
	"testing"
	"time"
)

func TestGetEnvironment(t *testing.T) {
	in := &Environment{
		ID:                          "one",
		OriginalUsage:               "Production",
		OriginalMaintenanceDay:      "Monday",
		OriginalMaintenanceTz:       "Europe/Vatican",
		OriginalMaintenanceDuration: "2:34",
		OriginalCurrentManifestSha:  "1.2.3",
		OriginalDesiredManifestSha:  "2.3.0",
		PHPVersion:                  "5.6",
	}
	api, ts := newMockDashboard(in, http.StatusOK)
	defer ts.Close()

	out, err := api.GetEnvironment("one", "prod")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if out.ID != "one" {
		t.Error("ID parsed incorrectly")
	}
	if out.Usage != UsageProduction {
		t.Error("Usage parsed incorrectly")
	}
	if out.MaintenanceDay != time.Monday {
		t.Error("MaintenanceDay parsed incorrectly")
	}
	if out.MaintenanceTz.String() != "Europe/Vatican" {
		t.Error("MaintenanceTz parsed incorrectly")
	}
	if out.MaintenanceDuration.Minutes() != 2*60+34 {
		t.Error("MaintenanceDuration parsed incorrectly")
	}
	if out.CurrentManifestSha.String() != "1.2.3" {
		t.Error("CurrentManifestSha parsed incorrectly")
	}
	if out.DesiredManifestSha.String() != "2.3.0" {
		t.Error("DesiredManifestSha parsed incorrectly")
	}
	if out.PHPVersion != in.PHPVersion {
		t.Errorf("PHPVersion parsed incorrectly, expected '%s', got '%s'", in.PHPVersion, out.PHPVersion)
	}
}

func TestUpdateInstanceType(t *testing.T) {
	api, ts := newMockDashboard(&Environment{}, http.StatusOK)
	defer ts.Close()
	s := &UpdateInstanceType{
		InstanceType: "t2.nano",
	}
	err := api.UpdateInstanceType("one", "prod", s)
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetEnvironmentUnspecifiedMaintenance(t *testing.T) {
	in := &Environment{
		ID: "one",
		OriginalMaintenanceDay:      "Unspecified",
		OriginalMaintenanceTz:       "Europe/Vatican",
		OriginalMaintenanceDuration: "2:34",
	}
	api, ts := newMockDashboard(in, http.StatusOK)
	defer ts.Close()

	out, err := api.GetEnvironment("one", "prod")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if !out.MaintenanceUnspecified {
		t.Error("Maintenance should be unspecified")
	}
}
