package aggendpoint

import (
	"context"
	"tolling/go-kit-example/aggsvc/aggservice"
	"tolling/types"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	AggregateEndpoint endpoint.Endpoint
	CalculateEndpoint endpoint.Endpoint
}

type AggregateRequest struct {
	Value float64 `json:"value"`
	OBUID int     `json:"ObuID"`
	Unix  int64   `json:"unix"`
}

type AggregateResponse struct {
	Err error `json:"error,omitempty"`
}

type CalculateRequest struct {
	OBUID int `json:"obu_id"`
}

type CalculateResponse struct {
	OBUID         int     `json:"obu_id"`
	TotalDistance float64 `json:"total_distance"`
	TotalAmount   float64 `json:"total_amount"`
	Error         error   `json:"error,omitempty"`
}

func (s *Set) Calculate(ctx context.Context, obuID int) (*types.Invoice, error) {
	resp, err := s.CalculateEndpoint(ctx, CalculateRequest{
		OBUID: obuID,
	})
	if err != nil {
		return nil, err
	}
	response := resp.(CalculateResponse)
	return &types.Invoice{
		OBUID:         response.OBUID,
		TotalDistance: response.TotalDistance,
		TotalAmount:   response.TotalAmount,
	}, response.Error
}

func (s *Set) Aggregate(ctx context.Context, dist types.Distance) error {
	_, err := s.AggregateEndpoint(ctx, AggregateRequest{
		Value: dist.Value,
		OBUID: dist.OBUID,
		Unix:  dist.Unix,
	})
	return err
}

func MakeAggregateEndpoint(svc aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AggregateRequest)
		err := svc.Aggregate(ctx, types.Distance{
			Value: req.Value,
			OBUID: req.OBUID,
			Unix:  req.Unix,
		})
		return AggregateResponse{Err: err}, nil
	}
}

func MakeCalculateEndpoint(svc aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CalculateRequest)
		inv, err := svc.Calculate(ctx, req.OBUID)
		if err != nil {
			return CalculateResponse{Error: err}, nil
		}
		return CalculateResponse{
			OBUID:         inv.OBUID,
			TotalDistance: inv.TotalDistance,
			TotalAmount:   inv.TotalAmount,
		}, nil
	}
}
