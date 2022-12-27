// Code generated by sqlgen. DO NOT EDIT.

package database

import (
	"context"

	"database/sql"
)

// *** select user ***

type User struct {
	ID      string
	Name    string
	Age     uint
	Comment sql.NullString
}

func (q Queries) SelectUser(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, `-- select user
select *
from users
where id = ?`, id)

	var out User
	err := row.Scan(
		&out.ID,
		&out.Name,
		&out.Age,
		&out.Comment,
	)

	return out, err
}

// *** select all users ***

func (q Queries) SelectAllUsers(ctx context.Context, f func(out User) error) error {
	rows, err := q.db.QueryContext(ctx, `-- select all users
select *
from users`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var out User
	for rows.Next() {
		if err := rows.Scan(
			&out.ID,
			&out.Name,
			&out.Age,
			&out.Comment,
		); err != nil {
			return err
		}

		if err := f(out); err != nil {
			return err
		}
	}

	if err := rows.Close(); err != nil {
		return err
	}

	return rows.Err()
}

// *** add new user ***

func (q Queries) AddNewUser(ctx context.Context, args User) error {
	_, err := q.db.ExecContext(ctx, `-- add new user
insert into users
  (id, name, age, comment)
values (?, ?, ?, ?)`,
		args.ID,
		args.Name,
		args.Age,
		args.Comment)

	return err
}

// *** update user ***

type UpdateUserParams struct {
	Name    string
	Age     uint
	Comment sql.NullString
	ID      string
}

func (q Queries) UpdateUser(ctx context.Context, args UpdateUserParams) error {
	result, err := q.db.ExecContext(ctx, `-- update user
update users
  set name = ?, age = ?, comment =?
where id = ?`,
		args.Name,
		args.Age,
		args.Comment,
		args.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNoRows
	}

	return nil
}