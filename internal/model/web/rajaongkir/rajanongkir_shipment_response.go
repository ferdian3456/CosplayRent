package rajaongkir

type RajaOngkirShipmentResponse struct {
	Rajaongkir struct {
		Query  interface{} `json:"query"`
		Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
		Results []struct {
			Code  string `json:"code"`
			Name  string `json:"name"`
			Costs []struct {
				Usecase     string `json:"Usecase"`
				Description string `json:"description"`
				Cost        []struct {
					Value int    `json:"value"`
					ETD   string `json:"etd"`
					Note  string `json:"note"`
				} `json:"cost"`
			} `json:"costs"`
		} `json:"results"`
	} `json:"rajaongkir"`
}
