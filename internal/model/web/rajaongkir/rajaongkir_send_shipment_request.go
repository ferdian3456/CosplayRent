package rajaongkir

type RajaOngkirSendShipmentRequest struct {
	Origin      string `validate:"required" json:"origin"`
	Destination string `validate:"required" json:"destination"`
	Weight      int    `validate:"required" json:"weight"`
	Courier     string `validate:"required" json:"courier"`
}
