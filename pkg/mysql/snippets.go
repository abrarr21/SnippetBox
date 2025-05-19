package mysql

import (
	"database/sql"

	"github.com/abrarr21/snippet/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

// This will Insert a new snippet into the DB.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}

// This will return a specific snippet based on the ID
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// This will return the 10 most recently created snippet
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
