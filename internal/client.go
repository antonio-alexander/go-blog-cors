package internal

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

func simulateCors(method, uri, origin string, headerKeys ...string) error {
	// Simulate preflight OPTIONS request
	request, err := http.NewRequest(http.MethodOptions, uri, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Origin", origin)
	request.Header.Set("Access-Control-Request-Method", method)
	request.Header.Set("Access-Control-Request-Headers", strings.Join(headerKeys, ","))

	//execute request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	switch response.StatusCode {
	default:
		return errors.Errorf("unexpected status code: %d", response.StatusCode)
	case http.StatusOK, http.StatusNoContent:
		allowedOrigins := response.Header.Get("Access-Control-Allow-Origin")
		allowedHeaders := response.Header.Get("Access-Control-Allow-Headers")
		allowedMethods := response.Header.Get("Access-Control-Allow-Methods")
		allowCredentials := response.Header.Get("Access-Control-Allow-Credentials")
		if allowedOrigins == "" || allowedHeaders == "" ||
			allowedMethods == "" {
			return errors.New("CORS error encountered")
		}
		fmt.Printf("origins allowed: %s\n", allowedOrigins)
		fmt.Printf("headers allowed: %s\n", allowedHeaders)
		fmt.Printf("methods allowed: %s\n", allowedMethods)
		fmt.Printf("credentials allowed: %s\n", allowCredentials)
	}
	return nil
}

func ClientHelloWorld() (string, error) {
	address := "http://" + configAddress
	if configPort != "" {
		address += ":" + configPort
	}
	uri, method := address+"/", http.MethodGet
	if err := simulateCors(method, uri, configCorsOrigin,
		"Origin"); err != nil {
		return "", err
	}
	request, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return "", err
	}
	request.Header.Add("Origin", configCorsOrigin)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ClientAuthorize(username, password string) (string, error) {
	address := "http://" + configAddress
	if configPort != "" {
		address += ":" + configPort
	}
	uri, method := address+"/authorize", http.MethodPost
	if err := simulateCors(method, uri, configCorsOrigin,
		"Origin", "Authorization"); err != nil {
		return "", err
	}
	request, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return "", err
	}
	request.Header.Add("Origin", configCorsOrigin)
	request.SetBasicAuth(username, password)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
