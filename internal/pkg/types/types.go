package types

type SSH struct {
	User     string `json:"user,omitempty"`
	Passwd   string `json:"passwd,omitempty"`
	PkName   string `json:"pkName,omitempty"`
	PkData   string `json:"pkData,omitempty"`
	Pk       string `json:"pk,omitempty"`
	PkPasswd string `json:"pkPasswd,omitempty"`
}
