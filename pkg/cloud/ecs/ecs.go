// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package ecs

type ECSAction interface {
	Reset() error
	List() error
}