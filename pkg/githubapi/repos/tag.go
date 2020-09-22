package repos

import (
	"encoding/json"
	"errors"

	"github.com/wangle201210/githubapi"
)

const TagUrl = "tags"

type Commit struct {
	Sha string `json:"sha"`
	Url string `json:"url"`
}

type Tag struct {
	Name       string `json:"name"`
	ZipballUrl string `json:"zipball_url"`
	TarballUrl string `json:"tarball_url"`
	Commit     Commit `json:"commit"`
}

type TagsList []Tag

func (re *Pkg) GetTagsList() (list TagsList, err error) {
	url := re.Url(TagUrl)
	body, err := githubapi.HttpGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &list)
	return
}

func (re *Pkg) LastTag() (tag Tag, err error) {
	list, err := re.GetTagsList()
	if err != nil {
		return
	}
	if len(list) < 1 {
		return tag, errors.New("there is no tag")
	}
	return list[0], nil
}
