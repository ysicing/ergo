package repos

import (
	"encoding/json"
	"time"

	"github.com/wangle201210/githubapi"
	"github.com/wangle201210/githubapi/user"
)

const (
	IssuesUrl  = "issues"
	CommentUrl = "issues/comments"
)

type Issues struct {
	URL                   string        `json:"url"`
	RepositoryURL         string        `json:"repository_url"`
	LabelsURL             string        `json:"labels_url"`
	CommentsURL           string        `json:"comments_url"`
	EventsURL             string        `json:"events_url"`
	HTMLURL               string        `json:"html_url"`
	ID                    int           `json:"id"`
	NodeID                string        `json:"node_id"`
	Number                int           `json:"number"`
	Title                 string        `json:"title"`
	User                  user.User     `json:"user"`
	Labels                []interface{} `json:"labels"`
	State                 string        `json:"state"`
	Locked                bool          `json:"locked"`
	Assignee              interface{}   `json:"assignee"`
	Assignees             []interface{} `json:"assignees"`
	Milestone             interface{}   `json:"milestone"`
	Comments              int           `json:"comments"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
	ClosedAt              interface{}   `json:"closed_at"`
	AuthorAssociation     string        `json:"author_association"`
	ActiveLockReason      interface{}   `json:"active_lock_reason"`
	Body                  string        `json:"body"`
	PerformedViaGithubApp interface{}   `json:"performed_via_github_app"`
}

type Comment struct {
	URL                   string      `json:"url"`
	HTMLURL               string      `json:"html_url"`
	IssueURL              string      `json:"issue_url"`
	ID                    int         `json:"id"`
	NodeID                string      `json:"node_id"`
	User                  user.User   `json:"user"`
	CreatedAt             time.Time   `json:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at"`
	AuthorAssociation     string      `json:"author_association"`
	Body                  string      `json:"body"`
	PerformedViaGithubApp interface{} `json:"performed_via_github_app"`
}

func (re *Pkg) GetIssuesList() (list []Issues, err error) {
	url := re.Url(IssuesUrl)
	body, err := githubapi.HttpGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &list)
	return
}

func (re *Pkg) GetCommentsList() (list []Comment, err error) {
	url := re.Url(CommentUrl)
	body, err := githubapi.HttpGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &list)
	return
}
