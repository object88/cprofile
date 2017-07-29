package bacon

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const baseURL = "http://baconipsum.com/api/?"

// GenerateBaconIpsum uses https://baconipsum.com/json-api/ to
// retreve random text content
func GenerateBaconIpsum(optionFns ...*OptionFn) string {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	o := defaultOptions()

	url := baseURL + makeUrlOptions(o)
	resp, err := client.Get(url)
	if err != nil {
		return "Error making request"
	}
	defer resp.Body.Close()

	fmt.Printf("Status code: %d\n", resp.StatusCode)

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return "Failed to read response"
	}

	return string(body)
}

func makeUrlOptions(o *Options) string {
	return "?type=meat-and-filler"
}
