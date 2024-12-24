package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ehrktia/lendbook/internal/data/pg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type OwnerRepo struct {
	connPool *pg.Postgres
}

func NewOwner(connPool *pg.Postgres) *OwnerRepo {
	return &OwnerRepo{
		connPool: connPool,
	}
}

var ErrDuplicate = errors.New("duplicate data found in store")

func (o OwnerRepo) GetById(
	ctx context.Context, id float64) (Owner, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	conn, err := o.connPool.GetConn(reqCtx)
	if err != nil {
		return Owner{}, err
	}
	defer func() {
		conn.Conn().Close(reqCtx)
		conn.Release()
	}()
	qOwnerById := `select
	*
from
	owner_with_books owb
where
	id = $1;`
	rows, err := conn.Query(reqCtx, qOwnerById, id)
	if err != nil {
		return Owner{}, errClassify(err)
	}
	owner, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Owner])
	if err != nil {
		return Owner{}, errClassify(err)
	}
	if ctx.Err() != nil {
		return owner, ctx.Err()
	}

	return owner, nil
}

func (o OwnerRepo) GetBookByOwnerId(
	ctx context.Context, ownerId float64) ([]Book, error) {
	conn, err := o.connPool.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		conn.Conn().Close(ctx)
		conn.Release()
	}()
	queryBookByOwnerId := `select * from books where owner_id=$1`
	rows, err := conn.Query(ctx, queryBookByOwnerId, ownerId)
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
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Edition   string `json:"edition"`
	Available bool   `json:"available"`
	OwnerID   int64  `json:"ownerId"`
	Added     string `json:"added"`
	Updated   string `json:"updated"`
}

type Owner struct {
	ID        int64      `json:"id"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Email     string     `json:"email"`
	Active    bool       `json:"active"`
	Version   int32      `json:"version"`
	Books     []BookList `json:"books"`
}

func (o OwnerRepo) GetOwners(ctx context.Context) ([]Owner, error) {
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
	qGetOwnersBooks := `select
	*
from
	owner_with_books owb;`
	rows, err := conn.Query(reqCtx, qGetOwnersBooks)
	if err != nil {
		return nil, errClassify(err)
	}
	owners, err := pgx.CollectRows(rows, pgx.RowToStructByName[Owner])
	if err != nil {
		return nil, errClassify(err)
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return owners, nil
}

var ErrCreateOwner = errors.New("can not create owner")

func (o OwnerRepo) Create(
	ctx context.Context, owner Owner) (int64, error) {
	var ownerId int64
	reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	conn, err := o.connPool.GetConn(reqCtx)
	if err != nil {
		return ownerId, errClassify(err)
	}
	qCreateUser := `INSERT INTO public."owner"
(first_name, last_name, email, active,version)
VALUES(@first_name, @last_name, @email, @active,@version) returning id;`
	args := pgx.NamedArgs{
		"first_name": owner.FirstName,
		"last_name":  owner.LastName,
		"email":      owner.Email,
		"active":     owner.Active,
		"version":    owner.Version,
	}
	f := func(tx pgx.Tx) error {
		rows, err := tx.Query(reqCtx, qCreateUser, args)
		if err != nil {
			return err
		}
		ownerId, err = pgx.CollectExactlyOneRow(rows, pgx.RowTo[int64])
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

type OwnerWithNoBooks struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Active    bool   `json:"active"`
	Version   int    `json:"version"`
}

func (owb OwnerWithNoBooks) String() string {
	return fmt.Sprintf("%#v", owb)

}

func (o OwnerRepo) Update(
	ctx context.Context, owner OwnerWithNoBooks) (OwnerWithNoBooks, error) {
	reqCtx, cancel := context.WithTimeoutCause(ctx,
		3*time.Second, errors.New("request timed out "))
	defer cancel()
	ow := OwnerWithNoBooks{}
	conn, err := o.connPool.GetConn(reqCtx)
	if err != nil {
		return ow, errClassify(err)
	}
	qUpdateUser := `update
	public."owner"
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
		ow, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[OwnerWithNoBooks])
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

func (o OwnerRepo) GetOwnerByEmail(
	ctx context.Context, email string) (int64, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	conn, err := o.connPool.GetConn(reqCtx)
	if err != nil {
		return int64(0), err
	}
	defer func() {
		conn.Conn().Close(reqCtx)
		conn.Release()
	}()
	args := pgx.NamedArgs{
		"email": email,
	}
	qGetOwnerByName := `select id from owner
	where email=@email;`
	row, err := conn.Query(reqCtx, qGetOwnerByName, args)
	if err != nil {
		return 0, errClassify(err)
	}
	id, err := pgx.CollectExactlyOneRow(row, pgx.RowTo[int64])
	if err != nil {
		return 0, errClassify(err)
	}
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	return id, nil

}
