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
	DeactivateAt int64  `json:"deactivate-at"`
}

type DeactivateResponse struct {
	Status string `json:"status"`
}

type DeprovisionRequest struct {
	QuickNodeId   string `json:"quicknode-id"`
	EndpointId    string `json:"endpoint-id"`
	DeprovisionAt int64  `json:"deprovision-at"`
}

type DeprovisionResponse struct {
	Status string `json:"status"`
}

func RequiresBasicAuth(url string, httpMethod string) (bool, error) {
	client := &http.Client{}

	payload := ProvisionRequest{}

	// Convert the payload to JSON
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)

	// Create the HTTP request
	req, err := http.NewRequest(httpMethod, url, payloadBuf)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-QN-TESTING", "true")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return true, nil
	} else {
		return false, nil
	}
}

func Provision(url string, payload ProvisionRequest, basicAuth string) (ProvisionResponse, error) {
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
	req.Header.Add("Authorization", "Basic "+basicAuth)
	req.Header.Add("X-QN-TESTING", "true")

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

func Update(url string, payload UpdateRequest, basicAuth string) (UpdateResponse, error) {
	client := &http.Client{}

	// Convert the payload to JSON
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)

	// Create the HTTP request
	req, err := http.NewRequest("PUT", url, payloadBuf)
	if err != nil {
		fmt.Println(err)
		return UpdateResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth)
	req.Header.Add("X-QN-TESTING", "true")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return UpdateResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return UpdateResponse{}, fmt.Errorf("HTTP Request failed with status code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return UpdateResponse{}, err
	}

	var response UpdateResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		fmt.Println(err)
		return UpdateResponse{}, err
	}

	return response, nil
}

func Deactivate(url string, payload DeactivateRequest, basicAuth string) (DeactivateResponse, error) {
	client := &http.Client{}

	// Convert the payload to JSON
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)

	// Create the HTTP request
	req, err := http.NewRequest("DELETE", url, payloadBuf)
	if err != nil {
		fmt.Println(err)
		return DeactivateResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth)
	req.Header.Add("X-QN-TESTING", "true")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return DeactivateResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return DeactivateResponse{}, fmt.Errorf("HTTP Request failed with status code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return DeactivateResponse{}, err
	}

	var response DeactivateResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		fmt.Println(err)
		return DeactivateResponse{}, err
	}

	return response, nil
}

func Deprovision(url string, payload DeprovisionRequest, basicAuth string) (DeprovisionResponse, error) {
	client := &http.Client{}

	// Convert the payload to JSON
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)

	// Create the HTTP request
	req, err := http.NewRequest("DELETE", url, payloadBuf)
	if err != nil {
		fmt.Println(err)
		return DeprovisionResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth)
	req.Header.Add("X-QN-TESTING", "true")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return DeprovisionResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return DeprovisionResponse{}, fmt.Errorf("HTTP Request failed with status code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return DeprovisionResponse{}, err
	}

	var response DeprovisionResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		fmt.Println(err)
		return DeprovisionResponse{}, err
	}

	return response, nil
}
