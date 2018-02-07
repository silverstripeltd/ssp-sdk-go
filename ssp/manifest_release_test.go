package ssp

import (
	"net/http"
	"reflect"
	"testing"
)

func TestListManifestReleases(t *testing.T) {
	in := []*ManifestRelease{
		{ID: "one"},
		{ID: "two"},
	}
	api, ts := newMockDashboard(in, http.StatusOK)
	defer ts.Close()

	out, err := api.ListManifestReleases()
	if err != nil {
		t.Errorf("%s", err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Error("Data returned is not matching the data sent")
	}
}
