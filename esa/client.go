package esa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	Endpoint = "api.esa.io"
)

type Post struct {
	Number   int    `json:"number"`
	Name     string `json:"name"`
	Category string `json:"category"`
	URL      string `json:"url"`
}

type Posts struct {
	Posts []*Post `json:"posts"`
}

type Client struct {
	team string
}

func NewClient(team string) *Client {
	return &Client{
		team: team,
	}
}

func (cli *Client) Get(token string, name string, category string) (*Post, error) {
	req, err := cli.buildRequest(token, name, category)

	if err != nil {
		return nil, err
	}

	httpCli := &http.Client{}
	res, err := httpCli.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", res.Status, body)
	}

	posts := &Posts{}
	err = json.Unmarshal(body, &posts)

	if err != nil {
		return nil, err
	}

	for _, p := range posts.Posts {
		if p.Name == name && p.Category == category {
			return p, nil
		}
	}

	return nil, nil
}

func (cli *Client) buildRequest(token string, name string, category string) (*http.Request, error) {
	category = strings.TrimPrefix(strings.TrimSuffix(category, "/"), "/")
	url := fmt.Sprintf("https://%s/v1/teams/%s/posts", Endpoint, cli.team)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("page", "1")
	query.Add("per_page", "1")
	query.Add("q", fmt.Sprintf(`title:"%s" on:"%s"`, name, category))

	req.URL.RawQuery = query.Encode()
	req.Header.Add("Authorization", "Bearer "+token)

	return req, nil
}
