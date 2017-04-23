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
	"io/ioutil"
)

func TestDriverCalledOnStartLoggingCall(t *testing.T) {
	var driver testDriver
	startLoggingRequestContent := StartLoggingRequest{File:"logs", Info:logger.Info{ContainerID:"abdud"}}

	outContent, err := json.Marshal(startLoggingRequestContent)
	if err != nil {	t.Fatal(err) }

	ht := executeHttpRequest(t, driver, START_LOGGING_REQUEST, startLoggingRequest(&driver), bytes.NewReader(outContent))


	assertCodeReceived(t, 200, ht)
	assert.Equal(t, true, driver.startLoggingCalled)
	assert.Equal(t, false, driver.stopLoggingCalled)
}

func TestDriverCalledOnStopLoggingCall(t *testing.T) {
	var driver testDriver
	stopLoggingRequestContent := StopLoggingRequest{File: "abcd"}

	outContent, err := json.Marshal(stopLoggingRequestContent)
	if err != nil {	t.Fatal(err) }

	ht := executeHttpRequest(t, driver, STOP_LOGGING_REQUEST, stopLoggingRequest(&driver), bytes.NewReader(outContent))

	assertCodeReceived(t, 200, ht)
	assert.Equal(t, false, driver.startLoggingCalled)
	assert.Equal(t, true, driver.stopLoggingCalled)
}

func TestDriverQueriedForCapabilities(t *testing.T) {
	var driver testDriver

	ht := executeHttpRequest(t, driver, LOG_CAPABILITIES_REQUEST, capabilitiesRequest(&driver), nil)

	assertCodeReceived(t, 200, ht)
	// Check that our test driver was indeed called
	assert.Equal(t, true, driver.capabilitiesCalled)

	// Test that the data was serialized correctly
	respContent, err := ioutil.ReadAll(ht.Body)
	if err != nil {	t.Fatal(err) }

	var capability CapabilitiesResponse
	err = json.Unmarshal(respContent, &capability)
	if err != nil {	t.Fatal(err) }
	assert.Equal(t, true, capability.Cap.ReadLogs)
}

func executeHttpRequest(t *testing.T, driver testDriver, uri string, fn func(w http.ResponseWriter, r *http.Request), body io.Reader) (*httptest.ResponseRecorder) {
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		t.Fatal(err)
	}
	ht := httptest.NewRecorder()
	handler := http.HandlerFunc(fn)
	handler.ServeHTTP(ht, req)
	return ht
}


func assertCodeReceived(t *testing.T, code int, recorder *httptest.ResponseRecorder) {
	if status := recorder.Code; status != code {
		t.Errorf("Wrong error code. Got %v expected %v", status, code)
		t.Fail()
	}
}


type testDriver struct {
	startLoggingCalled bool
	stopLoggingCalled bool
	capabilitiesCalled bool
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

func (d *testDriver) GetCapability() logger.Capability {
	d.capabilitiesCalled = true
	return logger.Capability{ReadLogs: true}
}

