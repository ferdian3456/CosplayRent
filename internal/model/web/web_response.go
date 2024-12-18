package web

type WebResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type AppResponse struct {
	AppVersion string `json:"app_version"`
}
