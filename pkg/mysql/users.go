package mysql

import (
	"database/sql"

	"github.com/abrarr21/snippet/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

// Insert Method to add users to the DB
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate method verifies whether user exist with provided email & password or not based on userID
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Get method to fetch a specific user based on userID
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
