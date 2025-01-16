package library

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Fn func(context.Context) error

type Group struct {
	eg *errgroup.Group
}

func NewGroup() *Group {
	return &Group{
		eg: &errgroup.Group{},
	}
}

func (g *Group) Go(ctx context.Context, funcs ...Fn) {
	for _, f := range funcs {
		fn := f
		g.eg.Go(func() error {
			return fn(ctx)
		})
	}
}

func (g *Group) Wait() error {
	return g.eg.Wait()
}

type GroupWithContext struct {
	eg  *errgroup.Group
	ctx context.Context
}

func NewGroupWithContext(ctx context.Context, funcs ...Fn) *GroupWithContext {
	eg, ctx := errgroup.WithContext(ctx)

	gwc := &GroupWithContext{
		eg:  eg,
		ctx: ctx,
	}

	gwc.Go(funcs...)

	return gwc
}

func (gwc *GroupWithContext) Go(funcs ...Fn) {
	for _, f := range funcs {
		fn := f
		gwc.eg.Go(func() error {
			return fn(gwc.ctx)
		})
	}
}

func (gwc *GroupWithContext) Wait() error {
	return gwc.eg.Wait()
}
