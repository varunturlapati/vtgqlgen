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
	GetFruit(ctx context.Context, id int) (*entity.Fruit, error)
	ListFruits(ctx context.Context) ([]*entity.Fruit, error)
	CreateFruit(ctx context.Context, arg *entity.CreateFruitParams) (*entity.Fruit, error)
	UpdateFruit(ctx context.Context, arg *entity.UpdateFruitParams) (*entity.Fruit, error)
	DeleteFruit(ctx context.Context, id int) (*entity.Fruit, error)

	// detail queries
	GetDetail(ctx context.Context, name string) (*entity.Detail, error)
	ListDetails(ctx context.Context) ([]*entity.Detail, error)

	// level queries
	GetLevel(ctx context.Context, level string) (*entity.Level, error)
	ListLevels(ctx context.Context) ([]*entity.Level, error)

	// rack queries
	GetRack(ctx context.Context, id int) (*entity.Rack, error)
	ListRacks(ctx context.Context) ([]*entity.Rack, error)

	ListRacksByFruitIDs(ctx context.Context, fruitIDs []int) ([]r.ListRacksByFruitIDsRow, error)
	ListFruitsByRackIDs(ctx context.Context, rackIDs []int) ([]d.ListFruitsByRackIDsRow, error)
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
