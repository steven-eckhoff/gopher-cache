package game

// ResponseKind indicates to clients what kind of information is contained in the response.
type ResponseKind string

const (
	LevelResponse ResponseKind = "level"
	ClueResponse  ResponseKind = "clue"
	EndResponse   ResponseKind = "end"
)

// Response represents the response from the game based on its current state and the player's input.
type Response struct {
	Kind             ResponseKind `json:"kind"`
	LevelTitle       string       `json:"levelTitle"`
	LevelDescription string       `json:"levelDescription"`
	Clue             string       `json:"clue"`
	EndMessage       string       `json:"endMessage"`
}

func newGameEndResponse(msg string) *Response {
	return &Response{
		Kind:       EndResponse,
		EndMessage: msg,
	}
}

func newLevelResponse(title, description string) *Response {
	return &Response{
		Kind:             LevelResponse,
		LevelTitle:       title,
		LevelDescription: description,
	}
}

func newClueResponse(clue string) *Response {
	return &Response{
		Kind: ClueResponse,
		Clue: clue,
	}
}
