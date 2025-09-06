package authors

const CreateAuthorQuery = `
	INSERT INTO authors (id, name, email)
	VALUES (:id, :name, :email)
`

const UpdateAuthorQuery = `
	UPDATE authors
	SET name = :name, email = :email, updated_at = :updated_at
	WHERE id = :id
`

const FindAuthorByIDQuery = `
	SELECT id, name, email FROM authors WHERE id = $1
`

const FindAuthorByIDListQuery = `
	SELECT id, name, email FROM authors WHERE id IN (?)
`

const GetIDAuthorsByNameQuery = `
	SELECT id, name FROM authors WHERE name = $1
`
