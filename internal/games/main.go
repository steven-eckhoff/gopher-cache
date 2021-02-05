// The games service.
package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"gopher-cache/internal/common/emulators"
	"gopher-cache/internal/common/logs"
	"gopher-cache/internal/common/server"
	"gopher-cache/internal/games/adapters"
	"gopher-cache/internal/games/app"
	"gopher-cache/internal/games/app/command"
	"gopher-cache/internal/games/app/query"
	"gopher-cache/internal/games/ports"
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
