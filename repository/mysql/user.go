package mysql

import (
	"database/sql"
	"fmt"
	"gameapp/entity"
)

func (d MySQLDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	user := entity.User{}
	//var createdAt time.Time
	var createdAt []uint8

	row := d.db.QueryRow(`SELECT * FROM users WHERE phone_number = ?`, phoneNumber)

	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return false, fmt.Errorf("unexpectd err: %w", err)
	}

	return false, nil
}

func (d MySQLDB) Register(u entity.User) (entity.User, error) {
	res, err := d.db.Exec(`INSERT INTO users(name, phone_number, password) values(?, ?, ?)`, u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("can't execute command: %w", err)
	}

	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}

func (d MySQLDB) GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error) {
	user := entity.User{}
	//var createdAt time.Time
	var createdAt []uint8

	row := d.db.QueryRow(`SELECT * FROM users WHERE phone_number = ?`, phoneNumber)

	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, false, nil
		}

		return entity.User{}, false, fmt.Errorf("unexpectd err: %w", err)
	}

	return user, true, nil
}
