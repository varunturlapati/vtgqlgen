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

	// server queries
	GetServerByName(ctx context.Context, name string) (*entity.Server, error)
	GetServerById(ctx context.Context, id int) (*entity.Server, error)
	ListServers(ctx context.Context) ([]*entity.Server, error)

	ListRacksByFruitIDs(ctx context.Context, fruitIDs []int) ([]r.ListRacksByFruitIDsRow, error)
	ListFruitsByRackIDs(ctx context.Context, rackIDs []int) ([]d.ListFruitsByRackIDsRow, error)
	ListRacksByIDs(ctx context.Context, IDs []int) ([]*entity.Rack, error)

	// TODO - is there an impact on helpers like this one on which we may not want to directly expose queries? The answer may be similar to
	// join cases like ListRacksByFruitIDs() above
	GetServerStatusById(ctx context.Context, id int) (*string, error)
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
	restRes, err := rs.GetRackFromNetbox(ctx, id) // GetRack - Rest netbox.go
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
	restRes, err := rs.ListRacksFromNetbox(ctx)
	if err != nil {
		log.Println("Couldn't get Netbox results for Racks")
	}
	rackList = append(rackList, restRes...)
	log.Printf("Num of rack entries after Netbox query: %d and Result length = %d\n", len(rackList), len(restRes))
	Id2RackMap := make(map[int]*entity.Rack)
	for _, i := range restRes {
		Id2RackMap[int(i.Id)] = i
	}
	res, err := rs.ListRacksFromDB(ctx)
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
			log.Printf("The entry with Id %d is in Netbox\n", rack.Id)
			rk.Id = rack.Id
			if tmp.Name == "" {
				rk.Name = rack.Name // tmp.Name
			} else {
				rk.Name = tmp.Name
			}
			rk.Ipaddr = rack.Ipaddr // tmp.Ipaddr
			rk.Live = rack.Live
			rk.CustomFields = tmp.CustomFields
			rk.Created = tmp.Created
			rackList = append(rackList, &rk)
		}
	}
	log.Printf("Returning rackList with final length = %d\n", len(rackList))
	return rackList, nil
}

func (rs *repoSvc) GetServerByName(ctx context.Context, name string) (*entity.Server, error) {
	log.Println("I reached GetServer")
	var rk entity.Server
	// TODO - discuss this. This is a reasonable way to workaround when fields are directly flattened into a struct but also the type is
	// different. For ex, the table has int that we directly get. Then we make another query for fetching the status string by id.
	res, statusInt, err := rs.GetServerByNameFromDB(ctx, name)
	if err != nil {
		log.Println("DB Query for Server threw an error")
	}
	if res != nil {
		rk.Id = res.Id
		rk.HostName = res.HostName
		rk.PublicIpAddress = res.PublicIpAddress
	}
	log.Printf("DB result for Server %s is %+v\n", name, rk)

	restRes, err := rs.GetServerFromNetbox(ctx, rk.Id) // GetServer - Rest netbox.go
	if err != nil {
		log.Println("Rest Query for Server threw and error")
	}
	if restRes != nil {
		rk.RackName = restRes.RackName
		rk.NetboxName = restRes.NetboxName
	}
	 s, err := rs.GetServerStatusById(ctx, statusInt)
	 if err != nil {
	 	log.Printf("Error from getting server status by id: %v\n", err)
	 }
	 rk.Status = *s
	return &rk, nil
}

func (rs *repoSvc) GetServerById(ctx context.Context, id int) (*entity.Server, error) {
	log.Println("I reached GetServer")
	var rk entity.Server
	res, err := rs.GetServerByIdFromDB(ctx, id)
	if err != nil {
		log.Println("DB Query for Server threw an error")
	}
	if res != nil && res.Id != 0 {
		id = res.Id
	}
	restRes, err := rs.GetServerFromNetbox(ctx, id) // GetServer - Rest netbox.go
	if err != nil {
		log.Println("Rest Query for Server threw and error")
	}

	if restRes.RackName != "" {
		rk.RackName = restRes.RackName
	}
	return &rk, nil
}

func (rs *repoSvc) ListServers(ctx context.Context) ([]*entity.Server, error) {
	log.Println("I reached ListServers")
	serverList := make([]*entity.Server, 0)
	restRes, err := rs.ListServersFromNetbox(ctx)
	if err != nil {
		log.Println("Couldn't get Netbox results for Servers")
	}
	serverList = append(serverList, restRes...)
	log.Printf("Num of server entries after Netbox query: %d and Result length = %d\n", len(serverList), len(restRes))
	Id2ServerMap := make(map[int]*entity.Server)
	for _, i := range restRes {
		Id2ServerMap[int(i.Id)] = i
	}
	res, err := rs.ListServersFromDB(ctx)
	if err != nil {
		log.Println("Couldn't get DB results for Servers")
		return serverList, err
	}
	for _, server := range res {
		var rk entity.Server
		tmp, ok := Id2ServerMap[int(server.Id)]
		if !ok {
			// log.Println("Found an entry not in Netbox")
			serverList = append(serverList, server)
		} else {
			log.Printf("The entry with Id %d is in Netbox\n", server.Id)
			rk.Id = server.Id
			rk.RackName = tmp.RackName
			rk.HostName = server.HostName
			rk.NetboxName = tmp.NetboxName
			rk.PublicIpAddress = server.PublicIpAddress
			//rk.ServerStatus = rs.GetServerStatusById()
			// TODO Server Status resolution and dataloaders
			serverList = append(serverList, &rk)
		}
	}
	log.Printf("Returning serverList with final length = %d\n", len(serverList))
	return serverList, nil
}
