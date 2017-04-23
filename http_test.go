package main

import (
	"testing"
	"github.com/docker/docker/daemon/logger"
	"io"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/http"
	"encoding/json"
	"bytes"
)

func TestDriverCalledOnStartLoggingCall(t *testing.T) {
	var driver testDriver

	startLoggingRequestContent := StartLoggingRequest{File:"logs", Info:logger.Info{ContainerID:"abdud"}}
	outContent, err := json.Marshal(startLoggingRequestContent)
	if err != nil {	t.Fatal(err) }

	req, err := http.NewRequest("POST", "/StartLogging", bytes.NewReader(outContent))
	if err != nil {	t.Fatal(err) }

	ht := httptest.NewRecorder()
	handler := http.HandlerFunc(startLoggingRequest(&driver))
	handler.ServeHTTP(ht, req)

	assertCodeReceived(t, 200, ht)
	assert.Equal(t, true, driver.startLoggingCalled)
	assert.Equal(t, false, driver.stopLoggingCalled)
}

func TestDriverCalledOnStopLoggingCall(t *testing.T) {
	var driver testDriver

	stopLoggingRequestContent := StopLoggingRequest{File: "abcd"}

	outContent, err := json.Marshal(stopLoggingRequestContent)
	if err != nil {	t.Fatal(err) }

	req, err := http.NewRequest("POST", "/StopLogging", bytes.NewReader(outContent))
	if err != nil {	t.Fatal(err) }

	ht := httptest.NewRecorder()
	handler := http.HandlerFunc(stopLoggingRequest(&driver))
	handler.ServeHTTP(ht, req)

	assertCodeReceived(t, 200, ht)
	assert.Equal(t, false, driver.startLoggingCalled)
	assert.Equal(t, true, driver.stopLoggingCalled)
}

type testDriver struct {
	startLoggingCalled bool
	stopLoggingCalled bool
}


func (d *testDriver) StartLogging(file string, logCtx logger.Info) error {
	d.startLoggingCalled = true
	return nil
}

func (d *testDriver) StopLogging(file string) error {
	d.stopLoggingCalled = true
	return nil
}


func (d *testDriver) ReadLogs(info logger.Info, config logger.ReadConfig) (io.ReadCloser, error) {
	//TODO
	return nil, nil
}

func assertCodeReceived(t *testing.T, code int, recorder *httptest.ResponseRecorder) {
	if status := recorder.Code; status != code {
		t.Errorf("Wrong error code. Got %v expected %v", status, code)
		t.Fail()
	}
}
