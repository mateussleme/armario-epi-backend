package viaonda

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func GetEntries(ctx context.Context) ([]TagEntry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ControllerAddress+"/getTagList.php", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response")
	}

	var data GetEntriesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if data.RecordCount == 0 {
		return nil, nil
	}
	return data.Records, nil
}

func GetGpio(ctx context.Context, port int) (bool, error) {
	form := url.Values{}
	form.Set("input", strconv.Itoa(port))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ControllerAddress+"/gpio.php?"+form.Encode(), nil)
	if err != nil {
		return false, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("bad response")
	}

	var data GpioResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return false, err
	}

	return (strings.ToUpper(data.Status) == "ON"), nil
}
