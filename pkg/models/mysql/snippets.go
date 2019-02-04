package mysql

import (
	"database/sql"

	"github.com/jcorry/snippetbox/pkg/models"
)

// SnippetModel wraps DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// Insert a new snippet into the DB
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}

// Get a snippet by ID
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// Latest snippet retrieved and returned.
func (m *SnippetModel) Latest() (*models.Snippet, error) {
	return nil, nil
}
