package authors

import (
	"context"

	"github.com/afif-musyayyidin/hertz-boilerplate/domain/infra/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthorRepo struct {
	db *sqlx.DB
}

func NewAuthorRepo(db *sqlx.DB) AuthorRepository {
	return &AuthorRepo{db: db}
}

func (r *AuthorRepo) Save(ctx context.Context, u *AuthorInput) (*uuid.UUID, error) {
	newAuthor := CreateNewAuthor(*u)
	_, err := r.db.NamedExecContext(ctx, CreateAuthorQuery, newAuthor)
	if err != nil {
		return nil, err
	}

	return &newAuthor.ID, nil
}

func (r *AuthorRepo) Update(ctx context.Context, u *AuthorInput, id uuid.UUID) (*uuid.UUID, error) {
	author := u.ToAuthorUpdate(id)
	_, err := r.db.NamedExecContext(ctx, UpdateAuthorQuery, author)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *AuthorRepo) FindByID(ctx context.Context, id uuid.UUID) (*Author, error) {
	var u Author
	if err := r.db.GetContext(ctx, &u, FindAuthorByIDQuery, id); err != nil {
		logger.Debug("error find by id", err)
		return nil, err
	}
	return &u, nil
}

func (r *AuthorRepo) FindByIDList(ctx context.Context, idList []uuid.UUID) ([]Author, error) {
	var uList []Author
	query, args, err := sqlx.In(FindAuthorByIDListQuery, idList)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)
	if err := r.db.SelectContext(ctx, &uList, query, args...); err != nil {
		logger.Debug("error find by id list", err)
		return nil, err
	}
	return uList, nil
}

func (r *AuthorRepo) FindIDNameByName(ctx context.Context, name string) ([]*AuthorIDName, error) {
	var idNameList []*AuthorIDName
	if err := r.db.SelectContext(ctx, &idNameList, GetIDAuthorsByNameQuery, name); err != nil {
		logger.Debug("error find by name", err)
		return nil, err
	}
	return idNameList, nil
}
