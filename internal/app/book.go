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

func (b Book) GetAll(ctx context.Context, o, l string) (model.BookList, error) {
	of, err := toInt(o)
	if err != nil {
		return model.BookList{}, err
	}
	limit, err := toInt(l)
	if err != nil {
		return model.BookList{}, err
	}
	books, err := b.query.GetBooks(ctx, of, limit)
	if err != nil {
		return model.BookList{}, err
	}
	return populateBookResults(books, of, limit), nil

}

func populateBookResults(books []data.Book, p, limit int) model.BookList {
	b := make([]*model.Book, len(books))
	for i, v := range books {
		bookModel := toBookModel(v)
		b[i] = &bookModel
	}
	r := model.BookList{
		Data: b,
		Prev: strconv.Itoa(p),
		Next: strconv.Itoa(limit + 1),
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
