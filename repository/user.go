package repository

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"

	"github.com/afandi-syaikhu/deall-job-interview/model"
)

type User struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &User{
		DB: db,
	}
}

func (_u *User) FindByCredential(ctx context.Context, data model.User) (*model.User, error) {
	query := `
			select id, username, password, role, created_at, updated_at
			from 
			    users
			where 
			    username = $1
				and password = $2
	`
	passHash := md5.Sum([]byte(data.Password))
	passMd5 := hex.EncodeToString(passHash[:])
	rows, err := _u.DB.QueryContext(ctx, query, data.Username, passMd5)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var user *model.User
	for rows.Next() {
		user = &model.User{}
		err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}

	if user == nil {
		return nil, sql.ErrNoRows
	}

	return user, nil
}
