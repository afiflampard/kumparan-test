package authors

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Author struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type AuthorIDName struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`
}

type AuthorInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginAuthorRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthorInputUpdate struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *Author) TableName() string {
	return "authors"
}

func CreateNewAuthor(input AuthorInput) Author {
	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return Author{}
	}
	return Author{
		ID:        uuid.New(),
		Name:      input.Name,
		Password:  string(password),
		Email:     input.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *AuthorInput) ToAuthorUpdate(id uuid.UUID) AuthorInputUpdate {
	return AuthorInputUpdate{
		ID:        id,
		Name:      u.Name,
		Email:     u.Email,
		UpdatedAt: time.Now(),
	}
}
