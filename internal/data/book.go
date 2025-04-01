package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ehrktia/lendbook/internal/data/pg"
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
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Edition   string    `json:"edition"`
	Available bool      `json:"available"`
	OwnerID   int64     `json:"ownerId"`
	Added     time.Time `json:"added"`
	Updated   time.Time `json:"updated"`
}

var ErrNodata = fmt.Errorf("%s", "no data found")

func (b BookRepo) GetBooks(ctx context.Context, of, limit int) ([]Book, error) {
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
	var query = `select * from books offset @offset limit @limit`
	args := pgx.NamedArgs{"offset": 0, "limit": 10}
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

var queryBookOwnerId = `select * from books where owner_id=$1`

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
