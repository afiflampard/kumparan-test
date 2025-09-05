package articles

const CreateArticleQuery = `
	INSERT INTO articles (id, title, body, author_id)
	VALUES (:id, :title, :body, :author_id)
`

const UpdateArticleQuery = `
	UPDATE articles
	SET title = :title, body = :body, updated_at = :updated_at
	WHERE id = :id
`

const FindArticleByIDQuery = `
	SELECT * FROM articles WHERE id = $1
`

const FindAllArticleByAuthorIDQuery = `
	SELECT * FROM articles WHERE author_id = $1
`

const FindAllArticleWithAuthorByAuthorIDQuery = `
	SELECT * FROM articles WHERE author_id = $1
`
