package app

import (
	"context"
	"log/slog"

	"github.com/ehrktia/lendbook/internal/data"
	"github.com/ehrktia/lendbook/internal/graph/model"
)

type User struct {
	repo data.UserRepo
	l    *slog.Logger
}

func NewUser(repo data.UserRepo, l *slog.Logger) User {
	l = l.WithGroup("app-layer-owner")
	return User{
		repo: repo,
		l:    l,
	}

}

func (o User) Create(ctx context.Context, user model.User) (model.User, error) {

	do := data.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Active:    user.Active,
		Books:     make([]data.BookList, 0),
	}
	userId, err := o.repo.Create(ctx, do)
	if err != nil {
		o.l.LogAttrs(ctx, slog.LevelError, "error creating owner",
			slog.String("error", err.Error()))
		return user, err
	}
	user.ID = float64(userId)
	return user, nil
}

func (o User) GetById(ctx context.Context, id float64) (model.User, error) {
	// operate
	ow, err := o.repo.GetById(ctx, id)
	if err != nil {
		o.l.LogAttrs(ctx, slog.LevelError, "error getting owner by id",
			slog.Float64("owner-id", id),
			slog.String("error", err.Error()))
		return model.User{}, err
	}
	// map to output format
	result := model.User{
		ID:        float64(ow.ID),
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
			ID:        float64(v.ID),
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

func (o User) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	uid, err := o.repo.GetUserByEmail(ctx, email)
	if err != nil {
		o.l.LogAttrs(ctx, slog.LevelError, "error fetching userid for email",
			slog.String("error", err.Error()),
			slog.String("user-email", email))
		return model.User{}, err
	}
	u, err := o.repo.GetById(ctx, float64(uid))
	if err != nil {
		o.l.LogAttrs(ctx, slog.LevelError, "error fetching user by id",
			slog.String("error", err.Error()),
			slog.String("user-email", email),
			slog.Float64("user-id", float64(uid)))
		return model.User{}, err
	}
	// map to output format
	result := model.User{
		ID:        float64(u.ID),
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
			ID:        float64(v.ID),
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

