package marketplace

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func Healthcheck(url string) (int, string, error) {
	client := &http.Client{}

	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return 0, "", err
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0, "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	bodyStr := fmt.Sprintf("%s", body)

	if res.StatusCode != http.StatusOK {
		return res.StatusCode, bodyStr, fmt.Errorf("HTTP Request failed with status code: %d", res.StatusCode)
	} else {
		return res.StatusCode, bodyStr, nil
	}
}
