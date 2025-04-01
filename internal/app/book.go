package app

import (
	"context"
	"log/slog"
	"strconv"

	"codeberg.org/ehrktia/lendbook/internal/data"
	"codeberg.org/ehrktia/lendbook/internal/graph/model"
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

func (b Book) GetAll(ctx context.Context, o, l string) ([]model.BookList, error) {
	of, err := toInt(o)
	if err != nil {
		return nil, err
	}
	limit, err := toInt(l)
	if err != nil {
		return nil, err
	}
	_, _ = b.query.GetBooks(ctx, of, limit)
	return []model.BookList{}, nil
}

func toInt(in string) (int, error) {
	return strconv.Atoi(in)
}

func toBookModel(in data.Book) model.Book {
	add := in.Added.String()
	upd := in.Updated.String()
	return model.Book{
		ID:        in.ID,
		Title:     in.Title,
		Author:    in.Author,
		Edition:   in.Edition,
		Available: in.Available,
		OwnerID:   float64(in.OwnerID),
		Added:     &add,
		Updated:   &upd,
	}
}
