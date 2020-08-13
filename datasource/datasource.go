package datasource

import (
	"context"
	"database/sql"
	"gorm.io/gorm"

	d "github.com/varunturlapati/vtgqlgen/datasource/db"
	r "github.com/varunturlapati/vtgqlgen/datasource/rest"
	"github.com/varunturlapati/vtgqlgen/pkg/entity"
)

type Repository interface {
	// fruit queries
	GetFruit(ctx context.Context, id int) (entity.Fruit, error)
	ListFruits(ctx context.Context) ([]entity.Fruit, error)

	// detail queries
	GetDetail(ctx context.Context, name string) (entity.Detail, error)
	ListDetails(ctx context.Context) ([]entity.Detail, error)

	// level queries
	GetLevel(ctx context.Context, level string) (entity.Level, error)
	ListLevels(ctx context.Context) ([]entity.Level, error)

	// rack queries
	GetRack(ctx context.Context, id int) (entity.ServerRack, error)
	ListRacks(ctx context.Context) ([]entity.ServerRack, error)
}

const (
	NetboxUrl = "http://localhost:8000/"
)

type repoSvc struct {
	db *sql.DB
	*d.Queries
	*r.RestRequests
}

func NewRepository(db *gorm.DB) Repository {
	sqlDB, err := db.DB()
	if err != nil {
		return nil
	}
	q := d.New(sqlDB)
	r := r.New(NetboxUrl)
	return &repoSvc{
		db:           sqlDB,
		Queries:      q,
		RestRequests: r,
	}
}
