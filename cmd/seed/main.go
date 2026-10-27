package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"time"

	"github.com/clary-work01/hot-post-api/internal/config"
	"github.com/clary-work01/hot-post-api/internal/domain/post"
	"github.com/clary-work01/hot-post-api/internal/domain/user"
	"github.com/clary-work01/hot-post-api/internal/infra/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	if err := run(); err != nil {
		slog.Error("seed exited with error", "error", err)
		os.Exit(1)
	}
	slog.Info("seed completed successfully")
}

func run() error {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// ---- Postgres ----
	gormDB, err := db.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect postgres failed: %w", err)
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB failed: %w", err)
	}
	defer sqlDB.Close()
	slog.Info("connect postgres success")

	// ---- DB：全部包在一個交易，失敗就整批回滾 ----
	err = gormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		users, err := seedUsers(tx)
		if err != nil {
			return fmt.Errorf("users: %w", err)
		}
		slog.Info("seed users success", "count", len(users))

		posts, err := seedPosts(tx, users)
		if err != nil {
			return fmt.Errorf("posts: %w", err)
		}
		slog.Info("seed posts success", "count", len(posts))

		return nil
	})

	if err != nil {
		return fmt.Errorf("seed db failed: %w", err)
	}

	return nil
}

func seedUsers(tx *gorm.DB) ([]user.User, error) {
	id, err := uuid.Parse("d2c5d597-4b6b-4d11-ac6a-06b5b66e1f47") // TODO
	if err != nil {
		return nil, err
	}

	users := []user.User{
		{ID: id, Username: "user01", Password: "123"},
		{Username: "user02", Password: "123"},
		{Username: "user03", Password: "123"},
		{Username: "user04", Password: "123"},
	}

	if err := tx.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&users).Error; err != nil {
		return nil, err
	}

	if err := tx.
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func seedPosts(tx *gorm.DB, users []user.User) ([]post.Post, error) {
	posts := []post.Post{
		{AuthorID: users[0].ID, Title: "測試文章01", Content: "測試文章01內容"},
		{AuthorID: users[1].ID, Title: "測試文章02", Content: "測試文章02內容"},
		{AuthorID: users[2].ID, Title: "測試文章03", Content: "測試文章03內容"},
		{AuthorID: users[3].ID, Title: "測試文章04", Content: "測試文章04內容"},
	}

	if err := tx.Create(&posts).Error; err != nil {
		return nil, err
	}

	if err := tx.Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}
