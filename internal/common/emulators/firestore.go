// Package emulators implements a library that assist in working with various emulators.
// The current emulators supported are the following:
// - Firestore
package emulators

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

const (
	firestoreEmulatorHost = "FIRESTORE_EMULATOR_HOST"
	defaultHost           = "localhost:9093"
)

var (
	lock sync.Mutex
)

// NewFirestoreClient returns the client connected to the Firestore emulator.
// If another goroutine is using the emulator this function will block until it is free.
// A clean up function will be returned that must be called after you are done to free
// the emulator for other goroutines. FIRESTORE_EMULATOR_HOST must be set in the
// environment or the default localhost:9093 will be used.
func NewFirestoreClient(ctx context.Context) (*firestore.Client, func()) {
	lock.Lock()

	if os.Getenv(firestoreEmulatorHost) == "" {
		err := os.Setenv(firestoreEmulatorHost, defaultHost)
		if err != nil {
			panic(err)
		}
	}

	client, err := firestore.NewClient(ctx, "local")
	if err != nil {
		log.Fatalf("firebase.NewClient err: %v", err)
	}

	return client, func() {
		clearData()
		lock.Unlock()
	}
}

func clearData() {
	url := fmt.Sprintf("http://%s/emulator/v1/projects/local/databases/(default)/documents", os.Getenv("FIRESTORE_EMULATOR_HOST"))
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic(errors.New("could not clear emulator data"))
	}
}
