package ssp

import (
	"errors"
	"fmt"
	"github.com/blang/semver"
	"github.com/google/jsonapi"
	"reflect"
	"time"
)

type ManifestRelease struct {
	ID          string `jsonapi:"primary,manifestreleases"`
	Sha         semver.Version
	Released    time.Time `jsonapi:"attr,released_unix"`
	OriginalSha string    `jsonapi:"attr,sha"`
}

func (a *Client) ListManifestReleases() ([]*ManifestRelease, error) {
	r, err := a.get("naut/manifestreleases")
	if err != nil {
		return nil, err
	}
	defer r.Close()

	items, err := jsonapi.UnmarshalManyPayload(r, reflect.TypeOf(new(ManifestRelease)))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed unmarshaling manifest releases: '%s'", err))
	}

	releases := make([]*ManifestRelease, len(items))
	for i, item := range items {
		releases[i] = item.(*ManifestRelease)
		releases[i].Sha, _ = semver.Make(releases[i].OriginalSha)
	}

	return releases, nil
}
