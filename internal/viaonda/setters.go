package viaonda

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func ClearEntries(ctx context.Context) error {
	form := url.Values{}
	form.Set("acao", "confirma")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ControllerAddress+"/limparRegistros.php", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response")
	}

	return nil
}

func SetGpio(ctx context.Context, port int, state bool) error {
	form := url.Values{}
	form.Set("output", strconv.Itoa(port))

	if state {
		form.Set("status", "on")
	} else {
		form.Set("status", "off")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ControllerAddress+"/gpio.php", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response")
	}

	return nil
}
