package entity

import (
	"time"
)

type UserProfile struct {
	ID        uint
	UserID    uint
	Bio       string
	AvatarURL string
	UpdatedAt time.Time
}

func NewUserProfile(userID uint, bio, avatarURL string) *UserProfile {
	return &UserProfile{
		UserID:    userID,
		Bio:       bio,
		AvatarURL: avatarURL,
		UpdatedAt: time.Now(),
	}
}

func (p *UserProfile) UpdateBio(newBio string) {
	p.Bio = newBio
	p.UpdatedAt = time.Now()
}
