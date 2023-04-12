// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package types

type SSH struct {
	User     string `json:"user,omitempty"`
	Passwd   string `json:"passwd,omitempty"`
	PkName   string `json:"pkName,omitempty"`
	PkData   string `json:"pkData,omitempty"`
	Pk       string `json:"pk,omitempty"`
	PkPasswd string `json:"pkPasswd,omitempty"`
}
