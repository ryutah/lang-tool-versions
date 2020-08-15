package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

func loadPythonVersions() ([]string, error) {
	c := colly.NewCollector()

	var versions []string
	c.OnHTML("div.row:nth-child(2) > ol:nth-child(4) > li > span.release-number", func(e *colly.HTMLElement) {
		versions = append(versions, strings.TrimLeft(e.Text, "Python "))
	})

	if err := c.Visit("https://www.python.org/downloads/"); err != nil {
		return nil, errors.WithStack(err)
	}

	return versions, nil
}

func loadGoVersions() ([]string, error) {
	tags, err := sendGitHubListTagRequestForGetAll("golang", "go")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return extractAsVersions(tags, `^go(\d+\.\d+(?:\.\d+)?)$`), nil
}

func loadNodeVersions() ([]string, error) {
	tags, err := sendGitHubListTagRequestForGetAll("nodejs", "node")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return extractAsVersions(tags, `^v(\d+\.\d+\.\d+)$`), nil
}

func loadDockerComposeVersions() ([]string, error) {
	tags, err := sendGitHubListTagRequestForGetAll("docker", "compose")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return extractAsVersions(tags, `^(\d+\.\d+\.\d+)$`), nil
}

func loadGotaskTask() ([]string, error) {
	tags, err := sendGitHubListTagRequestForGetAll("go-task", "task")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return extractAsVersions(tags, `^v(\d+\.\d+\.\d+)$`), nil
}

func loadGoogleCloudSDK() ([]string, error) {
	tags, err := sendGitHubListTagRequestForGetAll("GoogleCloudPlatform", "cloud-sdk-docker")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return extractAsVersions(tags, `^(\d+\.\d+\.\d+)$`), nil
}

func loadGolangciLint() ([]string, error) {
	tags, err := sendGitHubListTagRequestForGetAll("golangci", "golangci-lint")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return extractAsVersions(tags, `^v(\d+\.\d+\.\d+)$`), nil
}

func loadRustVersions() ([]string, error) {
	tags, err := sendGitHubListTagRequestForGetAll("rust-lang", "rust")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return extractAsVersions(tags, `^(\d+\.\d+\.\d+)$`), nil
}

func extractAsVersions(tags []tagResponse, extractRule string) []string {
	r := regexp.MustCompile(extractRule)
	versions := make([]string, 0, len(tags))
	for _, p := range tags {
		if matches := r.FindStringSubmatch(p.Name); len(matches) > 0 {
			versions = append(versions, matches[1])
		}
	}
	return versions
}

func sendGitHubListTagRequestForGetAll(org, repo string) (allTags []tagResponse, err error) {
	for i := 1; ; i++ {
		tags, e := sendGitHubListTagRequest(org, repo, i)
		if err = e; err != nil {
			err = errors.WithStack(err)
			return
		}
		if len(tags) == 0 {
			break
		}
		allTags = append(allTags, tags...)
	}
	return
}

// tagResponse resuponse definition of GitHubAPI /repos/{owner}/{repo}/tags(https://docs.github.com/en/rest/reference/repos#list-repository-tags)
// this is mimimam definition for use.
type tagResponse struct {
	Name string `json:"name"`
}

func sendGitHubListTagRequest(org, repo string, page int) (tags []tagResponse, err error) {
	uri := url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   fmt.Sprintf("repos/%s/%s/tags", org, repo),
	}
	query := url.Values{
		"per_page": []string{"100"},
	}

	if page > 1 {
		query["page"] = []string{strconv.Itoa(page)}
	}
	uri.RawQuery = query.Encode()

	req, e := http.NewRequest(http.MethodGet, uri.String(), nil)
	if e != nil {
		err = e
		return
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		err = e
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return
	}
	return
}
