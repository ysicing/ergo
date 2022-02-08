// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cloud

import (
	"strings"

	"github.com/ergoapi/util/exid"
	"github.com/manifoldco/promptui"
	"github.com/ysicing/ergo/pkg/config"
)

func addProvider() config.Provider {
	var provider, aid, akey, region string
	pprompt := promptui.Prompt{
		Label: "provider",
	}
	provider, _ = pprompt.Run()
	aidprompt := promptui.Prompt{
		Label: "aid",
	}
	aid, _ = aidprompt.Run()
	akeyprompt := promptui.Prompt{
		Label: "akey",
	}
	akey, _ = akeyprompt.Run()
	regionprompt := promptui.Prompt{
		Label: "region",
	}

	cfg := config.Provider{
		UUID:     exid.GenUUID(),
		Provider: strings.Trim(provider, " "),
		Secrets: config.Secrets{
			AID:  strings.Trim(aid, " "),
			AKey: strings.Trim(akey, ""),
		},
	}

	region, _ = regionprompt.Run()
	if len(region) != 0 {
		cfg.Regions = []string{strings.Trim(region, " ")}
	}
	return cfg
}
