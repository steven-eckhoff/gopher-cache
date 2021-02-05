// Package adapters provides adapters for games.
package adapters

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"gopher-cache/internal/games/app/query"
	"gopher-cache/internal/games/domain/game"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type firestoreGameModel struct {
	UUID        string                `firestore:"uuid"`
	CreatorUUID string                `firestore:"creatorUUID"`
	Title       string                `firestore:"title"`
	Description string                `firestore:"description"`
	Levels      []firestoreLevelModel `firestore:"levels"`
	Ending      string                `firestore:"ending"`
	Kind        string                `firestore:"kind"`
	City        string                `firestore:"city"`
	State       string                `firestore:"state"`
	Country     string                `firestore:"country"`
	Value       int                   `firestore:"value"`
}

type firestoreLevelModel struct {
	Title       string   `firestore:"title"`
	Description string   `firestore:"description"`
	Clues       []string `firestore:"clues"`
	Answers     []string `firestore:"answers"`
}

type firestorePlayerModel struct {
	UUID                 string `firestore:"uuid"`
	Number               string `firestore:"number"`
	GamesStarted         int    `firestore:"gamesStarted"`
	GamesFinished        int    `firestore:"gamesFinished"`
	TotalPoints          int    `firestore:"totalPoints"`
	CurrentGameStateUUID string `firestore:"currentGameStateUUID"`
}

type firestoreStateModel struct {
	UUID            string        `firestore:"uuid"`
	PlayerUUID      string        `firestore:"playerUUID"`
	GameUUID        string        `firestore:"gameUUID"`
	GameLevels      int           `firestore:"gameLevels"`
	Level           int           `firestore:"level"`
	Clue            int           `firestore:"clue"`
	Completed       bool          `firestore:"completed"`
	CurrentResponse game.Response `firestore:"currentResponse"`
}

var _ game.Repository = FirestoreGameRepository{}

// FirestoreGameRepository implements the Firestore game repository.
type FirestoreGameRepository struct {
	client *firestore.Client
}

// NewFirestoreGameRepository creates a new game repository using Firestore.
func NewFirestoreGameRepository(client *firestore.Client) (FirestoreGameRepository, error) {
	if client == nil {
		return FirestoreGameRepository{}, errors.New("nil firestore client")
	}

	return FirestoreGameRepository{client: client}, nil
}

func (r FirestoreGameRepository) AddGame(ctx context.Context, game *game.Game) error {
	model := firestoreGameModel{
		UUID:        game.UUID(),
		CreatorUUID: game.CreatorUUID(),
		Title:       game.Title(),
		Description: game.Description(),
		Ending:      game.Ending(),
		Kind:        game.Kind(),
		City:        game.City(),
		State:       game.State(),
		Country:     game.Country(),
		Value:       game.Value(),
	}

	for _, level := range game.Levels() {
		model.Levels = append(model.Levels, firestoreLevelModel{
			Title:       level.Title(),
			Description: level.Description(),
			Clues:       level.Clues(),
			Answers:     level.Answers(),
		})
	}

	doc := r.client.Doc("games/" + game.UUID())
	_, err := doc.Create(ctx, model)
	if err != nil {
		return err
	}

	return nil
}

func (r FirestoreGameRepository) GetGame(ctx context.Context, gameUUID string) (*game.Game, error) {
	docsnap, err := r.client.Doc("games/" + gameUUID).Get(ctx)
	if err != nil {
		return nil, err
	}

	model := new(firestoreGameModel)

	err = docsnap.DataTo(model)
	if err != nil {
		return nil, err
	}

	var levels []*game.Level
	for _, level := range model.Levels {
		levels = append(levels, game.UnmarshalLevelFromDatabase(
			level.Title,
			level.Description,
			level.Clues,
			level.Answers))
	}
	g, err := game.UnmarshalFromDataBase(
		model.UUID,
		model.CreatorUUID,
		model.Title,
		model.Description,
		levels,
		model.Ending,
		model.Kind,
		model.City,
		model.State,
		model.Country,
		model.Value)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (r FirestoreGameRepository) AddPlayer(ctx context.Context, player *game.Player) error {
	model := firestorePlayerModel{
		UUID:                 player.UUID(),
		Number:               player.Number(),
		GamesStarted:         player.GamesStarted(),
		GamesFinished:        player.GamesFinished(),
		TotalPoints:          player.TotalPoints(),
		CurrentGameStateUUID: player.CurrentGameStateUUID(),
	}

	doc := r.client.Doc("players/" + player.UUID())
	_, err := doc.Create(ctx, model)
	if err != nil {
		return err
	}

	return nil
}

func (r FirestoreGameRepository) GetPlayer(ctx context.Context, playerUUID string) (*game.Player, error) {
	docsnap, err := r.client.Doc("players/" + playerUUID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, game.ErrorPlayerNotFound
		} else {
			return nil, err
		}
	}

	model := new(firestorePlayerModel)

	err = docsnap.DataTo(model)
	if err != nil {
		return nil, err
	}

	return game.UnmarshalPlayerFromDatabase(
		model.UUID,
		model.Number,
		model.GamesStarted,
		model.GamesFinished,
		model.TotalPoints,
		model.CurrentGameStateUUID), nil
}

func (r FirestoreGameRepository) GetPlayerByNumber(ctx context.Context, playerNumber string) (*game.Player, error) {
	players := r.client.Collection("players")

	q := players.Where("number", "==", playerNumber).Limit(1)

	iter := q.Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err != nil {
		if err == iterator.Done {
			return nil, game.ErrorPlayerNotFound
		} else {
			return nil, err
		}
	}

	model := new(firestorePlayerModel)

	err = doc.DataTo(model)
	if err != nil {
		return nil, err
	}

	return game.UnmarshalPlayerFromDatabase(
		model.UUID,
		model.Number,
		model.GamesStarted,
		model.GamesFinished,
		model.TotalPoints,
		model.CurrentGameStateUUID), nil
}

func (r FirestoreGameRepository) AddState(ctx context.Context, state *game.State) error {
	model := firestoreStateModel{
		UUID:            state.UUID(),
		PlayerUUID:      state.PlayerUUID(),
		GameUUID:        state.GameUUID(),
		GameLevels:      state.GameLevels(),
		Level:           state.Level(),
		Clue:            state.Clue(),
		Completed:       state.Completed(),
		CurrentResponse: state.CurrentResponse(),
	}

	_, err := r.client.Doc("game-states/"+state.UUID()).Create(ctx, model)
	if err != nil {
		return err
	}

	return nil
}

func (r FirestoreGameRepository) GetState(ctx context.Context, uuid string) (*game.State, error) {
	docsnap, err := r.client.Doc("game-states/" + uuid).Get(ctx)
	if err != nil {
		return nil, err
	}

	model := new(firestoreStateModel)

	err = docsnap.DataTo(model)
	if err != nil {
		return nil, err
	}

	return game.UnmarshalGameStateFromDatabase(
		model.UUID,
		model.PlayerUUID,
		model.GameUUID,
		model.GameLevels,
		model.Level,
		model.Clue,
		model.Completed,
		model.CurrentResponse), nil
}

func (r FirestoreGameRepository) UpdateState(ctx context.Context, state *game.State) error {
	model := firestoreStateModel{
		UUID:            state.UUID(),
		PlayerUUID:      state.PlayerUUID(),
		GameUUID:        state.GameUUID(),
		GameLevels:      state.GameLevels(),
		Level:           state.Level(),
		Clue:            state.Clue(),
		Completed:       state.Completed(),
		CurrentResponse: state.CurrentResponse(),
	}

	_, err := r.client.Doc("game-states/"+state.UUID()).Set(ctx, model)
	if err != nil {
		return err
	}

	return nil
}

func (r FirestoreGameRepository) AddStateAndUpdatePlayer(ctx context.Context, state *game.State, player *game.Player) error {
	stateModel := firestoreStateModel{
		UUID:            state.UUID(),
		PlayerUUID:      state.PlayerUUID(),
		GameUUID:        state.GameUUID(),
		GameLevels:      state.GameLevels(),
		Level:           state.Level(),
		Clue:            state.Clue(),
		Completed:       state.Completed(),
		CurrentResponse: state.CurrentResponse(),
	}

	playerModel := firestorePlayerModel{
		UUID:                 player.UUID(),
		Number:               player.Number(),
		GamesStarted:         player.GamesStarted(),
		GamesFinished:        player.GamesFinished(),
		TotalPoints:          player.TotalPoints(),
		CurrentGameStateUUID: player.CurrentGameStateUUID(),
	}

	s := r.client.Doc("game-states/" + state.UUID())
	p := r.client.Doc("players/" + player.UUID())

	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		err := tx.Create(s, stateModel)
		if err != nil {
			return err
		}

		return tx.Set(p, playerModel)
	})
}

func (r FirestoreGameRepository) UpdateStateAndPlayer(ctx context.Context, state *game.State, player *game.Player) error {
	stateModel := firestoreStateModel{
		UUID:            state.UUID(),
		PlayerUUID:      state.PlayerUUID(),
		GameUUID:        state.GameUUID(),
		GameLevels:      state.GameLevels(),
		Level:           state.Level(),
		Clue:            state.Clue(),
		Completed:       state.Completed(),
		CurrentResponse: state.CurrentResponse(),
	}

	playerModel := firestorePlayerModel{
		UUID:                 player.UUID(),
		Number:               player.Number(),
		GamesStarted:         player.GamesStarted(),
		GamesFinished:        player.GamesFinished(),
		TotalPoints:          player.TotalPoints(),
		CurrentGameStateUUID: player.CurrentGameStateUUID(),
	}

	s := r.client.Doc("game-states/" + state.UUID())
	p := r.client.Doc("players/" + player.UUID())

	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		err := tx.Set(s, stateModel)
		if err != nil {
			return err
		}

		return tx.Set(p, playerModel)
	})
}

func (r FirestoreGameRepository) ReadGames(ctx context.Context, limit, offset int, options ...query.GameOption) ([]*query.Game, error) {
	q := r.client.Collection("games").Query

	for _, option := range options {
		q = q.Where(option.Key, option.Op, option.Value)
	}

	q = q.Offset(offset).Limit(limit).Select("uuid", "title", "description")
	iter := q.Documents(ctx)
	defer iter.Stop()

	// If no games are found return empty non-nil slice.
	results := []*query.Game{}

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return results, err
		}

		model := new(query.Game)

		err = doc.DataTo(model)
		if err != nil {
			return results, err
		}

		results = append(results, &query.Game{
			UUID:        model.UUID,
			Title:       model.Title,
			Description: model.Description,
		})
	}

	return results, nil
}

func (r FirestoreGameRepository) ReadState(ctx context.Context, uuid string) (*query.State, error) {
	q := r.client.Collection("game-states").
		Query.Limit(1).
		Select("currentResponse").
		Where("uuid", "==", uuid)

	iter := q.Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err != nil {
		return nil, err
	}

	st := new(query.State)

	err = doc.DataTo(st)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (r FirestoreGameRepository) ReadPlayer(ctx context.Context, uuid string) (*query.Player, error) {
	q := r.client.Collection("players").
		Query.Limit(1).
		Select("gamesFinished", "totalPoints").
		Where("uuid", "==", uuid)

	iter := q.Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err != nil {
		return nil, err
	}

	p := new(query.Player)

	err = doc.DataTo(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
