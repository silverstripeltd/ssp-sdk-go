package ssp

import (
	"net/http"
	"testing"
)

func TestListDeployments(t *testing.T) {
	in := []*Deployment{
		{ID: 12},
		{ID: 13},
	}
	api, ts := newMockDashboard(in, http.StatusOK)
	defer ts.Close()

	out, err := api.ListDeployments("one", "prod", nil)
	if err != nil {
		t.Errorf("%s", err)
	}

	for i, _ := range in {
		if in[i].ID != out[i].ID {
			t.Error("Data returned is not matching the data sent")
		}
	}
}

func TestGetDeployment(t *testing.T) {
	changes := make(map[string]interface{})
	changes["Infrastructure"] = &DeploymentChange{From: "5", To: "6", Description: "Changed"}

	in := &Deployment{
		ID:              123,
		OriginalChanges: map[string]interface{}(changes),
		OriginalState:   "Completed",
	}
	api, ts := newMockDashboard(in, http.StatusOK)
	defer ts.Close()

	out, err := api.GetDeployment("one", "prod", "123")
	if err != nil {
		t.Errorf("%s", err)
	}
	if out.ID != 123 {
		t.Error("ID parsed incorrectly")
	}
	if out.State != StateCompleted {
		t.Error("State parsed incorrectly")
	}
	ch, ok := out.Changes["Infrastructure"]
	if !ok {
		t.Error("Changes parsed incorrectly")
	}
	if ch.From != "5" || ch.To != "6" || ch.Description != "Changed" {
		t.Error("Change content parsed incorrectly")
	}
}

func TestCreateDeployment(t *testing.T) {
	api, ts := newMockDashboard(&Deployment{}, http.StatusCreated)
	defer ts.Close()
	_, err := api.CreateDeployment("one", "prod", &CreateDeployment{})
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestApproveDeployment(t *testing.T) {
	api, ts := newMockDashboard(&Deployment{}, http.StatusOK)
	defer ts.Close()
	_, err := api.ApproveDeployment("one", "prod", &ApproveDeployment{})
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestStartDeployment(t *testing.T) {
	api, ts := newMockDashboard(&Deployment{}, http.StatusOK)
	defer ts.Close()
	_, err := api.StartDeployment("one", "prod", &StartDeployment{})
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestInvalidateDeployment(t *testing.T) {
	api, ts := newMockDashboard(&Deployment{}, http.StatusOK)
	defer ts.Close()
	_, err := api.InvalidateDeployment("one", "prod", &InvalidateDeployment{})
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteDeployment(t *testing.T) {
	api, ts := newMockDashboard(&Deployment{}, http.StatusNoContent)
	defer ts.Close()
	err := api.DeleteDeployment("one", "prod", 1)
	if err != nil {
		t.Errorf("%s", err)
	}
}
