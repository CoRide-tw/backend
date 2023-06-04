package db

import (
	"context"

	"github.com/CoRide-tw/backend/internal/model"
)

const createUsersTableSQL = `
	CREATE TABLE IF NOT EXISTS users (
	    id SERIAL,
	    name VARCHAR(200) NOT NULL,
		email VARCHAR(200) NOT NULL,
	    google_id VARCHAR(200) NOT NULL UNIQUE,
		picture_url VARCHAR(200),
	    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
	    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
	    deleted_at TIMESTAMP WITH TIME ZONE
	);
`

func initUserTable() error {
	if _, err := DBClient.pgPool.Exec(context.Background(), createUsersTableSQL); err != nil {
		return err
	}
	return nil
}

const getUserSQL = `
	SELECT *
	FROM users
	WHERE id = $1 AND deleted_at IS NULL;
`

func GetUser(id int32) (*model.User, error) {
	var user model.User
	if err := DBClient.pgPool.QueryRow(context.Background(), getUserSQL, id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.GoogleId,
		&user.PictureUrl,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

const createUserSQL = `
	INSERT INTO users (name, email, google_id, picture_url)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (google_id)
	DO UPDATE SET name = $1, email = $2, picture_url = $4, updated_at = NOW(), deleted_at = NULL
	RETURNING id, created_at, updated_at;
`

func UpsertUser(user *model.User) (*model.User, error) {
	if err := DBClient.pgPool.QueryRow(context.Background(), createUserSQL,
		user.Name, user.Email, user.GoogleId, user.PictureUrl).Scan(
		&user.Id, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return user, nil
}

const updateUserSQL = `
	UPDATE users SET
	  email = COALESCE(NULLIF($2, ''), email),
	  updated_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
	RETURNING *;
`

func UpdateUser(id int32, user *model.User) (*model.User, error) {
	var updatedUser model.User
	if err := DBClient.pgPool.QueryRow(context.Background(), updateUserSQL,
		id, user.Email).Scan(
		&updatedUser.Id,
		&updatedUser.Name,
		&updatedUser.Email,
		&updatedUser.GoogleId,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
		&updatedUser.DeletedAt); err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

const deleteUserSQL = `
	UPDATE users SET deleted_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL;
`

func DeleteUser(id int32) error {
	if _, err := DBClient.pgPool.Exec(context.Background(), deleteUserSQL,
		id); err != nil {
		return err
	}
	return nil
}
