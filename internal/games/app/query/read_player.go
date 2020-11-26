package query

import "context"

// ReadPlayerHandler handles reading a player.
type ReadPlayerHandler struct {
	readModel PlayerReadModel
}

// NewReadPlayerHandler creates a new handler.
func NewReadPlayerHandler(readModel PlayerReadModel) ReadPlayerHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return ReadPlayerHandler{readModel: readModel}
}

// PlayerReadModel is the interface used for reading a Player for a client query.
type PlayerReadModel interface {
	ReadPlayer(ctx context.Context, uuid string) (*Player, error)
}

// Handle handles the use case for reading a player.
func (h ReadPlayerHandler) Handle(ctx context.Context, uuid string) (*Player, error) {
	return h.readModel.ReadPlayer(ctx, uuid)
}
