// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package compose

func ComposeDeploy(service string) {
	svc := NewCompose(service, ComposeConfig{
		Hosts:       Hosts,
		DeployLocal: DeployLocal,
		Service:     service,
	})
	svc.Check()
	svc.Write()
	svc.Up()
}
