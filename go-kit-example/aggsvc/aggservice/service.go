package aggservice

import (
	"context"
	"tolling/types"
)

const basePrice = 3.15

type Service interface {
	Aggregate(context.Context, types.Distance) error
	Calculate(context.Context, int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type BasicService struct {
	store Storer
}

func newBasicService(store Storer) Service {
	return &BasicService{store: store}
}

func (svc *BasicService) Aggregate(_ context.Context, dist types.Distance) error {
	return svc.store.Insert(dist)
}

func (svc *BasicService) Calculate(_ context.Context, obuId int) (*types.Invoice, error) {
	dist, err := svc.store.Get(obuId)
	if err != nil {
		return nil, err
	}

	inv := &types.Invoice{
		OBUID:         obuId,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}

	return inv, nil
}

func NewAggregatorService() Service {
	var svc Service
	{
		svc = newBasicService(makeStore())
		svc = newLoggerMiddleware()(svc)
		svc = newInstrumentationMiddleware()(svc)
	}
	return svc
}
