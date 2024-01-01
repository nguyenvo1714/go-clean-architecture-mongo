package auth

import "context"

type Profile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ProfileUseCase interface {
	GetProfileByID(ctx context.Context, id string) (*Profile, error)
}
