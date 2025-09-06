package authors

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthorMutation interface {
	CreateAuthor(ctx context.Context, u *AuthorInput) (*uuid.UUID, error)
	UpdateAuthor(ctx context.Context, u *AuthorInput, id uuid.UUID) (*uuid.UUID, error)
	GetAuthorByID(ctx context.Context, id uuid.UUID) (*Author, error)
	GetAuthorByIDList(ctx context.Context, idList []uuid.UUID) ([]Author, error)
	FindIDNameByName(ctx context.Context, name string) ([]*AuthorIDName, error)
}

type authorMutation struct {
	repo AuthorRepository
	db   *sqlx.DB
}

func NewAuthorMutation(repo AuthorRepository, db *sqlx.DB) AuthorMutation {
	return &authorMutation{repo: repo, db: db}
}

func (m *authorMutation) CreateAuthor(ctx context.Context, u *AuthorInput) (*uuid.UUID, error) {

	if u == nil {
		return nil, ErrInvalidInput
	}
	if u.Name == "" || u.Email == "" {
		return nil, ErrInvalidInput
	}
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	id, err := m.repo.Save(ctx, u)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return id, nil
}

func (m *authorMutation) UpdateAuthor(ctx context.Context, u *AuthorInput, id uuid.UUID) (*uuid.UUID, error) {
	if u == nil {
		return nil, ErrInvalidInput
	}
	if id == uuid.Nil {
		return nil, ErrInvalidInput
	}
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	idResult, err := m.repo.Update(ctx, u, id)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return idResult, nil
}

func (m *authorMutation) FindIDNameByName(ctx context.Context, name string) ([]*AuthorIDName, error) {
	return m.repo.FindIDNameByName(ctx, name)
}

func (m *authorMutation) GetAuthorByID(ctx context.Context, id uuid.UUID) (*Author, error) {
	return m.repo.FindByID(ctx, id)
}

func (m *authorMutation) GetAuthorByIDList(ctx context.Context, idList []uuid.UUID) ([]Author, error) {
	return m.repo.FindByIDList(ctx, idList)
}
