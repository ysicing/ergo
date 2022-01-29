// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package ecs

type Action interface {
	Reset(id string) error
	List() error
}
