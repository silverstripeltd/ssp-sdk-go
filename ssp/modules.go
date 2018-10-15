package ssp

import (
	"fmt"
	"reflect"

	"github.com/google/jsonapi"
)

type ModuleData struct {
	ID      string `jsonapi:"primary,module"`
	Name    string `jsonapi:"attr,name"`
	Version string `jsonapi:"attr,version"`
}

func (a *Client) ListModules(sID string, eID string) ([]*ModuleData, error) {
	url := fmt.Sprintf("naut/project/%s/environment/%s/modules", sID, eID)
	r, err := a.get(url)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	items, err := jsonapi.UnmarshalManyPayload(r, reflect.TypeOf(new(ModuleData)))
	if err != nil {
		return nil, fmt.Errorf("failed unmarshaling modules: '%s'", err)
	}

	modules := make([]*ModuleData, len(items))
	for i, item := range items {
		modules[i] = item.(*ModuleData)
	}

	return modules, nil
}
