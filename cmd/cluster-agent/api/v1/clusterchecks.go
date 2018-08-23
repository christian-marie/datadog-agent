// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

// +build clusterchecks

package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	apitypes "github.com/DataDog/datadog-agent/cmd/cluster-agent/api/types"
	cctypes "github.com/DataDog/datadog-agent/pkg/clusteragent/clusterchecks/types"
)

// Install registers v1 API endpoints
func installClusterCheckEndpoints(r *mux.Router, sc apitypes.ServerContext) {
	r.HandleFunc("/clusterchecks/status/{nodeName}", postCheckStatus(sc)).Methods("POST")
	r.HandleFunc("/clusterchecks/configs/{nodeName}", getCheckConfigs(sc)).Methods("GET")
	r.HandleFunc("/clusterchecks", getAllCheckConfigs(sc)).Methods("GET")
}

// postCheckStatus is used by the node-agent's config provider
func postCheckStatus(sc types.ServerContext) func(w http.ResponseWriter, r *http.Request) {
	if sc.ClusterCheckHandler == nil {
		return func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nodeName := vars["nodeName"]

		decoder := json.NewDecoder(r.Body)
		var status cctypes.NodeStatus
		err := decoder.Decode(&status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := sc.ClusterCheckHandler.PostStatus(nodeName, status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		slcB, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(slcB) != 0 {
			w.WriteHeader(http.StatusOK)
			w.Write(slcB)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}
}

// getCheckConfigs is used by the node-agent's config provider
func getCheckConfigs(sc types.ServerContext) func(w http.ResponseWriter, r *http.Request) {
	if sc.ClusterCheckHandler == nil {
		return func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nodeName := vars["nodeName"]
		response, err := sc.ClusterCheckHandler.GetConfigs(nodeName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		slcB, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(slcB) != 0 {
			w.WriteHeader(http.StatusOK)
			w.Write(slcB)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}
}

// getAllCheckConfigs is used by the clustercheck config
func getAllCheckConfigs(sc types.ServerContext) func(w http.ResponseWriter, r *http.Request) {
	if sc.ClusterCheckHandler == nil {
		return func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		response, err := sc.ClusterCheckHandler.GetAllConfigs()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		slcB, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(slcB) != 0 {
			w.WriteHeader(http.StatusOK)
			w.Write(slcB)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}
}
