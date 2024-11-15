package rajaongkir

type RajaOngkirCityResponse struct {
	Rajaongkir struct {
		Query  interface{} `json:"query"`
		Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
		Results []struct {
			CityID   string `json:"city_id"`
			CityName string `json:"city_name"`
		} `json:"results"`
	} `json:"rajaongkir"`
}
