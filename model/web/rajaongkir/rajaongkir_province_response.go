package rajaongkir

type RajaOngkirProvinceResponse struct {
	Rajaongkir struct {
		Query  interface{} `json:"query"`
		Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
		Results []struct {
			ProvinceID string `json:"province_id"`
			Province   string `json:"province"`
		} `json:"results"`
	} `json:"rajaongkir"`
}
