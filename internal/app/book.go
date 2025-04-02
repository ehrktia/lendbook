package app

import (
	"context"
	"log/slog"
	"strconv"
	"time"

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

func (b Book) GetAll(
	ctx context.Context, o, l string) (*model.BookList, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	of, err := toInt(o)
	if err != nil {
		return nil, err
	}
	limit, err := toInt(l)
	if err != nil {
		return nil, err
	}
	tot, err := b.query.GetBookCount(reqCtx)
	if err != nil {
		return nil, err
	}
	books, err := b.query.GetBooks(ctx, of, limit)
	if err != nil {
		return nil, err
	}
	return populateBookResults(books, of, limit, tot), nil
}

func populateBookResults(
	books []data.Book, p, limit, tot int) *model.BookList {
	b := make([]*model.Book, len(books))
	for i, v := range books {
		bookModel := toBookModel(v)
		b[i] = &bookModel
	}
	next := p + limit
	r := &model.BookList{
		Data:  b,
		Prev:  strconv.Itoa(p),
		Next:  strconv.Itoa(min(next, tot)),
		Total: strconv.Itoa(tot),
	}
	return r
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
		OwnerID:   in.OwnerID,
		Added:     &add,
		Updated:   &upd,
	}
}
