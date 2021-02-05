// Package ports contains ports into the games application.
package ports

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gopher-cache/internal/common/auth"
	"gopher-cache/internal/common/server/httperr"
	"gopher-cache/internal/games/app"
	"gopher-cache/internal/games/app/command"
	"gopher-cache/internal/games/app/query"
	"gopher-cache/internal/games/domain/game"
	"net/http"
	"strconv"
)

// HTTPServer maps HTTP request to application commands and queries.
type HTTPServer struct {
	app app.Application
}

// Creates a new HTTP server.
func NewHTTPServer(app app.Application) HTTPServer {
	return HTTPServer{app: app}
}

// CreateGame expects the body of the request to have JSON in the form of
// command.CreateGame.
func (h HTTPServer) CreateGame(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromContext(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	gameUser, err := game.NewUser(user.UUID, user.Number)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	cmd := new(command.CreateGame)

	err = render.Decode(r, cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	cmd.Creator = gameUser

	err = h.app.Commands.CreateGame.Handle(r.Context(), *cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

// CreateGameState expects the body of the request to have JSON in the form of
// command.CreateGameState.
func (h HTTPServer) CreateGameState(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromContext(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	gameUser, err := game.NewUser(user.UUID, user.Number)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	cmd := new(command.CreateGameState)

	err = render.Decode(r, cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	cmd.User = gameUser

	resp, err := h.app.Commands.CreateGameState.Handle(r.Context(), *cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Respond(w, r, resp)
}

// UpdateGameState expects the body of the request to have JSON in the form of
// command.UpdateGameState. A URL param player-number must also be present.
func (h HTTPServer) UpdateGameState(w http.ResponseWriter, r *http.Request) {
	// We'll use the user in the context to authenticate the request.
	_, err := auth.UserFromContext(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	cmd := new(command.UpdateGameState)

	err = render.Decode(r, cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	cmd.PlayerNumber = chi.URLParam(r, "player-number")

	resp, err := h.app.Commands.UpdateGameState.Handle(r.Context(), *cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Respond(w, r, resp)
}

func gameQueryParamsFromRequest(r *http.Request) (limit, offset int, options []query.GameOption, err error) {
	values := r.URL.Query()

	limitStr := values.Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return
		}
		values.Del("limit")
	} else {
		limit = 10
	}

	offsetStr := values.Get("offset")
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return
		}
		values.Del("offset")
	} else {
		offset = 0
	}

	for k := range values {
		options = append(options, query.GameOption{
			Key:   k,
			Op:    "==",
			Value: values.Get(k),
		})
	}

	return
}

// GetGames queries for games. Supported queries are camel cased member names of
// the game.Game type. limit and offset are also available to support pagination.
// If no limit is given then it defaults to 10.
func (h HTTPServer) GetGames(w http.ResponseWriter, r *http.Request) {
	// We'll use the user in the context to authenticate the request.
	_, err := auth.UserFromContext(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	limit, offset, options, err := gameQueryParamsFromRequest(r)
	if err != nil {
		httperr.BadRequest("query-params", err, w, r)
		return
	}

	games, err := h.app.Queries.GetGames.Handle(r.Context(), limit, offset, options...)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Respond(w, r, games)
}

// GetPlayer queries for a players UUID. The UUID is expressed in a URL param uuid.
func (h HTTPServer) GetPlayer(w http.ResponseWriter, r *http.Request) {
	// We'll use the user in the context to authenticate the request.
	_, err := auth.UserFromContext(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	player, err := h.app.Queries.GetPlayer.Handle(r.Context(), chi.URLParam(r, "uuid"))
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Respond(w, r, player)
}

// GetState queries for a game state by UUID. The UUID is expressed in a URL param uuid.
func (h HTTPServer) GetState(w http.ResponseWriter, r *http.Request) {
	// We'll use the user in the context to authenticate the request.
	_, err := auth.UserFromContext(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	state, err := h.app.Queries.GetState.Handle(r.Context(), chi.URLParam(r, "uuid"))
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Respond(w, r, state)
}
