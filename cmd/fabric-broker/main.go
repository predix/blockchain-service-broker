package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("fabric-sb")

var catalogJson = `
{
  "services": [{
    "name": "hyperledger-fabric",
    "id": "3D7690ED-C611-46F4-9E64-F3D2210EE194",
    "description": "Hyperledger fabric block chain service",
    "tags": ["blockchain"],
    "bindable": true,
    "metadata": {
      "provider": {
        "name": "Hyperledger fabric block chain"
      },
      "listing": {
        "blurb": "Hyperledger fabric",
        "longDescription": "Hyperledger fabric block chain - permissioned block chain"
      },
      "displayName": "Hyperledger service broker"
    },
    "plan_updateable": false,
    "plans": [{
      "name": "basic",
      "id": "15175506-D9F6-4CD8-AA1E-8F0AAFB99C07",
      "description": "Spins up 3 validating nodes in pbft based block chain",
      "metadata": {
        "cost": 99,
        "bullets": [{
          "content": "Dedicated 3 nodes block chain cluster"
        }]
      }
    }
    ]
  }]
}
`

var lastResponse = `
{
	"state": "succeeded",
	"description": "created fabric cluster"
}
`

var errAsyncResponse = `
{
  "error": "AsyncRequired",
  "description": "This service plan requires client support for asynchronous service operations."
}
`

var provisionResponse = `
{
 "dashboardurl": "http://example-dashboard.example.com/9189kdfsk0vfnku",
 "operation": "task10"
}
`

var deprovisionResponse = `
{
 "operation": "task10"
}
`

var (
	scheme       = "https"
	boshUsername = "admin"
	boshPassword = "admin"
	boshAddress  = "192.168.50.4"
	boshPort     = 25555
)

var boshDirectorUrl = fmt.Sprintf("%s://%s:%s@%s:%d", scheme, boshUsername, boshPassword, boshAddress, boshPort)

func main() {
	log.Debug("Starting fabric service broker")
	r := mux.NewRouter()
	r.HandleFunc("/v2/catalog", catalogHandler)
	r.HandleFunc("/v2/service_instances/{instanceId}", provisioningHandler).Methods("PUT")
	r.HandleFunc("/v2/service_instances/{instanceId}", deprovisioningHandler).Methods("DELETE")
	r.HandleFunc("/v2/service_instances/{instanceId}/last_operation", lastOperationHandler)

	http.ListenAndServe(":8999", r)
}

func catalogHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Serving /v2/catalog")
	w.Write([]byte(catalogJson))
}

func provisioningHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Handling PUT /v2/service_instances")
	// vars := mux.Vars(r)
	// instanceId := vars["instanceId"]

	query := r.URL.Query()
	async := query["accepts_incomplete"]
	if len(async) < 1 || async[0] != "true" {
		w.WriteHeader(422)
		w.Write([]byte(errAsyncResponse))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(provisionResponse))
}

func deprovisioningHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Handling DELETE /v2/service_instances")
	// vars := mux.Vars(r)
	// instanceId := vars["instanceId"]

	query := r.URL.Query()
	async := query["accepts_incomplete"]
	if len(async) < 1 || async[0] != "true" {
		w.WriteHeader(422)
		w.Write([]byte(errAsyncResponse))
		return
	}

	// TODO: get service_id and plan_id
	// 410 if it doesn't exist

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(deprovisionResponse))

}

func lastOperationHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Handling GET /v2/service_instances/:instanceId/last_operation")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(lastResponse))
}
