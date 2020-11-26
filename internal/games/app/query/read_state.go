package query

import "context"

// ReadStateHandler handles the reading of game states.
type ReadStateHandler struct {
	readModel StateReadModel
}

// NewReadStateHandler creates a new handler.
func NewReadStateHandler(readModel StateReadModel) ReadStateHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return ReadStateHandler{readModel: readModel}
}

// StateReadModel is the interface used for reading State for a client query.
type StateReadModel interface {
	ReadState(ctx context.Context, uuid string) (*State, error)
}

// Handle handles the use case for reading game states.
func (h ReadStateHandler) Handle(ctx context.Context, uuid string) (*State, error) {
	return h.readModel.ReadState(ctx, uuid)
}
