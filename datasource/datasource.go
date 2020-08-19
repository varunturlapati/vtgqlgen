package datasource

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"log"

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
	ListRacksByIDs(ctx context.Context, IDs []int) ([]*entity.Rack, error)
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
	rr := r.New(NetboxUrl)
	return &repoSvc{
		db:           sqlDB,
		Queries:      q,
		RestRequests: rr,
	}
}

func (rs *repoSvc) ListRacksByIDs(ctx context.Context, IDs []int) ([]*entity.Rack, error) {
	log.Println("I reached ListRacksByIDs")
	var rackList []*entity.Rack
	for _, i := range IDs {
		var rk entity.Rack
		restRes, err := rs.GetRack(ctx, i)
		if err != nil {
			log.Println("Rest Query for Rack threw and error")
		}
		rk.Created = restRes.Created
		rk.CustomFields = restRes.CustomFields
		res, err := rs.Queries.GetRackFromDB(ctx, i)
		if err != nil {
			log.Println("DB Query for Rack threw an error")
		}
		if restRes.Id != 0 {
			rk.Id = restRes.Id
		}
		if res.Name != "" {
			rk.Name = res.Name
		}

		rk.Ipaddr = res.Ipaddr
		rk.Live = res.Live
		rackList = append(rackList, &rk)
	}
	return rackList, nil
}

func (rs *repoSvc) GetRack(ctx context.Context, id int) (*entity.Rack, error) {
	log.Println("I reached GetRack")
	var rk entity.Rack
	restRes, err := rs.GetRackFromNetbox(ctx, id)
	if err != nil {
		log.Println("Rest Query for Rack threw and error")
	}
	rk.Created = restRes.Created
	rk.CustomFields = restRes.CustomFields
	res, err := rs.Queries.GetRackFromDB(ctx, id)
	if err != nil {
		log.Println("DB Query for Rack threw an error")
	}
	if restRes.Id != 0 {
		rk.Id = restRes.Id
	}
	if res.Name != "" {
		rk.Name = res.Name
	}
	rk.Ipaddr = res.Ipaddr
	rk.Live = res.Live
	return &rk, nil
}

func (rs *repoSvc) ListRacks(ctx context.Context) ([]*entity.Rack, error) {
	log.Println("I reached ListRacks")
	rackList := make([]*entity.Rack, 0)
	restRes, err := rs.RestRequests.ListRacksFromNetbox(ctx)
	if err != nil {
		log.Println("Couldn't get Netbox results for Racks")
	}
	rackList = append(rackList, restRes...)
	log.Printf("Num of rack entries after Netbox query: %d and Result length = %d\n", len(rackList), len(restRes))
	Id2RackMap := make(map[int]*entity.Rack)
	for _, i := range restRes {
		Id2RackMap[int(i.Id)] = i
	}
	res, err := rs.Queries.ListRacksFromDB(ctx)
	if err != nil {
		log.Println("Couldn't get DB results for Racks")
		return rackList, err
	}
	for _, rack := range res {
		var rk entity.Rack
		tmp, ok := Id2RackMap[int(rack.Id)]
		if !ok {
			log.Println("Found an entry not in Netbox")
			rackList = append(rackList, rack)
		} else {
			log.Println("This entry is in Netbox")
			rk.Id = rack.Id
			rk.Name = rack.Name
			rk.Ipaddr = rack.Ipaddr
			rk.Live = rack.Live
			rk.CustomFields = tmp.CustomFields
			rk.Created = tmp.Created
			rackList = append(rackList, &rk)
		}
	}
	log.Printf("Returning rackList with final length = %d\n", len(rackList))
	return rackList, nil
}
