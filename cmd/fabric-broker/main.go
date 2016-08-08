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
    "max_db_per_node": 5,
    "bindable": true,
    "metadata": {
      "provider": {
        "name": "The name"
      },
      "listing": {
        "imageUrl": "http://example.com/cat.gif",
        "blurb": "Add a blurb here",
        "longDescription": "A long time ago, in a galaxy far far away..."
      },
      "displayName": "The Fake Broker"
    },
    "dashboard_client": {
      "id": "398e2f8e-XXXX-XXXX-XXXX-19a71ecbcf64",
      "secret": "277cabb0-XXXX-XXXX-XXXX-7822c0a90e5d",
      "redirect_uri": "http://localhost:1234"
    },
    "plan_updateable": true,
    "plans": [{
      "name": "fake-plan",
      "id": "d3031751-XXXX-XXXX-XXXX-a42377d3320e",
      "description": "Shared fake Server, 5tb persistent disk, 40 max concurrent connections",
      "max_storage_tb": 5,
      "metadata": {
        "cost": 0,
        "bullets": [{
          "content": "Shared fake server"
        }, {
          "content": "5 TB storage"
        }, {
          "content": "40 concurrent connections"
        }]
      }
    }, {
      "name": "fake-async-plan",
      "id": "0f4008b5-XXXX-XXXX-XXXX-dace631cd648",
      "description": "Shared fake Server, 5tb persistent disk, 40 max concurrent connections. 100 async",
      "max_storage_tb": 5,
      "metadata": {
        "cost": 0,
        "bullets": [{
          "content": "40 concurrent connections"
        }]
      }
    }, {
      "name": "fake-async-only-plan",
      "id": "8d415f6a-XXXX-XXXX-XXXX-e61f3baa1c77",
      "description": "Shared fake Server, 5tb persistent disk, 40 max concurrent connections. 100 async",
      "max_storage_tb": 5,
      "metadata": {
        "cost": 0,
        "bullets": [{
          "content": "40 concurrent connections"
        }]
      }
    }]
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
