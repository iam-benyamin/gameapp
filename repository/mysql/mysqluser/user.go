package mysqluser

import (
	"context"
	"database/sql"
	"fmt"
	"gameapp/entity"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
	"gameapp/repository/mysql"
	"time"
)

func (d *DB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	const op = "mysql.IsPhoneNumberUnique"

	row := d.conn.Conn().QueryRow(`SELECT * FROM users WHERE phone_number = ?`, phoneNumber)
	_, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return false, richerror.New(op).
			WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryResult).
			WithKind(richerror.KindUnexpected)
	}

	return false, nil
}

func (d *DB) Register(u entity.User) (entity.User, error) {
	res, err := d.conn.Conn().Exec(`INSERT INTO users(name, phone_number, password, role) values(?, ?, ?, ?)`,
		u.Name, u.PhoneNumber, u.Password, u.Role.String())
	if err != nil {
		return entity.User{}, fmt.Errorf("can't execute command: %w", err)
	}

	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}

func (d *DB) GetUserByPhoneNumber(phoneNumber string) (entity.User, error) {
	const op = "mysql.GetUserByPhoneNumber"

	row := d.conn.Conn().QueryRow(`SELECT * FROM users WHERE phone_number = ?`, phoneNumber)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorMsgNotFound).
				WithKind(richerror.KindNotFound)
		}

		return entity.User{}, richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorMsgCantScanQueryResult).
			WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func (d *DB) GetUserByID(ctx context.Context, UserID uint) (entity.User, error) {
	const op = "mysql.GetUserByID"

	row := d.conn.Conn().QueryRowContext(ctx, `SELECT * FROM users WHERE id = ?`, UserID)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorMsgNotFound).
				WithKind(richerror.KindNotFound)
		}

		return entity.User{}, richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorMsgCantScanQueryResult).
			WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func scanUser(scanner mysql.Scanner) (entity.User, error) {
	var user entity.User
	var createdAt time.Time

	var roleStr string

	err := scanner.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAt, &user.Password, &roleStr)

	user.Role = entity.MapToRoleEntity(roleStr)
	return user, err
}

func (d *DB) DoseUserExistByPhoneNumber(phoneNumber string) (entity.User, error) {
	const op = "mysql.GetUserByPhoneNumber"

	row := d.conn.Conn().QueryRow(`SELECT * FROM users WHERE phone_number = ?`, phoneNumber)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorMsgNotFound).
				WithKind(richerror.KindNotFound)
		}

		return entity.User{}, richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorMsgCantScanQueryResult).
			WithKind(richerror.KindUnexpected)
	}

	return user, nil
}
