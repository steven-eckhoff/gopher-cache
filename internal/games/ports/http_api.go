package ports

import (
	"github.com/go-chi/chi"
	"net/http"
)

// ServerInterface represents the interface the server must have to satisfy the API.
type ServerInterface interface {
	// /games POST
	CreateGame(w http.ResponseWriter, r *http.Request)
	// /game-states POST
	CreateGameState(w http.ResponseWriter, r *http.Request)
	// /game-states/{player-number} PUT
	UpdateGameState(w http.ResponseWriter, r *http.Request)
	// /games GET
	GetGames(w http.ResponseWriter, r *http.Request)
	// /players/uuid GET
	GetPlayer(w http.ResponseWriter, r *http.Request)
	// /game-states/uuid GET
	GetState(w http.ResponseWriter, r *http.Request)
}

// APIHandler binds a server implementing the ServerInterface to the games API using the given router.
func APIHandler(si ServerInterface, r chi.Router) http.Handler {
	r.Post("/games", si.CreateGame)
	r.Post("/game-states", si.CreateGameState)
	r.Put("/game-states/{player-number}", si.UpdateGameState)
	r.Get("/games", si.GetGames)
	r.Get("/players/{uuid}", si.GetPlayer)
	r.Get("/game-states/{uuid}", si.GetState)

	return r
}
