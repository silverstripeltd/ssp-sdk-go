package ssp

import (
	"net/http"
	"reflect"
	"testing"
)

func TestListModules(t *testing.T) {
	in := []*ModuleData{
		{ID: "silverstripe/framework", Name: "silverstripe/framework", Version: "3.7.1"},
		{ID: "silverstripe/cms", Name: "silverstripe/cms", Version: "3.7.1"},
	}

	api, ts := newMockDashboard(in, http.StatusOK)
	defer ts.Close()

	out, err := api.ListModules("example", "production")
	if err != nil {
		t.Errorf("%s", err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Error("Data returned is not matching the data sent")
	}
}
