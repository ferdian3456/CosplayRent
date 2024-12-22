package web

type WebResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type OrderStatusResponse struct {
	Status_payment bool `json:"status_payment"`
}

type AppResponse struct {
	AppVersion string `json:"app_version"`
}
