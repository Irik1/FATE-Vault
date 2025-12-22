package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	ID string `json:"_id" bson:"_id"`

	Username       string `json:"username" bson:"username"`
	HashedPassword string `json:"hashedPassword" bson:"hashedPassword"`
	ProfilePicture string `json:"profilePicture,omitempty" bson:"profilePicture,omitempty"`
	Role           string `json:"role" bson:"role" validate:"oneof=admin user"`

	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func (u *Users) SetPassword(plain string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.HashedPassword = string(hash)
	return nil
}

func (u *Users) CheckPassword(plain string) bool {
	if u.HashedPassword == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(plain))
	return err == nil
}
