package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"integration-service/internal/model"
	"net/http"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func CallOrderService(orderRequest model.OrderRequest, url string) (model.OrderResponse, error) {
	var orderResponse model.OrderResponse

	body, err := json.Marshal(orderRequest)
	if err != nil {
		return orderResponse, err
	}

	resp, err := httpClient.Post(url+"/orders", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return orderResponse, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return orderResponse, errors.New("failed to call order service")
	}

	err = json.NewDecoder(resp.Body).Decode(&orderResponse)
	if err != nil {
		return orderResponse, err
	}

	return orderResponse, nil
}

func CallGetOrderService(orderID, url string) (model.OrderResponse, error) {
	var orderResponse model.OrderResponse

	resp, err := httpClient.Get(fmt.Sprintf("%s/orders/%s", url, orderID))
	if err != nil {
		return orderResponse, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return orderResponse, errors.New("failed to get order")
	}

	err = json.NewDecoder(resp.Body).Decode(&orderResponse)
	if err != nil {
		return orderResponse, err
	}

	return orderResponse, nil
}

func CallShippingService(shippingRequest model.ShippingRequest, url string) (model.ShippingResponse, error) {
	var shippingResponse model.ShippingResponse

	body, err := json.Marshal(shippingRequest)
	if err != nil {
		return shippingResponse, err
	}

	resp, err := httpClient.Post(url+"/shipping", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return shippingResponse, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return shippingResponse, errors.New("failed to call shipping service")
	}

	err = json.NewDecoder(resp.Body).Decode(&shippingResponse)
	if err != nil {
		return shippingResponse, err
	}

	return shippingResponse, nil
}

func CallInventoryService(inventoryRequest model.InventoryRequest, url string) (model.InventoryResponse, error) {
	var inventoryResponse model.InventoryResponse

	body, err := json.Marshal(inventoryRequest)
	if err != nil {
		return inventoryResponse, err
	}

	resp, err := httpClient.Post(url+"/inventory", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return inventoryResponse, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return inventoryResponse, errors.New("failed to call order service")
	}

	err = json.NewDecoder(resp.Body).Decode(&inventoryResponse)
	if err != nil {
		return inventoryResponse, err
	}

	return inventoryResponse, nil
}

func CallUsersService(userRequest model.UserRequest, url string) (model.UserResponse, error) {
	var userResponse model.UserResponse

	body, err := json.Marshal(userRequest)
	if err != nil {
		return userResponse, err
	}

	resp, err := httpClient.Post(url+"/users", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return userResponse, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return userResponse, errors.New("failed to call users service")
	}

	err = json.NewDecoder(resp.Body).Decode(&userResponse)
	if err != nil {
		return userResponse, err
	}

	return userResponse, nil
}

func CallSalesReportService(request model.SalesReportRequest, url string) ([]model.SalesReportResponse, error) {
	var response []model.SalesReportResponse

	body, err := json.Marshal(request)
	if err != nil {
		return response, err
	}

	resp, err := httpClient.Post(url+"/reports/sales", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, errors.New("failed to call sales report service")
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, err
	}

	return response, nil
}
