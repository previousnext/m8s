package stash

type Group struct {
	Name       string `json:"name"`
	Deleteable bool   `json:"deleteable,omitempty"`
}

type PagedGroup struct {
	Size       int      `json:"size"`
	Limit      int      `json:"limit"`
	Start      int      `json:"start"`
	IsLastPage bool     `json:"isLastPage"`
	Values     []*Group `json:"values"`
}

type PagedUsers struct {
	Size       int     `json:"size"`
	Limit      int     `json:"limit"`
	Start      int     `json:"start"`
	IsLastPage bool    `json:"isLastPage"`
	Values     []*User `json:"values"`
}

type User struct {
	Active                      bool                           `json:"active"`
	Slug                        string                         `json:"slug"`
	Deleteable                  bool                           `json:"deleatable"`
	ID                          int                            `json:"id"`
	DisplayName                 string                         `json:"displayName"`
	MutableDetails              bool                           `json:"mutableDetails"`
	Link                        map[string]string              `json:"link"`
	Name                        string                         `json:"name"`
	Type                        string                         `json:"type"`
	DirectoryName               string                         `json:"directoryName"`
	LastAuthenticationTimestamp int                            `json:"lastAuthenticationTimestamp"`
	Links                       map[string][]map[string]string `json:"links"`
	EmailAddress                string                         `json:"emailAddress"`
	MutableGroups               bool                           `json:"mutableGroups"`
}

type PagedRepos struct {
	Size       int     `json:"size"`
	Limit      int     `json:"limit"`
	Start      int     `json:"start"`
	IsLastPage bool    `json:"isLastPage"`
	Values     []*Repo `json:"values"`
}

type Project struct {
	Key    string                         `json:"key"`
	ID     int                            `json:"id"`
	Name   string                         `json:"name"`
	Public bool                           `json:"public"`
	Type   string                         `json:"type"`
	Link   map[string]string              `json:"link"`
	Links  map[string][]map[string]string `json:"links"`
}

type Repo struct {
	Project       Project                        `json:"project"`
	Public        bool                           `json:"public"`
	Link          map[string]string              `json:"link"`
	CloneURL      string                         `json:"cloneUrl"`
	Slug          string                         `json:"slug"`
	SCMID         string                         `json:"scmId"`
	State         string                         `json:"state"`
	Forkable      bool                           `json:"forkable"`
	ID            int                            `json:"id"`
	Name          string                         `json:"name"`
	StatusMessage string                         `json:"statusMessage"`
	Links         map[string][]map[string]string `json:"links"`
}

type Branch struct {
	ID              string `json:"id"`
	DisplayID       string `json:"displayId"`
	LatestChangeset string `json:"latestChangeset"`
	LatestCommit    string `json:"latestCommit"`
	IsDefault       bool   `json:"isDefault"`
}
