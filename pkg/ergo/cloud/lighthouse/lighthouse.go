// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package lighthouse

type Action interface {
	Reset(id string) error
	BindKey(id string) error
	List() error
}
