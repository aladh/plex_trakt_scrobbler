package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func TestWatchMovieWebhook(t *testing.T) {
	_, enabled := os.LookupEnv("ENABLE_E2E_TESTS")
	if !enabled {
		t.Skip("End to end tests are not enabled")
	}

	// Start server
	go main()

	sendRequest(t)
}

func sendRequest(t *testing.T) {
	payload, err := os.ReadFile("testdata/webhook.json")
	if err != nil {
		t.Fatalf("error opening fixture file: %s", err)
	}

	var formBuffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&formBuffer)
	payloadWriter, err := multipartWriter.CreateFormField("payload")
	if err != nil {
		t.Fatalf("error creating form field: %s", err)
	}

	_, err = payloadWriter.Write(payload)
	if err != nil {
		t.Fatalf("error writing form field: %s", err)
	}

	err = multipartWriter.Close()
	if err != nil {
		t.Fatalf("error closing form writer: %s", err)
	}

	response, err := http.Post("http://localhost:3000", multipartWriter.FormDataContentType(), &formBuffer)
	if err != nil {
		t.Fatalf("error making request to app server: %s", err)
	}

	if response.StatusCode != 200 {
		t.Fatalf("bad response from app server")
	}
}
