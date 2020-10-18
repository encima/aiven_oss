package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v32/github"
	"github.com/resilva87/stringmetric"
	"golang.org/x/oauth2"
)

var c chan bool

var licenses []string

func init() {
	c = make(chan bool)
	file, err := os.Open("./licenses.csv")
	handleErr(err)
	defer file.Close()
	csvr := csv.NewReader(file)
	osi, err := csvr.ReadAll()
	handleErr(err)
	for _, row := range osi {
		licenses = append(licenses, row[0])
	}
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Action struct {
	User    string
	Commits int
}

type Meta struct {
	Members []Action
	Commits int
	Stars   int
	Public  bool
	License string
}

func checkLicense(repo string) (license string, found bool) {
	if repo == "Other" || len(repo) == 0 {
		fmt.Println("No License")
		return "", false
	}

	match := 0.0
	matchingLicense := "None"
	for _, license := range licenses {
		res := stringmetric.RatcliffObershelpMetric(repo, license)
		if res > match {
			match = res
			matchingLicense = license
		}
	}
	if matchingLicense != "None" {
		return matchingLicense, true
	}
	return "", false
}

func main() {
	ctx := context.Background()
	fmt.Println(os.Getenv("AIVEN_OSS_GITHUB_PAT"))
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("AIVEN_OSS_GITHUB_PAT")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	contribs := make(map[string]Meta)
	// list public repositories for org "aiven"
	opt := &github.RepositoryListByOrgOptions{Type: "public"}
	repos, _, err := client.Repositories.ListByOrg(context.Background(), "aiven", opt)
	for _, repo := range repos {
		if val, ok := contribs[*repo.Name]; ok {
			fmt.Println(val)
		} else {
			contribs[*repo.Name] = Meta{Stars: *repo.StargazersCount}
		}
	}
	members, _, err := client.Organizations.ListMembers(context.Background(), "aiven", &github.ListMembersOptions{
		PublicOnly: false,
		Filter:     "",
		Role:       "",
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: 100,
		},
	})
	handleErr(err)
	for _, member := range members {
		repos, _, err := client.Repositories.List(
			ctx,
			*member.Login,
			&github.RepositoryListOptions{Visibility: "public"},
		)
		handleErr(err)
		for _, repo := range repos {
			commits, _, err := client.Repositories.ListCommits(ctx, *repo.Owner.Login, *repo.Name, &github.CommitsListOptions{Author: *member.Login})
			handleErr(err)
			if len(commits) > 0 {
				license := repo.License.GetName()
				match, public := checkLicense(license)
				if public {
					if val, ok := contribs[*repo.Name]; ok {
						val.Members = append(val.Members, Action{User: *member.Login, Commits: len(commits)})
					} else {
						contribs[*repo.Name] = Meta{Stars: *repo.StargazersCount, Members: []Action{Action{User: *member.Login, Commits: len(commits)}}, Public: true, License: match}
					}
				}
			}
		}

	}
	fmt.Println(contribs)
}
