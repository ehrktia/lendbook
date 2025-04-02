package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"codeberg.org/ehrktia/lendbook/internal/data/pg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepo struct {
	connPool *pg.Postgres
}

func NewUser(connPool *pg.Postgres) UserRepo {
	return UserRepo{
		connPool: connPool,
	}
}

var ErrDuplicate = errors.New("duplicate data found in store")

func (o UserRepo) GetById(
	ctx context.Context, id string) (User, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	conn, err := o.connPool.GetConn(reqCtx)
	if err != nil {
		return User{}, err
	}
	defer func() {
		conn.Conn().Close(reqCtx)
		conn.Release()
	}()
	q := `select
	*
from
	owner
where
	id = @id;`
	args := pgx.NamedArgs{
		"id": id,
	}
	rows, err := conn.Query(reqCtx, q, args)
	if err != nil {
		return User{}, errClassify(err)
	}
	owner, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		return User{}, errClassify(err)
	}
	if ctx.Err() != nil {
		return owner, ctx.Err()
	}
	return owner, nil
}

func (o UserRepo) GetBookByUserId(
	ctx context.Context, ownerId float64) ([]Book, error) {
	conn, err := o.connPool.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		conn.Conn().Close(ctx)
		conn.Release()
	}()
	queryBookByUserId := `select * from book where owner_id=$1`
	rows, err := conn.Query(ctx, queryBookByUserId, ownerId)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNodata
		default:
			return nil, err
		}
	}
	books, err := pgx.CollectRows(rows, pgx.RowToStructByName[Book])
	if err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return books, ctx.Err()
	default:
		return books, nil
	}
}

type BookList struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Edition   string `json:"edition"`
	Available bool   `json:"available"`
	OwnerID   int64  `json:"ownerId"`
	Added     string `json:"added"`
	Updated   string `json:"updated"`
}

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Added     time.Time `json:"added"`
	Updated   time.Time `json:"updated"`
}

func (o UserRepo) GetUsers(ctx context.Context) ([]User, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	conn, err := o.connPool.GetConn(reqCtx)
	if err != nil {
		return nil, err
	}
	defer func() {
		conn.Conn().Close(reqCtx)
		conn.Release()
	}()
	qGetUsersBooks := `select
	*
from
	owner owb;`
	rows, err := conn.Query(reqCtx, qGetUsersBooks)
	if err != nil {
		return nil, errClassify(err)
	}
	owners, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		return nil, errClassify(err)
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return owners, nil
}

var ErrCreateUser = errors.New("can not create owner")

func (o UserRepo) Create(
	ctx context.Context, owner User) (string, error) {
	var ownerId string
	reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	conn, err := o.connPool.GetConn(reqCtx)
	if err != nil {
		return ownerId, errClassify(err)
	}
	q := `insert
	into
	public."owner" (
	first_name,
	last_name,
	email)
values( @fname, @lname, @email) returning id;`
	args := pgx.NamedArgs{
		"fname": owner.FirstName,
		"lname": owner.LastName,
		"email": owner.Email,
	}
	f := func(tx pgx.Tx) error {
		rows, err := tx.Query(reqCtx, q, args)
		if err != nil {
			return err
		}
		ownerId, err = pgx.CollectExactlyOneRow(rows, pgx.RowTo[string])
		if err != nil {
			return err
		}
		return nil
	}
	if err := pg.RuninTx(reqCtx, conn, f); err != nil {
		return ownerId, errClassify(err)
	}
	return ownerId, nil
}

type UserWithNoBooks struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Active    bool   `json:"active"`
	Version   int    `json:"version"`
}

func (owb UserWithNoBooks) String() string {
	return fmt.Sprintf("%#v", owb)

}

func (o UserRepo) Update(ctx context.Context, owner UserWithNoBooks) (
	UserWithNoBooks, error) {
	reqCtx, cancel := context.WithTimeoutCause(ctx,
		3*time.Second, errors.New("request timed out "))
	defer cancel()
	ow := UserWithNoBooks{}
	conn, err := o.connPool.GetConn(reqCtx)
	if err != nil {
		return ow, errClassify(err)
	}

	qUpdateUser := `update
	public."lender"
set
	first_name = @first_name,
	last_name = @last_name,
	email = @email,
	active = @active,
	version = @version
where
	id = @id and 
	version < @version
	returning *;`

	args := pgx.NamedArgs{
		"first_name": owner.FirstName,
		"last_name":  owner.LastName,
		"email":      owner.Email,
		"active":     owner.Active,
		"id":         owner.ID,
		"version":    owner.Version,
	}
	f := func(tx pgx.Tx) error {
		rows, err := tx.Query(reqCtx, qUpdateUser, args)
		if err != nil {
			return err
		}
		ow, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[UserWithNoBooks])
		return errClassify(err)
	}
	if err := pg.RuninTx(reqCtx, conn, f); err != nil {
		return ow, errClassify(err)
	}
	if ctx.Err() != nil {
		return ow, ctx.Err()
	}
	return ow, nil
}

var PgErr = func(e string) error {
	er := errors.New("error from store")
	return fmt.Errorf("%w:%s", er, e)
}

func errClassify(err error) error {
	var pgConn *pgconn.ConnectError
	var pgErr *pgconn.PgError
	if err != nil {
		switch {
		case errors.As(err, &pgConn):
			return PgErr(pgConn.Error())
		case errors.As(err, &pgErr):
			return PgErr(pgErr.Error())
		case errors.Is(err, pgx.ErrNoRows):
			return ErrNodata
		case errors.Is(err, pgx.ErrTooManyRows):
			return ErrDuplicate
		case errors.Is(err, pgx.ErrTxClosed):
			return ErrClosedTx
		case errors.Is(err, pgx.ErrTxCommitRollback):
			return ErrTxRollBack
		default:
			return err
		}
	}
	return nil
}

var ErrClosedTx = errors.New("tx is closed")
var ErrTxRollBack = errors.New("tx commit rollback error")

func (o UserRepo) GetUserByEmail(
	ctx context.Context, email string) (string, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	conn, err := o.connPool.GetConn(reqCtx)
	if err != nil {
		return "", err
	}
	defer func() {
		conn.Conn().Close(reqCtx)
		conn.Release()
	}()
	args := pgx.NamedArgs{
		"email": email,
	}
	qGetUserByName := `select id from lender
	where email=@email;`
	row, err := conn.Query(reqCtx, qGetUserByName, args)
	if err != nil {
		return "", errClassify(err)
	}
	id, err := pgx.CollectExactlyOneRow(row, pgx.RowTo[string])
	if err != nil {
		return "", errClassify(err)
	}
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	return id, nil
}

func (o UserRepo) DeleteByID(ctx context.Context, id string) error {
	reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	conn, err := o.connPool.GetConn(reqCtx)
	if err != nil {
		return err
	}
	defer func() {
		conn.Conn().Close(reqCtx)
		conn.Release()
	}()
	args := pgx.NamedArgs{
		"id": id,
	}
	q := `delete from owner where id=@id;`
	if _, err := conn.Query(reqCtx, q, args); err != nil {
		return errClassify(err)
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return nil
}
