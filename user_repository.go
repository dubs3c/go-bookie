package gobookie

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
)

// UserRepositoryCreate - Creates a new user in DB
func (s *Server) UserRepositoryCreate(user *CreateUserRequest) (int, error) {
	var last int = 0
	err := s.DB.QueryRow(context.Background(), "INSERT INTO users(email, password) values($1, $2) RETURNING id", user.Email, user.Password).Scan(&last)
	return last, err
}

// UserRepositoryList - List all users
func (s *Server) UserRepositoryList() ([]*User, error) {

	var users []*User
	err := pgxscan.Select(context.Background(), s.DB, &users, `SELECT id, email, api_token, is_admin, created_at, updated_at FROM users ORDER BY id`)

	return users, err
}

func (s *Server) UserRepositoryExists(email string, hashed_password string) (int, error) {
	var found int = 0
	err := s.DB.QueryRow(context.Background(), "SELECT id FROM users WHERE email = $1 AND password = $2 LIMIT 1;", email, hashed_password).Scan(&found)
	return found, err
}

func (s *Server) UserRepositoryIsAuthenticated(userId int) (bool, error) {
	var id int = 0
	err := s.DB.QueryRow(context.Background(), "SELECT id FROM access_tokens WHERE user_fk = $1 LIMIT 1;", userId).Scan(&id)
	if err != nil {
		return false, err
	}

	if id > 0 {
		return true, err
	}

	return false, err
}

func (s *Server) UserRepositoryIsAuthenticatedByToken(token string) (User, bool) {

	var u User = User{}
	err := s.DB.QueryRow(context.Background(),
		`SELECT u.id, u.email, u.api_token, u.is_admin, u.created_at, u.modified_at, at.token AS access_token, at.created_at AS loggedin_at
	FROM access_tokens AS at
	LEFT JOIN users AS u ON u.id = at.user_fk
	WHERE at.token = $1 LIMIT 1;`, token).Scan(&u.ID, &u.Email, &u.APIToken, &u.IsAdmin, &u.CreatedAt, &u.AccessToken, &u.UpdatedAt)

	if err != nil {
		return u, false
	}

	if u == (User{}) {
		return u, false
	}

	return u, true
}

func (s *Server) UserRepositoryLogin(userId int) (string, error) {
	var token string = ""
	err := s.DB.QueryRow(context.Background(), "INSERT INTO access_tokens(user_fk) values($1) RETURNING token", userId).Scan(&token)
	return token, err
}

func (s *Server) UserRepositoryDeleteAccessToken(token string) error {
	_, err := s.DB.Query(context.Background(), "DELETE FROM access_tokens WHERE token = $1", token)
	return err
}
