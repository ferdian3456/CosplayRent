package rajaongkir

import (
	"bytes"
	"context"
	"cosplayrent/helper"
	"cosplayrent/model/web/rajaongkir"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const (
	rajaOngkirURL = "https://api.rajaongkir.com/starter"
)

type RajaOngkirServiceImpl struct {
	validate *validator.Validate
}

func NewRajaOngkirService(validate *validator.Validate) *RajaOngkirServiceImpl {
	return &RajaOngkirServiceImpl{
		validate: validate,
	}
}

func (service *RajaOngkirServiceImpl) FindProvince(ctx context.Context) (rajaongkir.RajaOngkirProvinceResponse, error) {
	err := godotenv.Load("../.env")
	helper.PanicIfError(err)

	rajaongkirAPIKEY := os.Getenv("RAJAONGKIR_SERVER_KEY")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/province", rajaOngkirURL), nil)
	if err != nil {
		return rajaongkir.RajaOngkirProvinceResponse{}, errors.New("failed to create request to RajaOngkir")
	}

	req.Header.Set("key", rajaongkirAPIKEY)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return rajaongkir.RajaOngkirProvinceResponse{}, errors.New("failed to reach RajaOngkir API")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rajaongkir.RajaOngkirProvinceResponse{}, errors.New("failed to read response body from RajaOngkir response")
	}

	var rajaongkirProvinceResponse rajaongkir.RajaOngkirProvinceResponse
	err = json.Unmarshal(body, &rajaongkirProvinceResponse)
	if err != nil {
		return rajaongkir.RajaOngkirProvinceResponse{}, errors.New("failed to unmarshal response body from RajaOngkir response")
	}

	return rajaongkirProvinceResponse, nil
}

func (service *RajaOngkirServiceImpl) FindCity(ctx context.Context, provinceID string) (rajaongkir.RajaOngkirCityResponse, error) {
	err := godotenv.Load("../.env")
	helper.PanicIfError(err)

	rajaongkirAPIKEY := os.Getenv("RAJAONGKIR_SERVER_KEY")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/city?province=%s", rajaOngkirURL, provinceID), nil)
	if err != nil {
		return rajaongkir.RajaOngkirCityResponse{}, errors.New("failed to create request to RajaOngkir")
	}

	req.Header.Set("key", rajaongkirAPIKEY)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return rajaongkir.RajaOngkirCityResponse{}, errors.New("failed to reach RajaOngkir API")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rajaongkir.RajaOngkirCityResponse{}, errors.New("failed to read response body from RajaOngkir response")
	}

	var rajaOngkirCityResponse rajaongkir.RajaOngkirCityResponse
	err = json.Unmarshal(body, &rajaOngkirCityResponse)
	if err != nil {
		return rajaongkir.RajaOngkirCityResponse{}, errors.New("failed to unmarshal response body from RajaOngkir response")
	}

	return rajaOngkirCityResponse, nil
}

func (service *RajaOngkirServiceImpl) CheckShippment(ctx context.Context, shipmentRequest rajaongkir.RajaOngkirSendShipmentRequest) (rajaongkir.RajaOngkirShipmentResponse, error) {
	sendRequest := url.Values{}
	finalWeight := strconv.Itoa(shipmentRequest.Weight)
	sendRequest.Set("origin", shipmentRequest.Origin)
	sendRequest.Set("destination", shipmentRequest.Destination)
	sendRequest.Set("weight", finalWeight)
	sendRequest.Set("courier", shipmentRequest.Courier)

	req, err := http.NewRequest("POST", "https://api.rajaongkir.com/starter/cost", bytes.NewBufferString(sendRequest.Encode()))
	if err != nil {
		return rajaongkir.RajaOngkirShipmentResponse{}, errors.New("failed to create request to RajaOngkir")
	}

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)
	rajaongkirAPIKEY := os.Getenv("RAJAONGKIR_SERVER_KEY")

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("key", rajaongkirAPIKEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return rajaongkir.RajaOngkirShipmentResponse{}, errors.New("failed to reach RajaOngkir API")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rajaongkir.RajaOngkirShipmentResponse{}, errors.New("failed to read response body from RajaOngkir response")
	}

	var RajaOngkirShipmentResponse rajaongkir.RajaOngkirShipmentResponse
	err = json.Unmarshal(body, &RajaOngkirShipmentResponse)
	if err != nil {
		return rajaongkir.RajaOngkirShipmentResponse{}, errors.New("failed to unmarshal response body from RajaOngkir response")
	}

	return RajaOngkirShipmentResponse, nil
}
