package app

import (
	"context"
	"log/slog"

	"github.com/ehrktia/lendbook/internal/data"
	"github.com/ehrktia/lendbook/internal/graph/model"
)

type Commander interface {
	GetID
	Create(ctx context.Context, owner data.User) (string, error)
	Update(ctx context.Context, owner data.UserWithNoBooks) (data.UserWithNoBooks, error)
}

type UserQuery interface {
	GetID
	GetBookByUserId(ctx context.Context, ownerId float64) ([]data.Book, error)
	GetUsers(ctx context.Context) ([]data.User, error)
	GetUserByEmail(ctx context.Context, email string) (string, error)
}

type BookQuery interface {
	GetBooks(ctx context.Context, of, limit int) ([]data.Book, error)
}

type GetID interface {
	GetById(ctx context.Context, id string) (data.User, error)
}

type User struct {
	command Commander
	query   UserQuery
	l       *slog.Logger
}

func NewUser(command Commander, query UserQuery, l *slog.Logger) User {
	l = l.WithGroup("app-layer-user")
	return User{
		command: command,
		query:   query,
		l:       l,
	}

}

func (o User) Create(ctx context.Context, user model.User) (
	model.User, error) {
	do := data.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Active:    user.Active,
		Books:     make([]data.BookList, 0),
	}
	userId, err := o.command.Create(ctx, do)
	if err != nil {
		o.l.LogAttrs(ctx, slog.LevelError, "error creating owner",
			slog.String("error", err.Error()))
		return user, err
	}
	user.ID = userId
	return user, nil
}

func (o User) GetById(ctx context.Context, id string) (
	model.User, error) {
	// operate
	ow, err := o.query.GetById(ctx, id)
	if err != nil {
		o.l.LogAttrs(ctx, slog.LevelError, "error getting user by id",
			slog.String("user-id", id),
			slog.String("error", err.Error()))
		return model.User{}, err
	}
	// map to output format
	result := model.User{
		ID:        ow.ID,
		FirstName: ow.FirstName,
		LastName:  ow.LastName,
		Email:     ow.Email,
		Active:    ow.Active,
		Version:   ow.Version,
	}
	if len(ow.Books) < 1 {
		return result, nil
	}
	bl := make([]*model.Book, len(ow.Books))
	for i, v := range ow.Books {
		bl[i] = &model.Book{
			ID:        v.ID,
			Title:     v.Title,
			Author:    v.Author,
			Edition:   v.Edition,
			Available: v.Available,
			OwnerID:   float64(v.OwnerID),
		}
		if v.Added != "" {
			bl[i].Added = &v.Added
		}
		if v.Updated != "" {
			bl[i].Updated = &v.Updated
		}

	}
	// result
	return result, nil
}

func (o User) GetUserByEmail(ctx context.Context, email string) (
	model.User, error) {
	uid, err := o.query.GetUserByEmail(ctx, email)
	if err != nil {
		o.l.LogAttrs(ctx, slog.LevelError, "error fetching userid for email",
			slog.String("error", err.Error()),
			slog.String("user-email", email))
		return model.User{}, err
	}
	u, err := o.query.GetById(ctx, uid)
	if err != nil {
		o.l.LogAttrs(ctx, slog.LevelError, "error fetching user by id",
			slog.String("error", err.Error()),
			slog.String("user-email", email),
			slog.String("user-id", uid))
		return model.User{}, err
	}
	// map to output format
	result := model.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Active:    u.Active,
		Version:   u.Version,
	}
	if len(u.Books) < 1 {
		return result, nil
	}
	bl := make([]*model.Book, len(u.Books))
	for i, v := range u.Books {
		bl[i] = &model.Book{
			ID:        v.ID,
			Title:     v.Title,
			Author:    v.Author,
			Edition:   v.Edition,
			Available: v.Available,
			OwnerID:   float64(v.OwnerID),
		}
		if v.Added != "" {
			bl[i].Added = &v.Added
		}
		if v.Updated != "" {
			bl[i].Updated = &v.Updated
		}

	}
	// result
	return result, err

}
