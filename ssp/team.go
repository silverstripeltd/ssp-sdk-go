package ssp

import (
	"fmt"
	"reflect"

	"github.com/google/jsonapi"
)

type User struct {
	Username  string `jsonapi:"attr,username"`
	Email     string `jsonapi:"attr,email"`
	FirstName string `jsonapi:"attr,first_name"`
	Surname   string `jsonapi:"attr,surname"`
	Role      string `jsonapi:"attr,role"`
	Stack     *Stack `jsonapi:"relation,stack"`
}

func (a *Client) ListTeam(sID string) ([]*User, error) {
	url := fmt.Sprintf("naut/project/%s/team", sID)
	r, err := a.get(url)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	items, err := jsonapi.UnmarshalManyPayload(r, reflect.TypeOf(new(User)))
	if err != nil {
		return nil, fmt.Errorf("failed unmarshaling users: '%s'", err)
	}

	users := make([]*User, len(items))
	for i, item := range items {
		users[i] = item.(*User)
	}

	return users, nil
}
