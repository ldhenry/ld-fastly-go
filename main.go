package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fastly/compute-sdk-go/fsthttp"

	kvdatastore "github.com/launchdarkly/fastly-go-example/kvdatasore"
	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
	ld "github.com/launchdarkly/go-server-sdk/v7"
	"github.com/launchdarkly/go-server-sdk/v7/ldcomponents"
)

const (
	LD_SDK_KEY        = "sdk-889edf5b-21a6-448e-b84e-0c08c5d54a5b"
	LD_CLIENT_SIDE_ID = "675aea6b1b327709c85da941"
)

func isLocal() bool {
	return os.Getenv("FASTLY_HOSTNAME") == "localhost"
}

func main() {
	// Log service version
	fmt.Println("FASTLY_SERVICE_VERSION:", os.Getenv("FASTLY_SERVICE_VERSION"))

	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
		// Filter requests that have unexpected methods.
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" || r.Method == "DELETE" {
			w.WriteHeader(fsthttp.StatusMethodNotAllowed)
			fmt.Fprintf(w, "This method is not allowed\n")
			return
		}

		clientSideID := "675aea6b1b327709c85da941"
		if isLocal() {
			clientSideID = "local"
		}

		// Initialize the LaunchDarkly SDK
		client, err := ld.MakeCustomClient(LD_SDK_KEY, ld.Config{
			DataSource: ldcomponents.ExternalUpdatesOnly(),
			DataStore: ldcomponents.PersistentDataStore(
				kvdatastore.DataStore().
					ClientSideID(clientSideID).
					KvStoreName("launchdarkly"),
			),
			Events: ldcomponents.NoEvents(),
		}, 5*time.Second)
		if err != nil {
			fmt.Println("Error initializing LaunchDarkly client:", err)
			w.WriteHeader(fsthttp.StatusInternalServerError)
			fmt.Fprintf(w, "Error initializing LaunchDarkly client\n")
			return
		}
		fmt.Println(client.Initialized())

		requestID := os.Getenv("FASTLY_TRACE_ID")
		if requestID == "" {
			requestID = "unknown"
		}

		requestContext := ldcontext.NewBuilder(requestID).
			Kind("fastly-request").
			SetString("fastly_service_version", os.Getenv("FASTLY_SERVICE_VERSION")).
			SetString("fastly_pop", os.Getenv("FASTLY_POP")).
			SetString("fastly_region", os.Getenv("FASTLY_REGION")).
			SetString("fastly_service_id", os.Getenv("FASTLY_SERVICE_ID"))

		ldContext := ldcontext.NewMultiBuilder().Add(ldcontext.New("user-123")).Add(requestContext.Build()).Build()

		animal, reason, err := client.StringVariationDetail("animal", ldContext, "default")
		if err != nil {
			fmt.Println("Error getting animal:", err)
			w.WriteHeader(fsthttp.StatusInternalServerError)
			fmt.Fprintf(w, "Error getting animal: %s", err)
			return
		}

		response := Response{
			Animal:         animal,
			Context:        ldContext,
			Reason:         reason,
			ServiceVersion: os.Getenv("FASTLY_SERVICE_VERSION"),
		}

		w.WriteHeader(fsthttp.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(response); err != nil {
			fmt.Println("Error encoding response:", err)
			w.WriteHeader(fsthttp.StatusInternalServerError)
			fmt.Fprintf(w, "Error encoding response\n")
			return
		}
	})
}
