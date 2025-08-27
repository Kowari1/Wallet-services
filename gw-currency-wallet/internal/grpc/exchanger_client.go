package grpcClient

import (
	"context"
	"gw-currency-wallet/gw-exchanger/proto/exchange"

	"google.golang.org/grpc"
)

type ExchangeClient interface {
	GetRate(ctx context.Context, from, to string) (float64, error)
	GetAllRates(ctx context.Context) ([]*exchange.GetRateResponse, error)
}

type ExchangeAdapter struct {
	client exchange.ExchangeServiceClient
}

func NewExchangeAdapter(conn *grpc.ClientConn) *ExchangeAdapter {
	return &ExchangeAdapter{
		client: exchange.NewExchangeServiceClient(conn),
	}
}

func (a *ExchangeAdapter) GetRate(ctx context.Context, from, to string) (float64, error) {
	resp, err := a.client.GetRate(ctx, &exchange.GetRateRequest{
		FromCurrency: from,
		ToCurrency:   to,
	})
	if err != nil {
		return 0, err
	}

	return resp.Rate, nil
}

func (a *ExchangeAdapter) GetAllRates(ctx context.Context) ([]*exchange.GetRateResponse, error) {
	resp, err := a.client.GetAllRates(ctx, &exchange.GetAllRatesRequest{})
	if err != nil {
		return nil, err
	}

	return resp.Rates, nil
}
