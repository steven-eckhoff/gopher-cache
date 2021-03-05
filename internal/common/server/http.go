// Package server implements a library helping with basic HTTP server operations.
package server

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"gopher-cache/internal/common/auth"
	"gopher-cache/internal/common/logs"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

// RunHTTPServer runs an HTTP server listening on the port specified by PORT in the environment.
// This function will block until the server is running.
// On SIGINT or SIGTERM the server will be shutdown cleanly.
func RunHTTPServer(ctx context.Context, createHandler func(router chi.Router) http.Handler) {
	apiRouter := chi.NewRouter()
	setMiddlewares(apiRouter)

	rootRouter := chi.NewRouter()
	rootRouter.Mount("/api", createHandler(apiRouter))

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: rootRouter,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		_ = srv.ListenAndServe()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		_ = <-sigs
		_ = srv.Shutdown(ctx)
	}()

	wg.Wait()
}

func setMiddlewares(router *chi.Mux) {
	router.Use(logs.NewStructuredLogger(logrus.StandardLogger()))
	addAuthMiddleware(router)
}

func addAuthMiddleware(router *chi.Mux) {
	if mockAuth, _ := strconv.ParseBool(os.Getenv("MOCK_AUTH")); mockAuth {
		logrus.Info("Using JWT mock auth")
		router.Use(auth.HttpMockMiddleware)
		return
	}

	var opts []option.ClientOption
	if file := os.Getenv("SERVICE_ACCOUNT_FILE"); file != "" {
		opts = append(opts, option.WithCredentialsFile(file))
	}

	config := &firebase.Config{ProjectID: os.Getenv("GCP_PROJECT")}
	firebaseApp, err := firebase.NewApp(context.Background(), config, opts...)
	if err != nil {
		logrus.Fatalf("error initializing app: %v\n", err)
	}

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		logrus.WithError(err).Fatal("Unable to create firebase Auth client")
	}

	router.Use(auth.FirebaseHttpMiddleware{authClient}.Middleware)
}
