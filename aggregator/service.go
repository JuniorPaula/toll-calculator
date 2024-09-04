package main

import (
	"tolling/types"

	"github.com/sirupsen/logrus"
)

const basePrice = 3.15

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoicerAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	logrus.WithFields(logrus.Fields{
		"obuid":    distance.OBUID,
		"distance": distance.Value,
		"unix":     distance.Unix,
	}).Info("aggregating distance")
	return i.store.Insert(distance)
}

func (i *InvoiceAggregator) CalculateInvoice(obuId int) (*types.Invoice, error) {
	dist, err := i.store.Get(obuId)
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
