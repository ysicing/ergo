/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package cloud

import "context"

type EcsCloud interface {
	Create(ctx context.Context, option CreateOption) error
	Destroy(ctx context.Context, option DestroyOption) error
	Snapshot(ctx context.Context, option SnapshotOption) error
	Status(ctx context.Context, option StatusOption) error
	Halt(ctx context.Context, option HaltOption) error
	Up(ctx context.Context, option UpOption) error
	List(ctx context.Context, option ListOption) error
}

type CreateOption struct{}

type EcsOption struct{}

type DestroyOption struct {
	EcsOption
}

type SnapshotOption struct {
	EcsOption
}

type StatusOption struct {
	EcsOption
}

type HaltOption struct {
	EcsOption
}

type UpOption struct {
	EcsOption
}

type ListOption struct {
	EcsOption
}
