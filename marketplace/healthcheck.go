package marketplace

import (
	"fmt"
	"net/http"
)

func Healthcheck(url string) (int, error) {
	client := &http.Client{}

	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return res.StatusCode, fmt.Errorf("HTTP Request failed with status code: %d", res.StatusCode)
	} else {
		return res.StatusCode, nil
	}
}
