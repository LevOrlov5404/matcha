package models

import (
	"golang.org/x/crypto/bcrypt"
)

type (
	UserToCreate struct {
		Email     string `json:"email" binding:"required"`
		Username  string `json:"username" binding:"required"`
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}
	UserToSignIn struct {
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		Fingerprint string `json:"fingerprint" binding:"required"`
	}
	User struct {
		ID               uint64 `json:"id" binding:"required" db:"id"`
		Email            string `json:"email" binding:"required" db:"email"`
		Username         string `json:"username" binding:"required" db:"username"`
		FirstName        string `json:"firstName" binding:"required" db:"first_name"`
		LastName         string `json:"lastName" binding:"required" db:"last_name"`
		Password         string `json:"-" db:"password"`
		IsEmailConfirmed bool   `json:"isEmailConfirmed" db:"is_email_confirmed"`
	}
	UserPassword struct {
		ID       uint64 `json:"id" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	UserPasswordToChange struct {
		ID          uint64 `json:"id" binding:"required"`
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}
	UserProfile struct {
		ID                uint64   `json:"id" binding:"required"`
		Email             string   `json:"email" binding:"required"`
		Username          string   `json:"username" binding:"required"`
		FirstName         string   `json:"firstName" binding:"required"`
		LastName          string   `json:"lastName" binding:"required"`
		IsEmailConfirmed  bool     `json:"isEmailConfirmed"`
		Gender            int      `json:"gender"`
		SexualPreferences int      `json:"sexualPreferences"`
		Biography         string   `json:"biography"`
		Tags              []string `json:"tags"`
		AvatarURL         string   `json:"avatarURL"`
		PicturesURL       []string `json:"picturesURL"`
		LikesNum          int      `json:"likesNum"`
		GPSPosition       string   `json:"gpsPosition"`
	}
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
