package repository

import (
	"context"
	"database/sql"
	"newgolang/proto/pb"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
func (repo *UserRepository) Save(ctx context.Context, user *pb.User) error {
	query := "INSERT INTO users(name, surname, email, password, role) VALUES ($1, $2, $3, $4, $5)"
	_, err := repo.db.ExecContext(ctx, query,
		user.Name, user.Surname, user.Email, user.Password, "STUDENT")
	return err
}

func (repo *UserRepository) GetByID(ctx context.Context, id int64) (*pb.User, error) {
	query := "SELECT id, name, surname, email, password, role, created_at FROM users WHERE id = $1"
	row := repo.db.QueryRowContext(ctx, query, id)

	var user pb.User
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) GetByEmail(ctx context.Context, email string) (*pb.User, error) {
	query := "SELECT id, name, surname, email, password, role FROM users WHERE email = $1"
	row := repo.db.QueryRowContext(ctx, query, email)

	var user pb.User
	if err := row.Scan(&user.Id, &user.Name, &user.Surname, &user.Email, &user.Password, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) GetAll(ctx context.Context) ([]*pb.User, error) {
	query := "SELECT id, name, surname, email, password, role, created_at FROM users"
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		var user pb.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Surname, &user.Email, &user.Password, &user.Role, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepository) UpdateByEmail(ctx context.Context, email string, user *pb.User) error {
	query := "UPDATE users SET name = $1, surname = $2 WHERE email = $3"
	_, err := repo.db.ExecContext(ctx, query, user.Name, user.Surname, email)
	return err
}

func (repo *UserRepository) DeleteByID(ctx context.Context, id int64) error {
	query := "DELETE FROM users where id = $1"
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *UserRepository) DeleteByEmail(ctx context.Context, email string) error {
	query := "DELETE FROM users where email = $1"
	_, err := repo.db.ExecContext(ctx, query, email)
	return err
}
