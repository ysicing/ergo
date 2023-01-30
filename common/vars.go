// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

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
