package marketplace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ProvisionRequest struct {
	QuickNodeId       string   `json:"quicknode-id"`
	EndpointId        string   `json:"endpoint-id"`
	Chain             string   `json:"chain"`
	Network           string   `json:"network"`
	Plan              string   `json:"plan"`
	WSSURL            string   `json:"wss-url"`
	HTTPURL           string   `json:"http-url"`
	Referers          []string `json:"referers"`
	ContractAddresses []string `json:"contract_addresses"`
}

type ProvisionResponse struct {
	Status       string `json:"status"`
	DashboardURL string `json:"dashboard-url"`
	AccessURL    string `json:"access-url"`
}

type UpdateRequest struct {
	QuickNodeId       string   `json:"quicknode-id"`
	EndpointId        string   `json:"endpoint-id"`
	Chain             string   `json:"chain"`
	Network           string   `json:"network"`
	Plan              string   `json:"plan"`
	WSSURL            string   `json:"wss-url"`
	HTTPURL           string   `json:"http-url"`
	Referers          []string `json:"referers"`
	ContractAddresses []string `json:"contract_addresses"`
}

type UpdateResponse struct {
	Status string `json:"status"`
}

type DeactivateRequest struct {
	QuickNodeId  string `json:"quicknode-id"`
	EndpointId   string `json:"endpoint-id"`
	Chain        string `json:"chain"`
	Network      string `json:"network"`
	DeactivateAt string `json:"deactivate-at"`
}

type DeactivateResponse struct {
	Status string `json:"status"`
}

type DeprovisionRequest struct {
	QuickNodeId   string `json:"quicknode-id"`
	EndpointId    string `json:"endpoint-id"`
	DeprovisionAt string `json:"deprovision-at"`
}

type DeprovisionResponse struct {
	Status string `json:"status"`
}

func Provision(url string, payload ProvisionRequest) (ProvisionResponse, error) {
	client := &http.Client{}

	// Convert the payload to JSON
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, payloadBuf)
	if err != nil {
		fmt.Println(err)
		return ProvisionResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ProvisionResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ProvisionResponse{}, fmt.Errorf("HTTP Request failed with status code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ProvisionResponse{}, err
	}

	var response ProvisionResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		fmt.Println(err)
		return ProvisionResponse{}, err
	}

	return response, nil
}
