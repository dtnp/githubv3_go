package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	inputFile    = "repos.txt"
	organization = "pantheon-systems"
)

type GithubRepo struct {
	ID            int64     `csv:"id"`
	Name          string    `csv:"repo_name"`
	Owner         string    `csv:"owner"`
	DefaultBranch string    `csv:"default_branch"`
	Language      string    `csv:"language"`
	IsPrivate     bool      `csv:"is_private"`
	IsArchived    bool      `csv:"is_archived"`
	HasIssues     bool      `csv:"has_issues"`
	ForkCount     int       `csv:"fork_count"`
	CreatedAt     time.Time `csv:"created_at"`
	UpdatedAt     time.Time `csv:"updated_at"`
	URL           string    `csv:"repo_url"`
	Description   string    `csv:"repo_description"`
	CommitDetails
}

type CommitDetails struct {
	SHA                 string `csv:"last_commit_sha"`
	DaysSinceLastCommit int    `csv:"days_since_last_commit"`
	AuthorName          string `csv:"commit_author_name"`
	AuthorEmail         string `csv:"commit_author_email"`
	AuthorAlias         string `csv:"commit_author_alias"`
	MergeMessage        string `csv:"commit_merge_message"`
	CommitURL           string `csv:"commit_message_url"`
}



func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("'GITHUB_TOKEN' is required")
	}

	ctx := context.Background()
	client := GetGithubClient(ctx, token)

	Run(ctx, client)
}

func Run(ctx context.Context, client *github.Client) {

	repos := GetAllRepos(ctx, client)

	// to download file inside downloads folder
	file, err := ioutil.TempFile(".", "repos.*.csv")
	if err != nil {
		fmt.Println("Err: ", err)
	}
	defer file.Close()

	err = gocsv.MarshalFile(&repos, file)
	if err != nil {
		fmt.Println("Err: ", err)
	}
}

func GetGithubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func GetAllRepos(ctx context.Context, client *github.Client) []GithubRepo {

	var allRepos []GithubRepo
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	counter := 0
	now := time.Now().UTC()

	// get all pages of results - loop until there are no more results
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, organization, opt)
		if err != nil {
			log.Fatal(err)
		}

		for _, val := range repos {

			var grd GithubRepo
			var cd CommitDetails

			// Days since last commit -- only for unarchived things
			updatedAt := val.GetUpdatedAt().UTC()

			if !val.GetArchived() {
				// Newer date subtracted by older date
				days := now.Sub(updatedAt).Hours() / 24
				cd = GetCommit(ctx, client, val.GetName(), days)
			}

			grd = GithubRepo{
				ID:            val.GetID(),
				Name:          val.GetName(),
				Owner:         val.GetOwner().GetLogin(),
				DefaultBranch: val.GetDefaultBranch(),
				IsPrivate:     val.GetPrivate(),
				IsArchived:    val.GetArchived(),
				Language:      val.GetLanguage(),
				Description:   val.GetDescription(),
				URL:           val.GetURL(),
				HasIssues:     val.GetHasIssues(),
				ForkCount:     val.GetForksCount(),
				CreatedAt:     val.GetCreatedAt().UTC(),
				UpdatedAt:     updatedAt,

				CommitDetails: cd,
			}

			allRepos = append(allRepos, grd)

			counter += 1
			fmt.Printf("Repo #%d: %s (%s) -- Days: %d, %s\n", counter, grd.Name, grd.DefaultBranch, grd.DaysSinceLastCommit, grd.AuthorAlias)
			//if counter > 10 {
			//	os.Exit(1)
			//}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos
}

func GetCommit(ctx context.Context, client *github.Client, repoName string, daysSinceLastCommit float64) CommitDetails {
	clo := &github.CommitsListOptions{
		ListOptions: github.ListOptions{Page: 1, PerPage: 1},
	}
	repoCommits, commitResp, err := client.Repositories.ListCommits(ctx, organization, repoName, clo)
	if err != nil {
		fmt.Printf("commit request fail (%d): %s\n", commitResp.StatusCode, err.Error())
	}

	var cd CommitDetails
	for _, commit := range repoCommits {
		cd = CommitDetails{
			SHA:                 commit.GetSHA(),
			AuthorName:          commit.GetCommit().GetAuthor().GetName(),
			AuthorEmail:         commit.GetCommit().GetAuthor().GetEmail(),
			AuthorAlias:         commit.GetAuthor().GetLogin(),
			MergeMessage:        commit.GetCommit().GetMessage(),
			CommitURL:           commit.GetCommit().GetURL(),
			DaysSinceLastCommit: int(daysSinceLastCommit),
		}
	}
	return cd
}
