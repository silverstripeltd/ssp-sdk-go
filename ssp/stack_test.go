package ssp

import (
	"net/http"
	"reflect"
	"testing"
)

func TestListStacks(t *testing.T) {
	in := []*Stack{
		{ID: "one"},
		{ID: "two"},
	}
	api, ts := newMockDashboard(in, http.StatusOK)
	defer ts.Close()

	out, err := api.ListStacks()
	if err != nil {
		t.Errorf("%s", err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Error("Data returned is not matching the data sent")
	}
}
