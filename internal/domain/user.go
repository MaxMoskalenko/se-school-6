package domain

import "github.com/google/uuid"

type User struct {
	id    uuid.UUID
	email string
}

func NewUser(email string) *User {
	return &User{email: email}
}

func (u *User) WithID(id uuid.UUID) *User {
	u.id = id
	return u
}

func (u *User) WithNewID() *User {
	u.id = uuid.New()
	return u
}

func (u User) ID() uuid.UUID {
	return u.id
}

func (u User) Email() string {
	return u.email
}
