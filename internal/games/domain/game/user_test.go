package game

import "github.com/google/uuid"

func newTestUser() User {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	return User{
		uuid:   id.String(),
		number: "15734497033",
	}
}
