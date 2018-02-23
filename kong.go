package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type kongClient struct {
	httpClient   *http.Client
	kongAdminURL string
}

func newKongClient(kongHost, kongAdminPort string) *kongClient {
	return &kongClient{
		kongAdminURL: fmt.Sprintf("%s:%s", kongHost, kongAdminPort),
		httpClient:   &http.Client{},
	}
}

type upstreamResponse struct {
	Data []upstream `json:"data"`
}

type upstream struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (kc *kongClient) upstreams() ([]upstream, error) {
	upstreams := []upstream{}

	respBytes, err := kc.doRequest(http.MethodGet, "upstreams", nil)
	if err != nil {
		return upstreams, err
	}

	upstreamResponse := &upstreamResponse{}

	err = json.Unmarshal(respBytes, upstreamResponse)
	if err != nil {
		return upstreams, err
	}

	return upstreamResponse.Data, nil
}

type target struct {
	ID     string `json:"id"`
	URL    string `json:"target"`
	Weight string `json:"weight"`
}

type targetResponse struct {
	Data []target `json:"data"`
}

func (kc *kongClient) targetsFor(upstreamID string) ([]target, error) {
	targets := []target{}

	respBytes, err := kc.doRequest(http.MethodGet, fmt.Sprintf("upstreams/%s/targets", upstreamID), nil)
	if err != nil {
		return targets, err
	}

	targetResponse := &targetResponse{}

	err = json.Unmarshal(respBytes, targetResponse)
	if err != nil {
		return targets, err
	}

	return targetResponse.Data, nil
}

func (kc *kongClient) setTargetWeightFor(upstreamID, targetID, weight string) error {
	return nil
}

func (kc *kongClient) doRequest(method, path string, body []byte) ([]byte, error) {
	var respBytes []byte

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", kc.kongAdminURL, path), bytes.NewBuffer(body))
	if err != nil {
		return respBytes, err
	}

	response, err := kc.httpClient.Do(req)
	if err != nil {
		return respBytes, err
	}

	defer response.Body.Close()

	respBytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return respBytes, err
	}

	return respBytes, nil
}
