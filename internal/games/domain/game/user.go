package game

import "errors"

// User holds all user info.
type User struct {
	uuid   string
	number string
}

func (u User) UUID() string   { return u.uuid }
func (u User) Number() string { return u.number }

// NewUser creates a new user.
func NewUser(uuid, number string) (User, error) {
	if uuid == "" {
		return User{}, errors.New("user has no uuid")
	}

	if number == "" {
		return User{}, errors.New("user has no number")
	}

	return User{
		uuid:   uuid,
		number: number,
	}, nil
}
