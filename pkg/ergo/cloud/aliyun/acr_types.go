// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package aliyun

// NamespacesRes 命名空间接口返回元数据
type NamespacesRes struct {
	Data struct {
		Namespaces []Namespace `json:"namespaces"`
	} `json:"data"`
}

// Namespace 命令空间
type Namespace struct {
	Namespace       string `json:"namespace"`
	NamespaceStatus string `json:"namespaceStatus"`
	AuthorizeType   string `json:"authorizeType"`
}

// ReposRes 仓库接口返回元数据
type ReposRes struct {
	Data struct {
		Total    int    `json:"total"`
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
		Repos    []Repo `json:"repos"`
	} `json:"data"`
}

// Repo 仓库元数据
type Repo struct {
	// Summary        string `json:"summary"`
	// RepoDomainList struct {
	// Internal string `json:"internal"`
	// Public   string `json:"public"`
	// Vpc      string `json:"vpc"`
	// } `json:"repoDomainList"`
	RepoAuthorizeType string `json:"repoAuthorizeType"`
	Downloads         int    `json:"downloads"`
	// Logo              string `json:"logo"`
	// Stars             int    `json:"stars"`
	RepoType      string `json:"repoType"`
	RepoNamespace string `json:"repoNamespace"`
	RepoName      string `json:"repoName"`
	RepoStatus    string `json:"repoStatus"`
	RepoID        int    `json:"repoId"`
	// RepoBuildType     string `json:"repoBuildType"`
	RegionID       string `json:"regionId"`
	RepoOriginType string `json:"repoOriginType"`
	GmtCreate      int64  `json:"gmtCreate"`
	GmtModified    int64  `json:"gmtModified"`
	LastTag        string `json:"lasttag"`
	ImageUpdate    int64  `json:"imageUpdate"`
}

// TagsRes 元数据
type TagsRes struct {
	Data struct {
		Total    int   `json:"total"`
		Page     int   `json:"page"`
		Tags     []Tag `json:"tags"`
		PageSize int   `json:"pageSize"`
	} `json:"data"`
}

// Tag 镜像标签
type Tag struct {
	Status      string `json:"status"`
	Digest      string `json:"digest"`
	ImageCreate int64  `json:"imageCreate"`
	ImageID     string `json:"imageId"`
	ImageUpdate int64  `json:"imageUpdate"`
	Tag         string `json:"tag"`
	ImageSize   int    `json:"imageSize"`
}
