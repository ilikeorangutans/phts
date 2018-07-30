package admin

import (
	"time"

	"github.com/ilikeorangutans/phts/db"
	"github.com/pkg/errors"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

type AdminRecord struct {
	db.UserRecord
	db.Timestamps
	UserID int64 `db:"user_id"`
}

type AdminDB interface {
	FindByID(id int64) (*AdminRecord, error)
	FindByEmail(email string) (*AdminRecord, error)
	Save(*AdminRecord) error
}

func NewAdminDB(db db.Queries) AdminDB {
	return &adminSQLDB{
		db:    db,
		clock: time.Now,
		sql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

type adminSQLDB struct {
	db    db.Queries
	clock db.Clock
	sql   sq.StatementBuilderType
}

func (a adminSQLDB) FindByID(id int64) (*AdminRecord, error) {
	return nil, nil
}

func (a adminSQLDB) FindByEmail(email string) (*AdminRecord, error) {
	record := &AdminRecord{UserRecord: db.UserRecord{}}
	query, args, err := a.sql.
		Select("a.user_id, u.*").
		From("admins as a").
		Join("users AS u ON a.user_id = u.id").
		Where(sq.Eq{"u.email": email}).
		ToSql()

	a.db.QueryRowx(query, args...).StructScan(record)
	return record, err
}

func (a adminSQLDB) Save(record *AdminRecord) error {
	record.Timestamps = db.JustCreated(a.clock)
	record.UserID = record.UserRecord.ID

	query, args, err := a.sql.
		Insert("admins").
		Columns("user_id", "updated_at", "created_at").
		Values(record.UserID, record.UpdatedAt.UTC(), record.CreatedAt.UTC()).
		ToSql()
	if err != nil {
		return nil
	}

	_, err = a.db.Exec(query, args...)
	return err
}

func NewAdminService(db db.DB) *AdminService {
	return &AdminService{
		db: db,
	}
}

type AdminService struct {
	db db.DB
}

func (a *AdminService) FindByEmailAndPassword(email, password string) (*AdminRecord, error) {
	admin, err := a.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if !admin.CheckPassword(password) {
		return nil, errors.New("not found")
	}

	return admin, nil
}

func (a *AdminService) FindByEmail(email string) (*AdminRecord, error) {
	return NewAdminDB(a.db).FindByEmail(email)
}

func (a *AdminService) Create(record *AdminRecord) error {
	tx, err := a.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "could not begin transaction")
	}
	userDB := db.NewUserDB(tx)
	if err := userDB.Save(&record.UserRecord); err != nil {
		return errors.Wrap(err, "could not save user")
	}
	record.UserID = record.UserRecord.ID

	adminDB := NewAdminDB(tx)
	if err := adminDB.Save(record); err != nil {
		return errors.Wrap(err, "could not save admin record")
	}

	return errors.Wrap(tx.Commit(), "failed to commit transaction")
}
