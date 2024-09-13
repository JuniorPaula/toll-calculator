package aggservice

import (
	"context"
	"tolling/types"
)

type Middleware func(Service) Service

type loggerMiddleware struct {
	next Service
}

func newLoggerMiddleware() Middleware {
	return func(next Service) Service {
		return loggerMiddleware{next: next}
	}
}

func (mw loggerMiddleware) Aggregate(_ context.Context, dist types.Distance) error {
	return nil
}

func (mw loggerMiddleware) Calculate(_ context.Context, i int) (*types.Invoice, error) {
	return nil, nil
}

type instrumentationMiddleware struct {
	next Service
}

func newInstrumentationMiddleware() Middleware {
	return func(next Service) Service {
		return instrumentationMiddleware{next: next}
	}
}

func (mw instrumentationMiddleware) Aggregate(_ context.Context, dist types.Distance) error {
	return nil
}

func (mw instrumentationMiddleware) Calculate(_ context.Context, i int) (*types.Invoice, error) {
	return nil, nil
}
