package repos

import (
	"encoding/json"
	"time"

	"github.com/wangle201210/githubapi"
	"github.com/wangle201210/githubapi/user"
)

const PullUrl = "pulls"

type Pull struct {
	URL                string        `json:"url"`
	ID                 int64         `json:"id"`
	NodeID             string        `json:"node_id"`
	HTMLURL            string        `json:"html_url"`
	DiffURL            string        `json:"diff_url"`
	PatchURL           string        `json:"patch_url"`
	IssueURL           string        `json:"issue_url"`
	Number             int64         `json:"number"`
	State              string        `json:"state"`
	Locked             bool          `json:"locked"`
	Title              string        `json:"title"`
	User               user.User     `json:"user"`
	Body               string        `json:"body"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
	ClosedAt           time.Time     `json:"closed_at"`
	MergedAt           time.Time     `json:"merged_at"`
	MergeCommitSha     string        `json:"merge_commit_sha"`
	Assignee           interface{}   `json:"assignee"`
	Assignees          []interface{} `json:"assignees"`
	RequestedReviewers []interface{} `json:"requested_reviewers"`
	RequestedTeams     []interface{} `json:"requested_teams"`
	Labels             []interface{} `json:"labels"`
	Milestone          interface{}   `json:"milestone"`
	Draft              bool          `json:"draft"`
	CommitsURL         string        `json:"commits_url"`
	ReviewCommentsURL  string        `json:"review_comments_url"`
	ReviewCommentURL   string        `json:"review_comment_url"`
	CommentsURL        string        `json:"comments_url"`
	StatusesURL        string        `json:"statuses_url"`
	Head               Head          `json:"head"`
	Base               Base          `json:"base"`
	Links              Links         `json:"_links"`
	AuthorAssociation  string        `json:"author_association"`
	ActiveLockReason   interface{}   `json:"active_lock_reason"`
}

type Base struct {
	Label string    `json:"label"`
	Ref   string    `json:"ref"`
	Sha   string    `json:"sha"`
	User  user.User `json:"user"`
	Repo  Repos     `json:"repos"`
}

type Links struct {
	Self           Self           `json:"self"`
	HTML           HTML           `json:"html"`
	Issue          Issue          `json:"issue"`
	Comments       Comments       `json:"comments"`
	ReviewComments ReviewComments `json:"review_comments"`
	ReviewComment  ReviewComment  `json:"review_comment"`
	Commits        Commits        `json:"commits"`
	Statuses       Statuses       `json:"statuses"`
}

type Self struct {
	Href string `json:"href"`
}
type HTML struct {
	Href string `json:"href"`
}
type Issue struct {
	Href string `json:"href"`
}
type Comments struct {
	Href string `json:"href"`
}
type ReviewComments struct {
	Href string `json:"href"`
}
type ReviewComment struct {
	Href string `json:"href"`
}
type Commits struct {
	Href string `json:"href"`
}
type Statuses struct {
	Href string `json:"href"`
}

type Head struct {
	Label string    `json:"label"`
	Ref   string    `json:"ref"`
	Sha   string    `json:"sha"`
	User  user.User `json:"user"`
	Repo  Repos     `json:"repos"`
}

type PullList []Pull

// List pull requests
func (re *Pkg) GetPullsList() (list PullList, err error) {
	url := re.Url(PullUrl)
	body, err := githubapi.HttpGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &list)
	return
}
