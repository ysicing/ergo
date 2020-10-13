package repos_test

import (
	"testing"

	"github.com/wangle201210/githubapi/repos"
)

var pkg = repos.Pkg{"beego", "bee"}

func TestUrl(t *testing.T) {
	url := pkg.Url("pulls")
	expect := "https://api.github.com/repos/beego/bee/" + "pulls"
	if url != expect {
		t.Errorf("make url err, expect '%s',got '%s'", expect, url)
	}
}

func TestGetRepo(t *testing.T) {
	re, err := pkg.GetRepos()
	if err != nil {
		t.Errorf("get repos err: %s", err)
		return
	}
	if re.ID == 0 {
		t.Errorf("repos is %+v", re)
		t.Errorf("repos must have id but got null, ")
	}
}

func TestGetPullList(t *testing.T) {
	list, err := pkg.GetPullsList()
	if err != nil {
		t.Errorf("get pull list err: %s", err)
		return
	}
	if len(list) > 0 && list[0].ID == 0 {
		t.Errorf("pull list is %+v", list)
		t.Errorf("pull must have id but got null")
	}
}

func TestGetTagsList(t *testing.T) {
	list, err := pkg.GetTagsList()
	if err != nil {
		t.Errorf("get tag list err: %s", err)
		return
	}
	if len(list) > 0 && list[0].Name == "" {
		t.Errorf("tag list is %+v", list)
		t.Errorf("tag must have name but got null")
	}
}

func TestLastTag(t *testing.T) {
	tag, err := pkg.LastTag()
	if err != nil {
		t.Errorf("get pull list err: %s", err)
		return
	}
	if tag.Name == "" {
		t.Errorf("tag list is %+v", tag)
		t.Errorf("tag must have name but got null")
	}
}

func TestGetIssuesList(t *testing.T) {
	list, err := pkg.GetIssuesList()
	if err != nil {
		t.Errorf("get issue list err: %s", err)
		return
	}
	if len(list) > 0 && list[0].ID == 0 {
		t.Errorf("issue list is %+v", list)
		t.Errorf("issue must have name but got null")
	}
}

func TestGetCommentList(t *testing.T) {
	list, err := pkg.GetCommentsList()
	if err != nil {
		t.Errorf("get comment list err: %s", err)
		return
	}
	if len(list) > 0 && list[0].ID == 0 {
		t.Errorf("comment list is %+v", list)
		t.Errorf("comment must have id but got null")
	}
}
