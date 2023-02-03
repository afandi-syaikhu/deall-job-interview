package repository

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"time"

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

func (_u *User) Create(ctx context.Context, data model.User) error {
	query := `
		insert into users
			(username, password, role, created_at, updated_at)
		values
			($1, $2, $3, $4, $5);
	`
	now := time.Now()
	passHash := md5.Sum([]byte(data.Password))
	passMd5 := hex.EncodeToString(passHash[:])
	_, err := _u.DB.ExecContext(ctx, query, data.Username, passMd5, data.Role, now, now)
	if err != nil {
		return err
	}

	return nil
}

func (_u *User) Fetch(ctx context.Context, data model.PaginationRequest) ([]*model.User, error) {
	query := `
		select id, username, password, role, created_at, updated_at
		from 
			users
		order by id
		limit $1
		offset $2
	`
	offset := (data.Page - 1) * data.Limit
	rows, err := _u.DB.QueryContext(ctx, query, data.Limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*model.User{}
	for rows.Next() {
		user := &model.User{}
		err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
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

func (_u *User) FindById(ctx context.Context, id int64) (*model.User, error) {
	query := `
		select id, username, password, role, created_at, updated_at
		from 
			users
		where 
			id = $1
	`
	rows, err := _u.DB.QueryContext(ctx, query, id)
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

func (_u *User) CountByRole(ctx context.Context, role model.Role) (int, error) {
	count := 0
	query := `
		select 
			count(*) 
		from
			users
		where
			role = $1
	`
	rows, err := _u.DB.QueryContext(ctx, query, role)
	if err != nil {
		return count, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (_u *User) IsExistByUsername(ctx context.Context, username string) (bool, error) {
	query := `
		select exists (
			select 1
			from users
			where username=$1
		) as exists
	`
	row := _u.DB.QueryRowContext(ctx, query, username)

	isExist := false
	err := row.Scan(&isExist)
	if err != nil {
		return isExist, err
	}

	return isExist, nil
}

func (_u *User) IsExistByUsernameAndIdNot(ctx context.Context, username string, id int64) (bool, error) {
	query := `
		select exists (
			select 1
			from users
			where username=$1
				and id != $2
		) as exists
	`
	row := _u.DB.QueryRowContext(ctx, query, username, id)

	isExist := false
	err := row.Scan(&isExist)
	if err != nil {
		return isExist, err
	}

	return isExist, nil
}

func (_u *User) Update(ctx context.Context, data model.User) error {
	query := `
		update 
			users 
		set
			username = $2,
			password = $3,
			role = $4,
			updated_at = $5
		where 
			id = $1
	`
	now := time.Now()
	passHash := md5.Sum([]byte(data.Password))
	passMd5 := hex.EncodeToString(passHash[:])
	_, err := _u.DB.ExecContext(ctx, query, data.ID, data.Username, passMd5, data.Role, now)
	if err != nil {
		return err
	}

	return nil
}

func (_u *User) DeleteById(ctx context.Context, id int64) error {
	query := `
		delete
		from 
			users
		where 
			id = $1
	`
	result, err := _u.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	ra, _ := result.RowsAffected()
	if ra == 0 {
		return sql.ErrNoRows
	}

	return nil
}
