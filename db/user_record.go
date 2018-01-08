package db

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRecord struct {
	Record
	Timestamps

	Email    string `db:"email"`
	Password []byte `db:"password"`
}

func (u *UserRecord) UpdatePassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	log.Printf("Password of user %s updated", u.Email)
	u.Password = hash
	return nil
}

func (u *UserRecord) CheckPassword(compareWith string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(compareWith))
	return err == nil
}

type UserDB interface {
	Save(UserRecord) (UserRecord, error)
	FindByEmail(string) (UserRecord, error)
	FindByID(int64) (UserRecord, error)
}

func NewUserDB(db DB) UserDB {
	return &userSQLDB{
		db:    db,
		clock: time.Now,
	}
}

type userSQLDB struct {
	db    DB
	clock Clock
}

func (u *userSQLDB) Save(record UserRecord) (UserRecord, error) {
	var err error
	if record.IsPersisted() {
		record.JustUpdated(u.clock)
		sql := "IMPLEMENT ME"
		err = checkResult(u.db.Exec(sql, record.Email, record.UpdatedAt.UTC(), record.ID))
	} else {
		record.Timestamps = JustCreated(u.clock)
		sql := "INSERT INTO users (email, password, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id"
		err = u.db.QueryRow(
			sql,
			record.Email,
			record.Password,
			record.CreatedAt.UTC(),
			record.UpdatedAt.UTC(),
		).Scan(&record.ID)
	}

	return record, err
}

func (u *userSQLDB) FindByEmail(email string) (UserRecord, error) {
	sql := "SELECT * FROM users WHERE email = $1 LIMIT 1"

	var record UserRecord
	err := u.db.QueryRowx(sql, email).StructScan(&record)

	return record, err
}

func (u *userSQLDB) FindByID(id int64) (UserRecord, error) {
	sql := "SELECT * FROM users WHERE id = $1 LIMIT 1"

	var record UserRecord
	err := u.db.QueryRowx(sql, id).StructScan(&record)

	return record, err
}
