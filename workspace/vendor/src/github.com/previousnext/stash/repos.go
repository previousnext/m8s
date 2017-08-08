package stash

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RepoService interface {
	GetRepos(string) (*PagedRepos, error)
	GetRepo(string, string) (*Repo, error)
	CreateRepo(string, string, string) (*Repo, error)
	CreateBranch(string, string, string, string) (*Branch, error)
}

type repoService struct {
	*Client
}

func (r *repoService) GetRepos(projectKey string) (*PagedRepos, error) {
	req, err := r.createReq(
		"GET", fmt.Sprintf("/rest/api/1.0/projects/%s/repos", projectKey), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var b PagedRepos

	if resp.StatusCode != http.StatusOK {
		return &b, errors.New("Recieved a non 200 response")
	}

	err = json.NewDecoder(resp.Body).Decode(&b)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (r *repoService) GetRepo(projectKey, repoKey string) (*Repo, error) {
	req, err := r.createReq(
		"GET", fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s", projectKey, repoKey), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var b Repo

	if resp.StatusCode != http.StatusOK {
		return &b, errors.New("Recieved a non 200 response")
	}

	err = json.NewDecoder(resp.Body).Decode(&b)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return &b, nil
}

// where scmId = git or hg
func (r *repoService) CreateRepo(projectKey, name, scmID string) (*Repo, error) {
	req, err := r.createReq(
		"POST",
		fmt.Sprintf("/rest/api/1.0/projects/%s/repos", projectKey),
		map[string]string{
			"name":  name,
			"scmId": scmID,
		},
	)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var b Repo

	if resp.StatusCode != http.StatusCreated {
		return &b, errors.New("Recieved a non 200 response")
	}

	err = json.NewDecoder(resp.Body).Decode(&b)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (r *repoService) CreateBranch(projKey, repo, name, startRef string) (*Branch, error) {
	req, err := r.createReq(
		"POST",
		fmt.Sprintf("/rest/branch-utils/1.0/projects/%s/repos/%s/branches",
			projKey, repo),
		map[string]string{
			"name":       name,
			"startPoint": startRef,
		},
	)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var b Branch

	if resp.StatusCode != http.StatusCreated {
		return &b, errors.New("Recieved a non 201 response")
	}

	err = json.NewDecoder(resp.Body).Decode(&b)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return &b, nil
}
