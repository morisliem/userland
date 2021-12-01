package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

type LoginLog struct {
	Username   string
	Ip_address string
	Created_at time.Time
}

type LogStoreInterface interface {
	SetLoginLog(ll LoginLog) error
}

func NewLoginStore(db *sql.DB) LogStoreInterface {
	return LogStore{
		db: db,
	}
}

type LogStore struct {
	db *sql.DB
}

func (ls LogStore) SetLoginLog(ll LoginLog) error {
	var tx *sql.Tx
	tx, err := ls.db.Begin()

	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return errors.New("failed to set user log audit")
	}

	defer func() {
		if rollBackErr := tx.Rollback(); rollBackErr == nil {
			log.Error().Err(err).Msg("rolling back changes")
		}
	}()

	var setLogAudit *sql.Stmt
	setLogAudit, err = tx.Prepare(`INSERT into "Login_audit" (username, ip_address, created_at) VALUES ($1,$2,$3)`)
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		return errors.New("failed to set user log audit")
	}
	defer setLogAudit.Close()

	result, err := setLogAudit.Exec(ll.Username, ll.Ip_address, ll.Created_at)
	rowsAff, _ := result.RowsAffected()

	if err != nil || rowsAff != 1 {
		log.Error().Err(err).Msg("error inserting user session")
		return errors.New("failed to set user log audit")
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("error committing changes")
		return errors.New("failed to set user log audit")
	}
	return nil
}
