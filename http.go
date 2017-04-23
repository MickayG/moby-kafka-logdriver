package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/go-plugins-helpers/sdk"
)

type StartLoggingRequest struct {
	File string
	Info logger.Info
}

type StopLoggingRequest struct {
	File string
}

type CapabilitiesResponse struct {
	Err string
	Cap logger.Capability
}

type ReadLogsRequest struct {
	Info   logger.Info
	Config logger.ReadConfig
}


const START_LOGGING_REQUEST string = "/LogDriver.StartLogging"
const STOP_LOGGING_REQUEST string = "/LogDriver.StopLogging"
const LOG_CAPABILITIES_REQUEST string = "/LogDriver.Capabilities"
const READ_LOGS_REQUEST string = "/LogDriver.ReadLogs"

func startLoggingRequest(d LogDriver) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		var req StartLoggingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.Info.ContainerID == "" {
			http.Error(w, "must provide container id in log context", http.StatusBadRequest)
			return
		}

		err := d.StartLogging(req.File, req.Info)
		respond(err, w)
	}
}

func stopLoggingRequest(d LogDriver) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req StopLoggingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := d.StopLogging(req.File)
		respond(err, w)
	}
}

func readLogsRequest(d LogDriver) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ReadLogsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		stream, err := d.ReadLogs(req.Info, req.Config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stream.Close()

		w.Header().Set("Content-Type", "application/x-json-stream")
		wf := ioutils.NewWriteFlusher(w)
		io.Copy(wf, stream)
	}
}


func capabilitiesRequest(d LogDriver) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(&CapabilitiesResponse{
			Cap: logger.Capability{ReadLogs: false},
		})
	}
}

func setupDockerHandlers(h *sdk.Handler, d LogDriver) {
	h.HandleFunc(START_LOGGING_REQUEST, startLoggingRequest(d))
	h.HandleFunc(STOP_LOGGING_REQUEST, stopLoggingRequest(d))
	h.HandleFunc(LOG_CAPABILITIES_REQUEST, capabilitiesRequest(d))
	h.HandleFunc(READ_LOGS_REQUEST, readLogsRequest(d))
}


type response struct {
	Err string
}

func respond(err error, w http.ResponseWriter) {
	var res response
	if err != nil {
		res.Err = err.Error()
	}
	json.NewEncoder(w).Encode(&res)
}
