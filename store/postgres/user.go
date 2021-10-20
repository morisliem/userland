package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"userland/store"

	"github.com/gofrs/uuid"
)

type User struct {
	Id         uint64    `json:"Id"`
	Fullname   string    `json:"fullname"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Location   string    `json:"location"`
	Bio        string    `json:"bio"`
	Web        string    `json:"web"`
	Picture    string    `json:"picture"`
	Created_at time.Time `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) store.UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) GetUser(ctx context.Context) error {

	return nil
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
	// var tx *sql.Tx
	// tx, err := us.db.Begin()
	// if err != nil {
	// 	fmt.Println("failed to begin transaction")
	// 	return err
	// }
	// defer tx.Rollback()

	// if us.EmailExist(ctx, u) != nil {
	// 	fmt.Println("user exists")
	// 	if rollbackErr := tx.Rollback(); rollbackErr != nil {
	// 		fmt.Println("Error occur rolling back changer")
	// 	}
	// 	return err
	// }

	// var inserUserStmt *sql.Stmt
	// inserUserStmt, err = tx.Prepare("INSERT INTO person (fullname, email, password, created_at, id) values ($1, $2, $3, $4, $5);")
	// if err != nil {
	// 	fmt.Println("error preparing statement")
	// 	if rollbackErr := tx.Rollback(); rollbackErr != nil {
	// 		fmt.Println("Error occur rolling back changer")
	// 	}
	// 	return err
	// }
	// defer inserUserStmt.Close()

	// userId, err := uuid.NewV4()
	// if err != nil {
	// 	return errors.New("failed to generate user id")
	// }

	// var result sql.Result
	// result, err = inserUserStmt.Exec(u.Fullname, u.Email, u.Password, time.Now().UTC(), userId.String())

	// rowsAff, _ := result.RowsAffected()

	// if err != nil || rowsAff != 1 {
	// 	fmt.Println("error inserting new user")
	// 	if rollbackErr := tx.Rollback(); rollbackErr != nil {
	// 		fmt.Println("Error occur rolling back changer")
	// 	}
	// 	return err
	// }

	// rand.Seed(time.Now().UnixNano())
	// rn := rand.Intn(100000)
	// fmt.Println("random number:", rn)
	// var insertEmailVerStmt *sql.Stmt
	// insertEmailVerStmt, err = tx.Prepare("INSERT INTO email_ver (fullname, email, ver_code) VALUES ($1, $2, $3);")
	// if err != nil {
	// 	fmt.Println("failed to prepare statement")
	// 	if rollbackErr := tx.Rollback(); rollbackErr != nil {
	// 		fmt.Println("Error occur rolling back changer")
	// 	}
	// 	return err
	// }
	// defer insertEmailVerStmt.Close()
	// result, err = insertEmailVerStmt.Exec(u.Fullname, u.Email, rn)
	// rowsAff, _ = result.RowsAffected()

	// if err != nil || rowsAff != 1 {
	// 	fmt.Println("error inserting into email_ver")
	// 	if rollbackErr := tx.Rollback(); rollbackErr != nil {
	// 		fmt.Println("Error occur rolling back changer")
	// 	}
	// 	return err
	// }

	// err = helper.SendEmail(u.Email, rn)
	// if err != nil {
	// 	fmt.Println("error emailing verification code")
	// 	if rollbackErr := tx.Rollback(); rollbackErr != nil {
	// 		fmt.Println("Error occur rolling back changer")
	// 	}
	// 	return err
	// }

	// err = tx.Commit()
	// if err != nil {
	// 	fmt.Println("error committing changes")
	// 	if rollbackErr := tx.Rollback(); rollbackErr != nil {
	// 		fmt.Println("Error occur rolling back changer")
	// 	}
	// 	return err
	// }

	// fmt.Println("success")

	psqlStatement := `INSERT INTO PERSON (Fullname,Email,Password,Created_at,Id) values ($1,$2,$3,$4,$5)`

	userId, err := uuid.NewV4()
	if err != nil {
		return errors.New("failed to generate user id")
	}

	_, err = us.db.Exec(psqlStatement, u.Fullname, u.Email, u.Password, time.Now().UTC(), userId.String())

	if err != nil {
		return errors.New("failed to store user")
	}
	return nil
}
