package subscription

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"syscall"
	"time"
)

type SubscriptionRequest struct {
	Interval     int `json:"interval"`
	SamplePeriod int `json:"samplePeriod"`
}

type SubscriptionConfigration interface {
	UpdateConfig(cfg SubscriptionRequest) error
	DeleteConfig() error
}

type SubscriptionListener struct {
	Address             string
	ConfigDir           string
	listener            net.Listener
	statsConfigHandlers map[string]SubscriptionConfigration
	stateConfigHandlers map[string]SubscriptionConfigration
}

// Standard HTTP response definition
type HttpResponse struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Resource     string `json:"resource"`
}

// Method to send standard HTTP response
// Parameters:
//  * w http.ResponseWriter: HTTP output stream
//  * rc int: HTTP return code
//  * msg string: Error/success message with additional details
//  * resource string: request URL
// Returns:
//  * None
func response(w http.ResponseWriter, rc int, msg string, resource string) {
	w.WriteHeader(rc)
	w.Header().Set("Content-Type", "application/json")
	resp := make([]HttpResponse, 1)
	resp[0] = HttpResponse{ErrorCode: rc, ErrorMessage: msg, Resource: resource}
	jpl, err := json.Marshal(resp)
	if err != nil {
		log.Printf("E! Failed to encode JSON %+v\n", resp)
		return
	}
	w.Write(jpl)
}

func (h *SubscriptionListener) subscribeStats(metric string, w http.ResponseWriter, r *http.Request) {
	handler, ok := h.statsConfigHandlers[metric]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var pl SubscriptionRequest
	err := json.NewDecoder(r.Body).Decode(&pl)
	if err != nil {
		log.Printf("E! Failed to parse payload: %+v\n", err)
		response(w, http.StatusBadRequest, err.Error(), r.URL.Path)
		return
	}
	err = handler.UpdateConfig(pl)
	if err == nil {
		// telegraf re-read config when receives SIGHUP
		syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
		response(w, http.StatusOK, "Success", r.URL.Path)
	} else {
		response(w, http.StatusBadRequest, err.Error(), r.URL.Path)
	}
}

func (h *SubscriptionListener) stopStats(metric string, w http.ResponseWriter, r *http.Request) {
	handler, ok := h.statsConfigHandlers[metric]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err := handler.DeleteConfig()
	if err == nil {
		// telegraf re-read config when receives SIGHUP
		syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
		response(w, http.StatusOK, "Success", r.URL.Path)
	} else {
		response(w, http.StatusBadRequest, err.Error(), r.URL.Path)
	}
}

func (h *SubscriptionListener) subscribeState(metric string, w http.ResponseWriter, r *http.Request) {
	handler, ok := h.stateConfigHandlers[metric]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var pl SubscriptionRequest
	err := json.NewDecoder(r.Body).Decode(&pl)
	if err != nil {
		log.Printf("E! Failed to parse payload: %+v\n", err)
		response(w, http.StatusBadRequest, err.Error(), r.URL.Path)
		return
	}
	err = handler.UpdateConfig(pl)
	if err == nil {
		// telegraf re-read config when receives SIGHUP
		syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
		response(w, http.StatusOK, "Success", r.URL.Path)
	} else {
		response(w, http.StatusBadRequest, err.Error(), r.URL.Path)
	}
}

func (h *SubscriptionListener) stopState(metric string, w http.ResponseWriter, r *http.Request) {
	handler, ok := h.stateConfigHandlers[metric]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err := handler.DeleteConfig()
	if err == nil {
		// telegraf re-read config when receives SIGHUP
		syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
		response(w, http.StatusOK, "Success", r.URL.Path)
	} else {
		response(w, http.StatusBadRequest, err.Error(), r.URL.Path)
	}
}

func (h *SubscriptionListener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	if len(path) == 5 && path[1] == "v1" && path[2] == "subscribe" {
		if path[3] == "stats" {
			switch r.Method {
			case "POST":
				h.subscribeStats(path[4], w, r)
				return
			case "DELETE":
				h.stopStats(path[4], w, r)
				return
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
		}
		if path[3] == "state" {
			switch r.Method {
			case "POST":
				h.subscribeState(path[4], w, r)
				return
			case "DELETE":
				h.stopState(path[4], w, r)
				return
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h *SubscriptionListener) Start() error {
	h.statsConfigHandlers = map[string]SubscriptionConfigration{}
	h.statsConfigHandlers["cpu"] = &SubscriptionConfigrationStatsCpu{ConfigDir: h.ConfigDir}

	var listener net.Listener
	listener, err := net.Listen("tcp", h.Address)
	if err != nil {
		return err
	}
	h.listener = listener
	log.Printf("I! Starting REST API subscription server at: %s", h.Address)

	server := &http.Server{
		Addr:         h.Address,
		Handler:      h,
		ReadTimeout:  time.Duration(30 * time.Second),
		WriteTimeout: time.Duration(30 * time.Second),
	}

	if err := server.Serve(h.listener); err != nil {
		if !errors.Is(err, net.ErrClosed) {
			fmt.Printf("Serve failed: %v", err)
		}
	}
	return nil
}

func (h *SubscriptionListener) Stop() {
	if h.listener != nil {
		// Ignore the returned error as we cannot do anything about it anyway
		//nolint:errcheck,revive
		h.listener.Close()
	}
}
