package postgres

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"
	"userland/api/helper"
	"userland/store"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) store.UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) ResetPassword(ctx context.Context, uid string, u store.User) error {
	var tx *sql.Tx
	tx, err := us.db.Begin()
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return errors.New("failed to update password")
	}
	defer tx.Rollback()

	var updatePassword *sql.Stmt
	updatePassword, err = tx.Prepare("UPDATE person SET Password = $1 WHERE id = $2")
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error().Err(err).Msg("error occur rolling back changer")
		}
		return errors.New("failed to update password")
	}
	defer updatePassword.Close()

	var result sql.Result
	result, err = updatePassword.Exec(u.Password, uid)

	rowsAff, _ := result.RowsAffected()

	if err != nil || rowsAff != 1 {
		log.Error().Err(err).Msg("error updating user password")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error().Err(err).Msg("error occur rolling back changer")
		}
		return errors.New("failed to update password")
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("error committing changes")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error().Err(err).Msg("error occur rolling back changer")
		}
		return errors.New("error committing changes")
	}

	return nil
}

func (us *UserStore) GetUserId(ctx context.Context, u store.User) (string, error) {
	psqlStatement := `SELECT id, password FROM PERSON WHERE EMAIL = $1`

	res := us.db.QueryRow(psqlStatement, u.Email)
	var id string
	var password string

	err := res.Scan(&id, &password)

	if err != nil {
		return "", errors.New("unable to find the user")
	}
	if !helper.ComparePasswordHash(u.Password, password) {
		return "", errors.New("password incorrect")
	}
	return id, nil
}

func (us *UserStore) GetUserCode(ctx context.Context, u store.User) (int, error) {
	psqlStatement := `SELECT ver_code FROM email_ver WHERE EMAIL = $1`

	res := us.db.QueryRow(psqlStatement, u.Email)
	var code int

	err := res.Scan(&code)

	if err != nil {
		return 0, errors.New("unable to find the user")
	}
	return code, nil

}

func (us *UserStore) EmailExist(ctx context.Context, u store.User) error {
	psqlStatement := `SELECT email FROM PERSON WHERE EMAIL = $1`

	res := us.db.QueryRow(psqlStatement, u.Email)
	var email string

	switch err := res.Scan(&email); err {
	case sql.ErrNoRows:
		return nil
	case nil:
		return errors.New("email is registed")
	default:
		return errors.New("server error")
	}
}

func (us *UserStore) RegisterUser(ctx context.Context, u store.User) error {
	var tx *sql.Tx
	tx, err := us.db.Begin()

	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return errors.New("failed to register user")
	}
	defer tx.Rollback()

	if us.EmailExist(ctx, u) != nil {
		return errors.New("email is used")
	}

	// Prepare statement for inserting new user
	var inserUserStmt *sql.Stmt
	inserUserStmt, err = tx.Prepare(`INSERT INTO person 
									(fullname, email, password, created_at, id, is_active) 
									values ($1, $2, $3, $4, $5, $6);`)
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error().Err(err).Msg("error occur rolling back changer")
		}
		return errors.New("failed to register user")
	}
	defer inserUserStmt.Close()

	userId, err := uuid.NewV4()
	if err != nil {
		log.Error().Err(err).Msg("failed to generate user id")
		return errors.New("failed to register user")
	}

	var result sql.Result
	result, err = inserUserStmt.Exec(u.Fullname, u.Email, u.Password, time.Now().UTC(), userId.String(), 0)

	rowsAff, _ := result.RowsAffected()

	if err != nil || rowsAff != 1 {
		log.Error().Err(err).Msg("error inserting new user")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error().Err(err).Msg("error occur rolling back changer")
		}
		return errors.New("failed to register user")
	}

	rand.Seed(time.Now().UnixNano())
	rn := rand.Intn(100000)

	// Prepare statement for inserting new verfication code
	var insertEmailVerStmt *sql.Stmt
	insertEmailVerStmt, err = tx.Prepare(`INSERT INTO email_ver 
										(fullname, email, ver_code) 
										VALUES ($1, $2, $3);`)
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error().Err(err).Msg("error occur rolling back changer")
		}
		return errors.New("failed to register user")
	}
	defer insertEmailVerStmt.Close()
	result, err = insertEmailVerStmt.Exec(u.Fullname, u.Email, rn)
	rowsAff, _ = result.RowsAffected()

	if err != nil || rowsAff != 1 {
		log.Error().Err(err).Msg("error inserting email varification")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error().Err(err).Msg("error occur rolling back changer")
		}
		return errors.New("failed to register user")
	}

	go helper.SendEmail(u.Email, rn)

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("error committing changes")
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error().Err(err).Msg("error occur rolling back changer")
		}
		return errors.New("failed to register user")
	}

	return nil
}

func (us *UserStore) ValidateCode(ctx context.Context, u store.User) error {
	code, err := us.GetUserCode(ctx, u)
	if err != nil {
		return err
	}

	if code == u.VerCode {
		var tx *sql.Tx
		tx, err = us.db.Begin()
		if err != nil {
			log.Error().Err(err).Msg("failed to begin transaction")
			return errors.New("failed to update password")
		}
		defer tx.Rollback()

		var updateUser *sql.Stmt
		updateUser, err := tx.Prepare(`UPDATE PERSON SET is_active = 1 WHERE email = $1`)
		if err != nil {
			log.Error().Err(err).Msg("error preparing statement")
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(err).Msg("error occur rolling back changer")
			}
			return errors.New("failed to update user state")
		}
		defer updateUser.Close()

		var result sql.Result
		result, err = updateUser.Exec(u.Email)

		rowsAff, _ := result.RowsAffected()

		if err != nil || rowsAff != 1 {
			log.Error().Err(err).Msg("error updating user state")
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(err).Msg("error occur rolling back changer")
			}
			return errors.New("failed to update user state")
		}

		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("error committing changes")
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(err).Msg("error occur rolling back changer")
			}
			return errors.New("error committing changes")
		}
		return nil

	} else {
		return errors.New("incorrect verification code")
	}
}
