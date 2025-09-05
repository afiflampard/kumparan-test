package authors

import (
	"context"

	"github.com/google/uuid"
)

type AuthorMutation interface {
	CreateAuthor(ctx context.Context, u *AuthorInput) (*uuid.UUID, error)
	UpdateAuthor(ctx context.Context, u *AuthorInput, id uuid.UUID) (*uuid.UUID, error)
	GetAuthorByID(ctx context.Context, id uuid.UUID) (*Author, error)
	GetAuthorByIDList(ctx context.Context, idList []uuid.UUID) ([]Author, error)
	FindIDNameByName(ctx context.Context, name string) ([]*AuthorIDName, error)

	Commit(ctx context.Context) error
	Cancel(ctx context.Context)
}

type authorMutation struct {
	repo AuthorRepository
}

func NewAuthorMutation(repo AuthorRepository) AuthorMutation {
	return &authorMutation{repo: repo}
}

func (m *authorMutation) CreateAuthor(ctx context.Context, u *AuthorInput) (*uuid.UUID, error) {
	if u == nil {
		return nil, ErrInvalidInput
	}
	if u.Name == "" || u.Email == "" {
		return nil, ErrInvalidInput
	}
	return m.repo.Save(ctx, u)
}

func (m *authorMutation) UpdateAuthor(ctx context.Context, u *AuthorInput, id uuid.UUID) (*uuid.UUID, error) {
	if u == nil {
		return nil, ErrInvalidInput
	}
	if id == uuid.Nil {
		return nil, ErrInvalidInput
	}
	return m.repo.Update(ctx, u, id)
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

func (m *authorMutation) Commit(ctx context.Context) error {
	return m.repo.Commit(ctx)
}

func (m *authorMutation) Cancel(ctx context.Context) {
	m.repo.Cancel(ctx)
}
