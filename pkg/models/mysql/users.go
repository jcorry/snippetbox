package mysql

import (
	"database/sql"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-sql-driver/mysql"
	"github.com/jcorry/snippetbox/pkg/models"
)

// UserModel struct
type UserModel struct {
	DB *sql.DB
}

// Insert a new user
func (u *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
		VALUES (?, ?, ?, UTC_TIMESTAMP())`

	_, err = u.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "uk_users_email") {
				return models.ErrDuplicateEmail
			}
		}
	}
	return err
}

// Authenticate the user given the email/password pair
func (u *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	row := u.DB.QueryRow("SELECT id, hashed_password FROM users WHERE email = ?", email)
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	return id, nil
}

// Get a single user by ID
func (u *UserModel) Get(id int) (*models.User, error) {
	s := &models.User{}

	stmt := `SELECT id, name, email, created FROM users WHERE id = ?`

	err := u.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Name, &s.Email, &s.Created)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}
