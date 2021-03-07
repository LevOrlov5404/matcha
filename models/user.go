package models

type (
	UserToCreate struct {
		Name     string `json:"name" binding:"required"`
		Surname  string `json:"surname" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	UserToSignIn struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	UserToGet struct {
		ID      int64  `json:"id" db:"id"`
		Name    string `json:"name" db:"name"`
		Surname string `json:"surname" db:"surname"`
		Email   string `json:"email" db:"email"`
	}
)
