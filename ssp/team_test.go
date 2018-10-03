package ssp

import (
	"net/http"
	"reflect"
	"testing"
)

func TestListTeam(t *testing.T) {
	in := []*User{
		{Username: "janebloggs"},
		{Username: "joebloggs"},
	}
	api, ts := newMockDashboard(in, http.StatusOK)
	defer ts.Close()

	out, err := api.ListTeam("asdasd")
	if err != nil {
		t.Errorf("%s", err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Error("Data returned is not matching the data sent")
	}
}
