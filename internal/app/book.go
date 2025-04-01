package app

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/ehrktia/lendbook/internal/graph/model"
)

type Book struct {
	query BookQuery
	l     *slog.Logger
}

func NewBook(query BookQuery, l *slog.Logger) Book {
	l = l.WithGroup("app-layer-book")
	return Book{
		query: query,
		l:     l,
	}
}

func (b Book) GetAll(ctx context.Context, o, l string) ([]model.Book, error) {
	of, err := strconv.Atoi(o)
	if err != nil {
		return nil, err
	}
	limit, err := strconv.Atoi(l)
	if err != nil {
		return nil, err
	}
	_, _ = b.query.GetBooks(ctx, of, limit)
	return []model.Book{}, nil
}
