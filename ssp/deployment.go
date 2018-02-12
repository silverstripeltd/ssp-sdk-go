package ssp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/google/jsonapi"
	"github.com/mitchellh/mapstructure"
	"io"
	"reflect"
	"time"
)

const (
	IgnoreConfigChangeOption = "ignore_config_changes"
	ForceFullOption          = "force_full"
)

const (
	DeploymentTypeFull     = "full"
	DeploymentTypeCodeOnly = "code-only"
)

type State string

const (
	StateNew       = "New"
	StateSubmitted = "Submitted"
	StateInvalid   = "Invalid"
	StateApproved  = "Approved"
	StateRejected  = "Rejected"
	StateQueued    = "Queued"
	StateDeploying = "Deploying"
	StateAborting  = "Aborting"
	StateCompleted = "Completed"
	StateFailed    = "Failed"
	StateDeleted   = "Deleted"
)

type DeploymentFilter struct {
	State           State  `url:"state"`
	LastEditedFrom  int64  `url:"lastedited_from_unix"`
	DateStartedFrom int64  `url:"datestarted_from_unix"`
	DateStartedTo   int64  `url:"datestarted_to_unix"`
	DeployerEmail   string `url:"deployer_email"`
	Title           string `url:"title"`
	Summary         string `url:"summary"`
}

type Deployment struct {
	ID              int       `jsonapi:"primary,deployments"`
	DateCreated     time.Time `jsonapi:"attr,date_created_unix"`
	DateStarted     time.Time `jsonapi:"attr,date_started_unix"`
	DateRequested   time.Time `jsonapi:"attr,date_requested_unix"`
	DateUpdated     time.Time `jsonapi:"attr,date_updated_unix"`
	ScheduleStart   time.Time `jsonapi:"attr,schedule_start_unix"`
	ScheduleEnd     time.Time `jsonapi:"attr,schedule_end_unix"`
	Title           string    `jsonapi:"attr,title"`
	Summary         string    `jsonapi:"attr,summary"`
	RefType         string    `jsonapi:"attr,ref_type"`
	RefName         string    `jsonapi:"attr,ref_name"`
	RejectedReason  string    `jsonapi:"attr,rejected_reason"`
	Tags            []string  `jsonapi:"attr,tags"`
	DeploymentType  string    `jsonapi:"attr,deployment_type"`
	SHA             string    `jsonapi:"attr,sha"`
	ShortSHA        string    `jsonapi:"attr,short_sha"`
	IsCurrentBuild  bool      `jsonapi:"attr,is_current_build"`
	Changes         map[string]DeploymentChange
	OriginalChanges map[string]interface{} `jsonapi:"attr,changes"`
	State           State
	OriginalState   string       `jsonapi:"attr,state"`
	Environment     *Environment `jsonapi:"relation,environment"`
	Stack           *Stack       `jsonapi:"relation,stack"`
}

type DeploymentChange struct {
	From        string
	To          string
	Description string
}

type CreateDeployment struct {
	Ref            string   `json:"ref"`
	RefType        string   `json:"ref_type"`
	Title          string   `json:"title"`
	Summary        string   `json:"summary"`
	Options        []string `json:"options"`
	ScheduleStart  int64    `json:"schedule_start_unix"`
	ScheduleEnd    int64    `json:"schedule_end_unix"`
	Bypass         bool     `json:"bypass"`
	BypassAndStart bool     `json:"bypass_and_start"`
	Locked         bool     `json:"locked"`
}

type StartDeployment struct {
	ID int `json:"id"`
}

type InvalidateDeployment struct {
	ID int `json:"id"`
}

type ApproveDeployment struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

func (a *Client) GetDeploymentCurrent(sID string, eID string) (*Deployment, error) {
	return a.GetDeployment(sID, eID, "current")
}

func (a *Client) GetDeploymentCurrentFull(sID string, eID string) (*Deployment, error) {
	return a.GetDeployment(sID, eID, "currentfull")
}

func (a *Client) ListDeployments(sID string, eID string, filter *DeploymentFilter) ([]*Deployment, error) {
	url := fmt.Sprintf("naut/project/%s/environment/%s/deploys", sID, eID)

	// todo: move further back to api.Get
	q, err := query.Values(filter)
	if err != nil {
		return nil, err
	}
	url += "?" + q.Encode()

	r, err := a.get(url)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	data, err := jsonapi.UnmarshalManyPayload(r, reflect.TypeOf(new(Deployment)))
	if err != nil {
		return nil, fmt.Errorf("failed unmarshaling deployments: '%s'", err)
	}

	deployments := make([]*Deployment, len(data))
	for i, deployment := range data {
		deployments[i] = deployment.(*Deployment)
		err = postProcessDeployment(deployments[i])
		if err != nil {
			return nil, fmt.Errorf("failed post-processing deployments: '%s'", err)
		}
	}

	return deployments, nil
}

func (a *Client) GetDeployment(sID string, eID string, dID string) (*Deployment, error) {
	url := fmt.Sprintf("naut/project/%s/environment/%s/deploys/%s", sID, eID, dID)
	resp, err := a.get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	d, unmarshalErr := responseToDeployment(resp)
	if unmarshalErr != nil {
		return nil, fmt.Errorf("failed unmarshaling deployment: '%s'", unmarshalErr)
	}

	return d, nil
}

func (a *Client) CreateDeployment(sID string, eID string, cd *CreateDeployment) (*Deployment, error) {
	req, err := json.Marshal(cd)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("naut/project/%s/environment/%s/deploys", sID, eID)
	resp, err := a.post(url, bytes.NewReader(req))
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	d, unmarshalErr := responseToDeployment(resp)
	if unmarshalErr != nil {
		return nil, fmt.Errorf("failed unmarshaling deployment: '%s'", unmarshalErr)
	}

	return d, nil
}

func (a *Client) ApproveDeployment(sID string, eID string, ad *ApproveDeployment) (*Deployment, error) {
	req, err := json.Marshal(ad)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("naut/project/%s/environment/%s/approvals/approve", sID, eID)
	resp, readErr := a.post(url, bytes.NewReader(req))
	if readErr != nil {
		return nil, readErr
	}
	defer resp.Close()

	d, unmarshalErr := responseToDeployment(resp)
	if unmarshalErr != nil {
		return nil, fmt.Errorf("failed unmarshaling deployment: '%s'", unmarshalErr)
	}

	return d, nil
}

func (a *Client) StartDeployment(sID string, eID string, sd *StartDeployment) (*Deployment, error) {
	req, err := json.Marshal(sd)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("naut/project/%s/environment/%s/deploys/start", sID, eID)
	resp, readErr := a.post(url, bytes.NewReader(req))
	if readErr != nil {
		return nil, readErr
	}
	defer resp.Close()

	d, unmarshalErr := responseToDeployment(resp)
	if unmarshalErr != nil {
		return nil, fmt.Errorf("failed unmarshaling deployment: '%s'", unmarshalErr)
	}

	return d, nil
}

func (a *Client) InvalidateDeployment(sID string, eID string, id *InvalidateDeployment) (*Deployment, error) {
	req, err := json.Marshal(id)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("naut/project/%s/environment/%s/deploys/invalidate", sID, eID)
	resp, readErr := a.post(url, bytes.NewReader(req))
	if readErr != nil {
		return nil, readErr
	}
	defer resp.Close()

	d, unmarshalErr := responseToDeployment(resp)
	if unmarshalErr != nil {
		return nil, fmt.Errorf("failed unmarshaling deployment: '%s'", unmarshalErr)
	}

	return d, nil
}

func (a *Client) DeleteDeployment(sID string, eID string, dID int) error {
	url := fmt.Sprintf("naut/project/%s/environment/%s/deploys/%d", sID, eID, dID)
	resp, readErr := a.delete(url, nil)
	if readErr != nil {
		return readErr
	}
	defer resp.Close()

	return nil
}

func responseToDeployment(r io.Reader) (*Deployment, error) {
	d := &Deployment{}
	err := jsonapi.UnmarshalPayload(r, d)
	if err != nil {
		return nil, err
	}

	err = postProcessDeployment(d)
	if err != nil {
		return nil, err
	}

	return d, err
}

func postProcessDeployment(d *Deployment) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &d.Changes,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	err = decoder.Decode(d.OriginalChanges)
	if err != nil {
		return err
	}

	var states = map[string]State{
		"New":       StateNew,
		"Submitted": StateSubmitted,
		"Invalid":   StateInvalid,
		"Approved":  StateApproved,
		"Rejected":  StateRejected,
		"Queued":    StateQueued,
		"Deploying": StateDeploying,
		"Aborting":  StateAborting,
		"Completed": StateCompleted,
		"Failed":    StateFailed,
		"Deleted":   StateDeleted,
	}
	state, ok := states[d.OriginalState]
	if ok {
		d.State = state
	}

	return nil
}
