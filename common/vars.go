// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package common

var (
	ValidPrefixes        = []string{"ergo", "kubectl", "docker"}
	ListOutput           string
	PluginRunTypeCurl    = "curl"
	PluginRunTypeShell   = "shell"
	PluginRunTypeCompose = "compose"
	PluginRunTypeKube    = "kube"
	PluginRunTypeBin     = "bin"
)

type Protocol string

func (p Protocol) String() string {
	return string(p)
}

var TCPProtocol Protocol = "TCP"
var UDPProtocol Protocol = "UDP"
var ICMPProtocol Protocol = "ICMP"
var ALLProtocol Protocol = "ALL"

type FirewallRuleAction string

func (f FirewallRuleAction) String() string {
	return string(f)
}

var ACCEPTFirewallRuleAction FirewallRuleAction = "ACCEPT"
var DROPFirewallRuleAction FirewallRuleAction = "DROP"

var (
	Version       string
	BuildDate     string
	GitCommitHash string
)
