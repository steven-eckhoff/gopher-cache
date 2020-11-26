// The games service.
package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/steven-eckhoff/gopher-cache-open/internal/common/emulators"
	"github.com/steven-eckhoff/gopher-cache-open/internal/common/logs"
	"github.com/steven-eckhoff/gopher-cache-open/internal/common/server"
	"github.com/steven-eckhoff/gopher-cache-open/internal/games/adapters"
	"github.com/steven-eckhoff/gopher-cache-open/internal/games/app"
	"github.com/steven-eckhoff/gopher-cache-open/internal/games/app/command"
	"github.com/steven-eckhoff/gopher-cache-open/internal/games/app/query"
	"github.com/steven-eckhoff/gopher-cache-open/internal/games/ports"
	"net/http"
)

func init() {
	logs.Init()
}

func main() {
	ctx := context.Background()

	application, cleanup := newLocalApplication(ctx)
	defer cleanup()

	logrus.Info("Starting HTTP server")

	server.RunHTTPServer(ctx, func(router chi.Router) http.Handler {
		return ports.APIHandler(ports.NewHTTPServer(application), router)
	})
}

func newLocalApplication(ctx context.Context) (app.Application, func()) {
	client, cleanup := emulators.NewFirestoreClient(ctx)

	gamesRepository, err := adapters.NewFirestoreGameRepository(client)
	if err != nil {
		panic(err)
	}

	return app.Application{
			Commands: app.Commands{
				CreateGame:      command.NewCreateGameHandler(gamesRepository),
				CreateGameState: command.NewCreateGameStateHandler(gamesRepository),
				UpdateGameState: command.NewUpdateGameStateHandler(gamesRepository),
			},
			Queries: app.Queries{
				GetGames:  query.NewReadGamesHandler(gamesRepository),
				GetPlayer: query.NewReadPlayerHandler(gamesRepository),
				GetState:  query.NewReadStateHandler(gamesRepository),
			},
		}, func() {
			_ = client.Close()
			cleanup()
		}
}
