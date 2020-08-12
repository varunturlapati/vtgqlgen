package db

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Repository interface {
	// fruit queries
	GetFruit(ctx context.Context, id int) (Fruit, error)
	ListFruits(ctx context.Context) ([]Fruit, error)

	// detail queries
	GetDetail(ctx context.Context, name string) (Detail, error)
	ListDetails(ctx context.Context) ([]Detail, error)

	// level queries
	GetLevel(ctx context.Context, level string) (Level, error)
	ListLevels(ctx context.Context) ([]Level, error)

	// rack queries
	GetRack(ctx context.Context, id int) (ServerRack, error)
	ListRacks(ctx context.Context) ([]ServerRack, error)
}

type repoSvc struct {
	*Queries
	db *sql.DB
}

func (r *repoSvc) withTx(ctx context.Context, txFn func(*Queries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = txFn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			err = fmt.Errorf("tx failed: %v unable to rollback: %v", err, rbErr)
		}
	} else {
		err = tx.Commit()
	}
	return err
}

func NewRepository(db *gorm.DB) Repository {
	abc, err := db.DB()
	if err != nil {
		return nil
	}
	xyz := New(abc)
	return &repoSvc{
		db:      abc,
		Queries: xyz,
	}
}

func Open(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

