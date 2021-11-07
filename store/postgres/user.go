package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"userland/api/helper"
	"userland/store"

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

func (us *UserStore) UpdatePassword(ctx context.Context, uid string, u store.User) error {
	var tx *sql.Tx
	tx, err := us.db.Begin()
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return errors.New("failed to update password")
	}
	defer tx.Rollback()
	defer func() {
		if rollBackErr := tx.Rollback(); rollBackErr == nil {
			log.Error().Err(err).Msg("rolling back changes")
		}
	}()

	// Updating user password
	var updatePassword *sql.Stmt
	updatePassword, err = tx.Prepare("UPDATE person SET Password = $1 WHERE id = $2")
	if err != nil {
		log.Error().Err(err).Msg("error preparing update password statement")
		return errors.New("failed to update password")
	}
	defer updatePassword.Close()

	var result sql.Result
	result, err = updatePassword.Exec(u.Password, uid)

	rowsAff, _ := result.RowsAffected()

	if err != nil || rowsAff != 1 {
		log.Error().Err(err).Msg("error updating user password")
		return errors.New("failed to update password")
	}

	// Adding new password to user_password table
	var insertNewPassword *sql.Stmt
	insertNewPassword, err = tx.Prepare(`INSERT INTO user_password (id, password, created_at) values ($1, $2, $3)`)
	if err != nil {
		log.Error().Err(err).Msg("error preparing insert new password statement")
		return errors.New("failed to update password")
	}
	defer insertNewPassword.Close()

	result, err = insertNewPassword.Exec(uid, u.Password, time.Now())

	rowsAff, _ = result.RowsAffected()

	if err != nil || rowsAff != 1 {
		log.Error().Err(err).Msg("error updating user password")
		return errors.New("failed to update password")
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("error committing changes")
		return errors.New("error committing changes")
	}

	return nil
}

func (us *UserStore) GetUserId(ctx context.Context, u store.User) (string, error) {
	psqlStatement := `SELECT id, password FROM PERSON WHERE EMAIL = $1`

	res := us.db.QueryRowContext(ctx, psqlStatement, u.Email)
	var id string
	var password string

	err := res.Scan(&id, &password)

	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return "", errors.New("unable to find the user")
	}
	if !helper.ComparePasswordHash(u.Password, password) {
		log.Error().Err(err).Msg("password incorrect")
		return "", errors.New("password incorrect")
	}
	return id, nil
}

func (us *UserStore) GetUserid(ctx context.Context, email string) (string, error) {
	psqlStatement := `SELECT id FROM PERSON WHERE EMAIL = $1`

	res := us.db.QueryRowContext(ctx, psqlStatement, email)
	var userId string

	err := res.Scan(&userId)
	if err != nil {
		return "", errors.New("unable to find the user")
	}
	return userId, nil
}

func (us *UserStore) GetUserCode(ctx context.Context, u store.User) (int, error) {
	psqlStatement := `SELECT ver_code FROM email_ver WHERE EMAIL = $1`

	res := us.db.QueryRowContext(ctx, psqlStatement, u.Email)
	var code int

	err := res.Scan(&code)

	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return 0, errors.New("unable to find the user")
	}
	return code, nil
}

func (us *UserStore) EmailActive(ctx context.Context, u store.User) (bool, error) {
	psqlStatement := `SELECT is_active FROM Person WHERE Email = $1`

	res := us.db.QueryRowContext(ctx, psqlStatement, u.Email)
	var is_active bool

	err := res.Scan(&is_active)

	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return false, errors.New("unable to find the user")
	}
	return is_active, nil
}

func (us *UserStore) GetPassword(ctx context.Context, uid string) (string, error) {
	psqlStatement := `SELECT password from person where id = $1`

	res := us.db.QueryRowContext(ctx, psqlStatement, uid)
	var pwd string

	err := res.Scan(&pwd)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return "", errors.New("unable to find user")
	}
	return pwd, nil
}

func (us *UserStore) GetPasswords(ctx context.Context, uid string) ([]string, error) {
	psqlStatement := `SELECT password FROM user_password WHERE id = $1 order by created_at desc limit 3`

	res, err := us.db.QueryContext(ctx, psqlStatement, uid)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return nil, errors.New("unable to find user")
	}

	defer res.Close()

	var pwd []string

	for res.Next() {
		var tmp string
		err = res.Scan(&tmp)
		if err != nil {
			log.Error().Err(err).Msg("error getting the password")
			return nil, errors.New("unable to find user")
		}
		pwd = append(pwd, tmp)
	}

	return pwd, nil
}

func (us *UserStore) EmailExist(ctx context.Context, u store.User) error {
	psqlStatement := `SELECT email FROM PERSON WHERE EMAIL = $1`

	res := us.db.QueryRowContext(ctx, psqlStatement, u.Email)
	var email string

	err := res.Scan(&email)

	if err != nil {
		return errors.New("unable to find the user")
	} else {
		return nil
	}
}

func (us *UserStore) RegisterUser(ctx context.Context, u store.User, rn int) error {
	var tx *sql.Tx
	tx, err := us.db.Begin()

	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return errors.New("failed to register user")
	}
	defer tx.Rollback()
	defer func() {
		if rollBackErr := tx.Rollback(); rollBackErr == nil {
			log.Error().Err(err).Msg("rolling back changes")
		}
	}()

	if us.EmailExist(ctx, u) == nil {
		return errors.New("email is used")
	}

	// Prepare statement for inserting new user
	var inserUserStmt *sql.Stmt
	inserUserStmt, err = tx.Prepare(`INSERT INTO person 
									(fullname, email, password, created_at, id, is_active) 
									values ($1, $2, $3, $4, $5, $6);`)
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		return errors.New("failed to register user")
	}
	defer inserUserStmt.Close()

	var result sql.Result
	result, err = inserUserStmt.Exec(u.Fullname, u.Email, u.Password, time.Now().UTC(), u.Id, false)

	rowsAff, _ := result.RowsAffected()

	if err != nil || rowsAff != 1 {
		log.Error().Err(err).Msg("error inserting new user")
		return errors.New("failed to register user")
	}

	// Prepare statement for inserting new password
	var insertNewPassword *sql.Stmt
	insertNewPassword, err = tx.Prepare(`INSERT INTO user_password 
										(id, password, created_at) 
										VALUES($1, $2, $3)`)
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		return errors.New("failed to register user")
	}
	defer insertNewPassword.Close()
	result, err = insertNewPassword.Exec(u.Id, u.Password, time.Now())
	rowsAff, _ = result.RowsAffected()

	if err != nil || rowsAff != 1 {
		log.Error().Err(err).Msg("error inserting new user")
		return errors.New("failed to register user")
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("error committing changes")
		return errors.New("failed to register user")
	}
	return nil
}

func (us *UserStore) ValidateCode(ctx context.Context, u store.User) error {
	var tx *sql.Tx
	tx, err := us.db.Begin()

	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return errors.New("failed to update password")
	}
	defer tx.Rollback()
	defer func() {
		if rollBackErr := tx.Rollback(); rollBackErr == nil {
			log.Error().Err(err).Msg("rolling back changes")
		}
	}()

	var updateUser *sql.Stmt
	updateUser, err = tx.Prepare(`UPDATE PERSON SET is_active = $3, email = $2 WHERE id = $1`)
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		return errors.New("failed to update user state")
	}
	defer updateUser.Close()

	var result sql.Result
	result, err = updateUser.Exec(u.Id, u.Email, true)
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		return errors.New("failed to update user state")
	}

	rowsAff, _ := result.RowsAffected()

	if err != nil || rowsAff != 1 {
		log.Error().Err(err).Msg("error updating user state")
		return errors.New("failed to update user state")
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("error committing changes")
		return errors.New("error committing changes")
	}
	return nil
}

func (us *UserStore) GetUserDetail(ctx context.Context, uid string) (store.User, error) {
	var response store.User
	psqlstatement := `SELECT * FROM person WHERE id = $1`

	res := us.db.QueryRowContext(ctx, psqlstatement, uid)
	var id, fullname, email, password string
	var location, bio, picture, web sql.NullString
	var created_at time.Time
	var is_active bool

	err := res.Scan(&id, &fullname, &email, &password, &location, &bio, &web, &picture, &created_at, &is_active)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return response, errors.New("unable to find the user")
	}

	loc := ""
	bioo := ""
	webs := ""
	pict := ""

	response.Id = id
	response.Fullname = fullname

	if !location.Valid {
		response.Location = loc
	} else {
		var tmp string
		err := us.db.QueryRowContext(ctx, `SELECT location from PERSON where id = $1`, id).Scan(&tmp)
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
			return response, errors.New("unable to find the user")
		}
		response.Location = tmp
	}

	if !bio.Valid {
		response.Bio = bioo
	} else {
		var tmp string
		err := us.db.QueryRowContext(ctx, `SELECT bio from PERSON where id = $1`, id).Scan(&tmp)
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
			return response, errors.New("unable to find the user")
		}
		response.Bio = tmp
	}

	if !web.Valid {
		response.Web = webs
	} else {
		var tmp string
		err := us.db.QueryRowContext(ctx, `SELECT web from PERSON where id = $1`, id).Scan(&tmp)
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
			return response, errors.New("unable to find the user")
		}
		response.Web = tmp
	}

	if !picture.Valid {
		response.Picture = pict
	} else {
		var tmp string
		err := us.db.QueryRowContext(ctx, `SELECT picture from PERSON where id = $1`, id).Scan(&tmp)
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
			return response, errors.New("unable to find the user")
		}
		response.Picture = tmp
	}

	response.Created_at = created_at

	return response, nil
}

func (us *UserStore) UpdateUserDetail(ctx context.Context, u store.User, uid string) error {
	var tx *sql.Tx
	tx, err := us.db.Begin()
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return errors.New("failed to update user detail")
	}
	defer tx.Rollback()
	defer func() {
		if rollBackErr := tx.Rollback(); rollBackErr == nil {
			log.Error().Err(err).Msg("rolling back changes")
		}
	}()

	var updateUserDetail *sql.Stmt
	updateUserDetail, err = tx.Prepare(`UPDATE person 
										SET Fullname = $2, Location = $3, Bio = $4, Web = $5 
										WHERE id = $1`)

	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		return errors.New("failed to update user detail")
	}
	defer updateUserDetail.Close()

	var result sql.Result
	result, err = updateUserDetail.Exec(uid, u.Fullname, u.Location, u.Bio, u.Web)

	rowsAff, _ := result.RowsAffected()

	if err != nil || rowsAff != 1 {
		log.Error().Err(err).Msg("error updating user user detail")
		return errors.New("failed to update user detail")
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("error committing changes")
		return errors.New("error committing changes")
	}

	return nil
}

func (us *UserStore) GetUserEmail(ctx context.Context, uid string) (store.User, error) {
	var email string
	var response store.User
	psqlstatement := `SELECT email FROM person WHERE id = $1`

	err := us.db.QueryRowContext(ctx, psqlstatement, uid).Scan(&email)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return response, errors.New("unable to find the user")
	}

	response.Email = email

	return response, nil
}

func (us *UserStore) DeleteAccount(ctx context.Context, uid string) error {
	var tx *sql.Tx
	tx, err := us.db.Begin()

	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return errors.New("failed to delete account")
	}
	defer tx.Rollback()
	defer func() {
		if rollBackErr := tx.Rollback(); rollBackErr == nil {
			log.Error().Err(err).Msg("rolling back changes")
		}
	}()

	var updateAccountState *sql.Stmt
	updateAccountState, err = tx.Prepare(`UPDATE person SET is_active = $2 WHERE id = $1`)
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		return errors.New("failed to delete account")
	}
	defer updateAccountState.Close()

	result, err := updateAccountState.Exec(uid, false)
	rowsAff, _ := result.RowsAffected()

	if err != nil || rowsAff != 1 {
		log.Error().Err(err).Msg("error inserting email varification")
		return errors.New("failed to delete account")
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("error committing changes")
		return errors.New("failed to delete account")
	}

	return nil
}

func (us *UserStore) GetUserProfilePicture(ctx context.Context, uid string) (string, error) {
	psqlStatement := `SELECT picture from PERSON WHERE id = $1`

	res := us.db.QueryRowContext(ctx, psqlStatement, uid)
	var pictureName string
	err := res.Scan(&pictureName)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return "", errors.New("failed to get user picture")
	}

	return pictureName, nil
}

func (us *UserStore) DeleteUserPicture(ctx context.Context, uid string) error {
	var tx *sql.Tx
	tx, err := us.db.Begin()

	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return errors.New("failed to remove user picture")
	}

	defer tx.Rollback()
	defer func() {
		if rollBackErr := tx.Rollback(); rollBackErr == nil {
			log.Error().Err(err).Msg("rolling back changes")
		}
	}()

	var setPicture *sql.Stmt
	setPicture, err = tx.Prepare(`UPDATE person Set Picture = $2 Where id = $1`)
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		return errors.New("failed to remove user picture")
	}
	defer setPicture.Close()

	result, err := setPicture.Exec(uid, nil)
	rowsAff, _ := result.RowsAffected()

	if err != nil || rowsAff != 1 {
		log.Error().Err(err).Msg("error inserting user session")
		return errors.New("failed to remove user picture")
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("error committing changes")
		return errors.New("failed to remove user picture")
	}

	return nil
}

func (us *UserStore) SetUserPicture(ctx context.Context, uid string, pict string) error {
	var tx *sql.Tx
	tx, err := us.db.Begin()

	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return errors.New("failed to add user picture")
	}

	defer tx.Rollback()
	defer func() {
		if rollBackErr := tx.Rollback(); rollBackErr == nil {
			log.Error().Err(err).Msg("rolling back changes")
		}
	}()

	var setPicture *sql.Stmt
	setPicture, err = tx.Prepare(`UPDATE person Set Picture = $2 Where id = $1`)
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		return errors.New("failed to add user picture")
	}
	defer setPicture.Close()

	result, err := setPicture.Exec(uid, pict)
	rowsAff, _ := result.RowsAffected()

	if err != nil || rowsAff != 1 {
		log.Error().Err(err).Msg("error inserting user session")
		return errors.New("failed to add user picture")
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("error committing changes")
		return errors.New("failed to add user picture")
	}

	return nil
}

func (us *UserStore) SetUserSession(ctx context.Context, t store.TokenDetails, uid string, ip string, device string) error {
	var tx *sql.Tx
	tx, err := us.db.Begin()

	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return errors.New("failed to set user session")
	}

	defer tx.Rollback()
	defer func() {
		if rollBackErr := tx.Rollback(); rollBackErr == nil {
			log.Error().Err(err).Msg("rolling back changes")
		}
	}()

	var setSession *sql.Stmt
	setSession, err = tx.Prepare(`INSERT into SESSION (id, userid, ip_address, created_at, device) VALUES ($1,$2,$3,$4,$5)`)
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		return errors.New("failed to set user session")
	}
	defer setSession.Close()

	result, err := setSession.Exec(t.AccessUuid, uid, ip, time.Now(), device)
	rowsAff, _ := result.RowsAffected()

	if err != nil || rowsAff != 1 {
		log.Error().Err(err).Msg("error inserting user session")
		return errors.New("failed to set user session")
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("error committing changes")
		return errors.New("failed to set user session")
	}

	return nil
}

func (us *UserStore) GetUserSession(ctx context.Context, uid string, sessionId string) (store.UserSession, error) {
	psqlToGetSessionInfo := `SELECT ip_address, created_at, updated_at FROM session WHERE id = $1`
	psqlToGetClientInfo := `SELECT id, device FROM session WHERE userid = $1`
	var userSessionResponse store.UserSession
	var userInfo store.UserInfo
	var ip string
	var created_at time.Time
	var updated_at sql.NullTime
	var tmpUat time.Time

	// first sql query to get the user current session info
	sessionInfo := us.db.QueryRowContext(ctx, psqlToGetSessionInfo, sessionId)
	err := sessionInfo.Scan(&ip, &created_at, &updated_at)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return userSessionResponse, errors.New("unable to find the user session")
	}

	userSessionResponse.Ip = ip
	userSessionResponse.Created_at = created_at
	userSessionResponse.Is_current = true

	if !updated_at.Valid {
		userSessionResponse.Updated_at = tmpUat
	} else {
		var tmp time.Time
		err := us.db.QueryRowContext(ctx, `SELECT updated_at from SESSION where id = $1`, sessionId).Scan(&tmp)
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
			return userSessionResponse, errors.New("unable to find the user session")
		}
		userSessionResponse.Updated_at = tmp
	}

	// Second sql query to get the client info
	clientInfo, err := us.db.QueryContext(ctx, psqlToGetClientInfo, uid)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return userSessionResponse, errors.New("unable to find the user session")
	}

	defer clientInfo.Close()

	for clientInfo.Next() {
		var sessionId, device string
		err := clientInfo.Scan(&sessionId, &device)
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
			return userSessionResponse, errors.New("unable to get the user session")
		}

		userInfo.Name = device
		userInfo.SessionId = sessionId
		userSessionResponse.Client = append(userSessionResponse.Client, userInfo)
	}

	return userSessionResponse, nil
}

func (us *UserStore) UpdateUserSession(ctx context.Context, sessionId string) error {
	var tx *sql.Tx
	tx, err := us.db.Begin()

	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return errors.New("failed to update session")
	}

	defer tx.Rollback()
	defer func() {
		if rollBackErr := tx.Rollback(); rollBackErr == nil {
			log.Error().Err(err).Msg("rolling back changes")
		}
	}()

	var updateSession *sql.Stmt
	updateSession, err = tx.Prepare(`UPDATE session Set Updated_at = $2 Where id = $1`)
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		return errors.New("failed to update session")
	}
	defer updateSession.Close()

	result, err := updateSession.Exec(sessionId, time.Now())
	rowsAff, _ := result.RowsAffected()

	if err != nil || rowsAff != 1 {
		log.Error().Err(err).Msg("error updeting user session")
		return errors.New("failed to update session")
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("error committing changes")
		return errors.New("failed to update session")
	}

	return nil
}

func (us *UserStore) DeleteCurrentSession(ctx context.Context, sessionId string) error {
	var tx *sql.Tx
	tx, err := us.db.Begin()

	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return errors.New("failed to delete current session")
	}

	defer tx.Rollback()
	defer func() {
		if rollBackErr := tx.Rollback(); rollBackErr == nil {
			log.Error().Err(err).Msg("rolling back changes")
		}
	}()

	var deleteCurrentSession *sql.Stmt
	deleteCurrentSession, err = tx.Prepare(`DELETE from Session Where id = $1`)
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		return errors.New("failed to delete current session")
	}
	defer deleteCurrentSession.Close()

	result, err := deleteCurrentSession.Exec(sessionId)
	rowsAff, _ := result.RowsAffected()

	if err != nil || rowsAff != 1 {
		log.Error().Err(err).Msg("error deleting user session")
		return errors.New("failed to delete current session")
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("error committing changes")
		return errors.New("failed to delete current session")
	}

	return nil
}

func (us *UserStore) DeleteOtherSession(ctx context.Context, uid string, sid string) error {
	var tx *sql.Tx
	tx, err := us.db.Begin()

	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return errors.New("failed to delete other session")
	}

	defer tx.Rollback()
	defer func() {
		if rollBackErr := tx.Rollback(); rollBackErr == nil {
			log.Error().Err(err).Msg("rolling back changes")
		}
	}()

	var deleteOtherSession *sql.Stmt
	deleteOtherSession, err = tx.Prepare(`DELETE from Session Where userid = $1 AND id <> $2`)
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		return errors.New("failed to delete other session")
	}
	defer deleteOtherSession.Close()

	result, err := deleteOtherSession.Exec(uid, sid)
	rowsAff, _ := result.RowsAffected()

	if err != nil || rowsAff < 1 {
		log.Error().Err(err).Msg("error deleting user session")
		return errors.New("failed to delete other session")
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("error committing changes")
		return errors.New("failed to delete other session")
	}

	return nil
}

func (us *UserStore) GetSessionsId(ctx context.Context, uid string) ([]string, error) {
	var listOfSessionId []string
	psqlStatement := `SELECT id FROM SESSION WHERE userid = $1`

	res, err := us.db.QueryContext(ctx, psqlStatement, uid)

	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return nil, errors.New("unable to find session id")
	}

	defer res.Close()

	for res.Next() {
		var sid string
		err := res.Scan(&sid)
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
			return nil, errors.New("unable to find session id")
		}

		listOfSessionId = append(listOfSessionId, sid)
	}

	return listOfSessionId, nil

}
