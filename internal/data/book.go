package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"codeberg.org/ehrktia/lendbook/internal/data/pg"
	"github.com/jackc/pgx/v5"
)

type BookRepo struct {
	connPool *pg.Postgres
}

func NewBook(connPool *pg.Postgres) *BookRepo {
	return &BookRepo{
		connPool: connPool,
	}
}

type Book struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Edition   string    `json:"edition"`
	Available bool      `json:"available"`
	OwnerID   string    `json:"ownerId"`
	Added     time.Time `json:"added"`
	Updated   time.Time `json:"updated"`
}

var ErrNodata = fmt.Errorf("%s", "no data found")

func (b BookRepo) GetBooks(
	ctx context.Context, of, limit int) ([]Book, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	conn, err := b.connPool.GetConn(reqCtx)
	if err != nil {
		return nil, err
	}
	defer func() {
		conn.Conn().Close(reqCtx)
		conn.Release()
	}()
	var query = `select * from book offset @offset limit @limit`
	args := pgx.NamedArgs{"offset": of, "limit": limit}
	rows, err := conn.Query(reqCtx, query, args)
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
	return books, nil
}

func (b BookRepo) GetBookCount(ctx context.Context) (int, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	conn, err := b.connPool.GetConn(reqCtx)
	if err != nil {
		return 0, err
	}
	defer func() {
		conn.Conn().Close(reqCtx)
		conn.Release()
	}()
	q := `select count(id) from book;`
	rows, err := conn.Query(reqCtx, q)
	if err != nil {
		return 0, errClassify(err)
	}
	tot, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, errClassify(err)
	}
	return tot, nil
}

var queryBookOwnerId = `select * from book where owner_id=$1`

func (b BookRepo) GetBooksByOwnerId(
	ctx context.Context, ownerId float64) ([]Book, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	conn, err := b.connPool.GetConn(reqCtx)
	if err != nil {
		return nil, err
	}
	defer func() {
		conn.Conn().Close(reqCtx)
		conn.Release()
	}()
	rows, err := conn.Query(reqCtx, queryBookOwnerId, ownerId)
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
	return books, nil
}
