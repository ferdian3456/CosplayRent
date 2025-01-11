package usecase

import (
	"bytes"
	"context"
	"cosplayrent/internal/model/web/rajaongkir"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-playground/validator"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	rajaOngkirURL = "https://api.rajaongkir.com/starter"
)

type RajaOngkirUsecase struct {
	Cache    *memcache.Client
	Validate *validator.Validate
	Log      *zerolog.Logger
	Config   *koanf.Koanf
}

func NewRajaOngkirUsecase(cache *memcache.Client, validate *validator.Validate, zerolog *zerolog.Logger, koanf *koanf.Koanf) *RajaOngkirUsecase {
	return &RajaOngkirUsecase{
		Cache:    cache,
		Validate: validate,
		Log:      zerolog,
		Config:   koanf,
	}
}

func (usecase *RajaOngkirUsecase) FindProvince(ctx context.Context) rajaongkir.RajaOngkirProvinceResponse {
	cachedData, err := usecase.Cache.Get("RajaOngkirProvinceCache")
	if err == nil && cachedData != nil {
		usecase.Log.Info().Msg(("Hit province cache"))

		var cachedResponse rajaongkir.RajaOngkirProvinceResponse

		err := json.Unmarshal(cachedData.Value, &cachedResponse)
		if err != nil {
			respErr := errors.New("failed to unmarshal cached data")
			usecase.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return cachedResponse
	}

	rajaongkirAPIKEY := usecase.Config.String("RAJAONGKIR_SERVER_KEY")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/province", rajaOngkirURL), nil)
	if err != nil {
		respErr := errors.New("failed to create request to RajaOngkir")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	req.Header.Set("key", rajaongkirAPIKEY)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respErr := errors.New("failed to reach RajaOngkir API")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		respErr := errors.New("failed to read response body from RajaOngkir response")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	var rajaongkirProvinceResponse rajaongkir.RajaOngkirProvinceResponse
	err = json.Unmarshal(body, &rajaongkirProvinceResponse)
	if err != nil {
		respErr := errors.New("failed to unmarshal response body from RajaOngkir response")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	cacheData, err := json.Marshal(rajaongkirProvinceResponse)
	if err == nil {
		err = usecase.Cache.Set(&memcache.Item{
			Key:   "RajaOngkirProvinceCache",
			Value: cacheData,
		})
		if err != nil {
			respErr := errors.New("failed to set cache")
			usecase.Log.Panic().Err(err).Msg(respErr.Error())
		} else {
			usecase.Log.Info().Msg(("Success to create cache for RajaOngkirProvince's response"))
		}
	}

	return rajaongkirProvinceResponse
}

func (usecase *RajaOngkirUsecase) FindCity(ctx context.Context, provinceID string) rajaongkir.RajaOngkirCityResponse {
	cacheKey := fmt.Sprintf("RajaOngkirCityCache_%s", provinceID)

	cachedData, err := usecase.Cache.Get(cacheKey)
	if err == nil && cachedData != nil {
		usecase.Log.Info().Msg(("Hit city cache for province:" + provinceID))

		var cachedResponse rajaongkir.RajaOngkirCityResponse

		err := json.Unmarshal(cachedData.Value, &cachedResponse)
		if err != nil {
			respErr := errors.New("failed to unmarshal cached data")
			usecase.Log.Panic().Err(err).Msg(respErr.Error())
		}

		return cachedResponse
	}

	rajaongkirAPIKEY := usecase.Config.String("RAJAONGKIR_SERVER_KEY")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/city?province=%s", rajaOngkirURL, provinceID), nil)
	if err != nil {
		respErr := errors.New("failed to create request to RajaOngkir")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	req.Header.Set("key", rajaongkirAPIKEY)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		respErr := errors.New("failed to reach RajaOngkir API")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		respErr := errors.New("failed to read response body from RajaOngkir response")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	var rajaOngkirCityResponse rajaongkir.RajaOngkirCityResponse
	err = json.Unmarshal(body, &rajaOngkirCityResponse)
	if err != nil {
		respErr := errors.New("failed to unmarshal response body from RajaOngkir response")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	cacheData, err := json.Marshal(rajaOngkirCityResponse)
	if err == nil {
		err = usecase.Cache.Set(&memcache.Item{
			Key:   cacheKey,
			Value: cacheData,
		})
		if err != nil {
			respErr := errors.New("failed to set city")
			usecase.Log.Panic().Err(err).Msg(respErr.Error())
		} else {
			usecase.Log.Info().Msg(("Success to create cache for RajaOngkirCities's response"))
		}
	}

	return rajaOngkirCityResponse
}

func (usecase *RajaOngkirUsecase) CheckShippment(ctx context.Context, shipmentRequest rajaongkir.RajaOngkirSendShipmentRequest) (rajaongkir.RajaOngkirShipmentResponse, error) {
	err := usecase.Validate.Struct(shipmentRequest)
	if err != nil {
		respErr := errors.New("invalid request body")
		usecase.Log.Warn().Err(respErr).Msg(err.Error())
		return rajaongkir.RajaOngkirShipmentResponse{}, respErr
	}

	sendRequest := url.Values{}
	finalWeight := strconv.Itoa(shipmentRequest.Weight)
	sendRequest.Set("origin", shipmentRequest.Origin)
	sendRequest.Set("destination", shipmentRequest.Destination)
	sendRequest.Set("weight", finalWeight)
	sendRequest.Set("courier", shipmentRequest.Courier)

	req, err := http.NewRequest("POST", "https://api.rajaongkir.com/starter/cost", bytes.NewBufferString(sendRequest.Encode()))
	if err != nil {
		respErr := errors.New("failed to create request to RajaOngkir")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	rajaongkirAPIKEY := usecase.Config.String("RAJAONGKIR_SERVER_KEY")

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("key", rajaongkirAPIKEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		respErr := errors.New("failed to reach RajaOngkir API")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		respErr := errors.New("failed to read response body from RajaOngkir response")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	var RajaOngkirShipmentResponse rajaongkir.RajaOngkirShipmentResponse
	err = json.Unmarshal(body, &RajaOngkirShipmentResponse)
	if err != nil {
		respErr := errors.New("failed to unmarshal response body from RajaOngkir response")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	return RajaOngkirShipmentResponse, nil
}
