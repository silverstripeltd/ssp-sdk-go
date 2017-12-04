package ssp

import (
	"errors"
	"fmt"
	"github.com/google/jsonapi"
	"reflect"
	"time"
)

type Stack struct {
	ID           string         `jsonapi:"primary,stacks"`
	Name         string         `jsonapi:"attr,name"`
	Title        string         `jsonapi:"attr,title"`
	Created      time.Time      `jsonapi:"attr,created_unix"`
	Environments []*Environment `jsonapi:"relation,environments"`
}

func (a *Client) ListStacks() ([]*Stack, error) {
	r, err := a.get("naut/projects")
	if err != nil {
		return nil, err
	}
	defer r.Close()

	items, err := jsonapi.UnmarshalManyPayload(r, reflect.TypeOf(new(Stack)))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed unmarshaling stacks: '%s'", err))
	}

	stacks := make([]*Stack, len(items))
	for i, item := range items {
		stacks[i] = item.(*Stack)
	}

	return stacks, nil
}
