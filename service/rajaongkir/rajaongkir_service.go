package rajaongkir

import (
	"context"
	"cosplayrent/model/web/rajaongkir"
)

type RajaOngkirService interface {
	FindProvince(ctx context.Context) (rajaongkir.RajaOngkirProvinceResponse, error)
	FindCity(ctx context.Context, provinceID string) (rajaongkir.RajaOngkirCityResponse, error)
	CheckShippment(ctx context.Context, shipmentRequest rajaongkir.RajaOngkirSendShipmentRequest) (rajaongkir.RajaOngkirShipmentResponse, error)
}
