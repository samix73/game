package ecs

import (
	"context"
	"runtime/trace"
	"sync/atomic"
)

var nextID = atomic.Uint64{}

func NextID(ctx context.Context) ID {
	region := trace.StartRegion(ctx, "ecs.NextID")
	defer region.End()

	return ID(nextID.Add(1))
}

type ID uint64

const UndefinedID ID = 0
