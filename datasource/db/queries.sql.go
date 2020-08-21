package db

import (
	"context"
	"log"

	"github.com/varunturlapati/vtgqlgen/pkg/entity"
)

const (
	getFruit = `-- name: GetFruit: one
SELECT id, name, quantity FROM fruits
WHERE id = ?
`
	listFruits = `-- name: ListFruits: many
SELECT id, name, quantity FROM fruits
ORDER BY id
`
	getDetail = `-- name: GetDetail: one
SELECT id, name, color, taste FROM fruitinfo
where name = ?
`
	listDetails = `-- name: ListDetails: many
SELECT id, name, color, taste FROM fruitinfo
ORDER BY name
`
	getLevel = `-- name: GetLevel: one
SELECT color, level FROM colorkey
where color = ?
`
	listLevels = `-- name: ListLevels: many
SELECT color, level FROM colorkey
ORDER BY color
`
	createFruit = `-- name: CreateFruit: one
INSERT INTO fruits (name, quantity)
VALUES (?, ?)
`
	updateFruit = `-- name: UpdateFruit: one
UPDATE fruits
SET name = ?, quantity = ?
WHERE id = ?
`
	deleteFruit = `-- name: DeleteFruit: one
DELETE FROM fruits
WHERE id = ?
`
	getRack = `-- name: GetRack: one
SELECT id, name, ipaddr, live FROM racks
WHERE id = ?
`
	listRacks = `-- name: ListRacks: many
SELECT id, name, ipaddr, live FROM racks
`
	// TODO - DB level handling of joins is a-OK?
	getServerByName = `-- name: GetServerByName: one
SELECT s.id, s.hostname, ss.name, s.publicipaddress FROM servers AS s, ServerStatus AS ss
WHERE s.serverstatus = ss.id AND s.hostname = ?
`
	getServerById = `-- name: GetServerById: one
SELECT id, hostname, serverstatus, publicipaddress FROM servers
WHERE id = ?
`
	// TODO - DB level handling of joins is a-OK?
	listServers = `-- name: ListServers: many
SELECT s.id, s.hostname, ss.name, s.publicipaddress FROM servers AS s, ServerStatus AS ss WHERE s.serverstatus = ss.id
`
	getServerStatusById = `-- name: GetServerStatusById: one
SELECT name FROM serverstatus
WHERE id = ?
`
)

type GetFruitParams struct {
	id int64
}

type GetDetailParams struct {
	name string
}

type GetLevelParams struct {
	color string
}

type Result struct {
	Count   int           `json:"count"`
	Results []entity.Rack `json:"results"`
}

func (q *Queries) GetFruit(ctx context.Context, id int) (*entity.Fruit, error) {
	row := q.db.QueryRowContext(ctx, getFruit, id)
	var f entity.Fruit
	err := row.Scan(&f.Id, &f.Name, &f.Quantity)
	return &f, err
}

func (q *Queries) ListFruits(ctx context.Context) ([]*entity.Fruit, error) {
	rows, err := q.db.QueryContext(ctx, listFruits)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var fs []*entity.Fruit
	for rows.Next() {
		var f entity.Fruit
		if err := rows.Scan(&f.Id, &f.Name, &f.Quantity); err != nil {
			return nil, err
		}
		fs = append(fs, &f)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return fs, err
}

func (q *Queries) CreateFruit(ctx context.Context, arg *entity.CreateFruitParams) (*entity.Fruit, error) {
	row := q.db.QueryRowContext(ctx, createFruit, arg.Name, arg.Quantity)
	var f entity.Fruit
	err := row.Scan(&f.Id, &f.Name, &f.Quantity)
	return &f, err
}

func (q *Queries) UpdateFruit(ctx context.Context, arg *entity.UpdateFruitParams) (*entity.Fruit, error) {
	row := q.db.QueryRowContext(ctx, updateFruit, arg.Name, arg.Quantity, arg.Id)
	var f entity.Fruit
	err := row.Scan(&f.Id, &f.Name, &f.Quantity)
	return &f, err
}

func (q *Queries) DeleteFruit(ctx context.Context, id int) (*entity.Fruit, error) {
	row := q.db.QueryRowContext(ctx, deleteFruit, id)
	var f entity.Fruit
	err := row.Scan(&f.Id, &f.Name, &f.Quantity)
	return &f, err
}

func (q *Queries) GetDetail(ctx context.Context, fruitName string) (*entity.Detail, error) {
	log.Println(fruitName)
	row := q.db.QueryRowContext(ctx, getDetail, fruitName)
	var f entity.Detail
	err := row.Scan(&f.Id, &f.Name, &f.Color, &f.Taste)
	log.Printf("Queries.GetDetail result: f is %+v\nerr is %v\n", f, err)
	return &f, err
}

func (q *Queries) ListDetails(ctx context.Context) ([]*entity.Detail, error) {
	rows, err := q.db.QueryContext(ctx, listDetails)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var fs []*entity.Detail
	for rows.Next() {
		var f entity.Detail
		if err := rows.Scan(&f.Id, &f.Name, &f.Color, &f.Taste); err != nil {
			return nil, err
		}
		fs = append(fs, &f)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return fs, err
}

func (q *Queries) GetLevel(ctx context.Context, color string) (*entity.Level, error) {
	row := q.db.QueryRowContext(ctx, getLevel, color)
	var f entity.Level
	err := row.Scan(&f.Color, &f.Level)
	return &f, err
}

func (q *Queries) ListLevels(ctx context.Context) ([]*entity.Level, error) {
	rows, err := q.db.QueryContext(ctx, listLevels)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var fs []*entity.Level
	for rows.Next() {
		var f entity.Level
		if err := rows.Scan(&f.Color, &f.Level); err != nil {
			return nil, err
		}
		fs = append(fs, &f)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return fs, err
}

func (q *Queries) GetRackFromDB(ctx context.Context, id int) (*entity.Rack, error) {
	row := q.db.QueryRowContext(ctx, getRack, id)
	var f entity.Rack
	err := row.Scan(&f.Id, &f.Name, &f.Ipaddr, &f.Live)
	return &f, err
}

func (q *Queries) ListRacksFromDB(ctx context.Context) ([]*entity.Rack, error) {
	rows, err := q.db.QueryContext(ctx, listRacks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var fs []*entity.Rack
	for rows.Next() {
		var f entity.Rack
		if err := rows.Scan(&f.Id, &f.Name, &f.Ipaddr, &f.Live); err != nil {
			return nil, err
		}
		fs = append(fs, &f)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return fs, err
}

func (q *Queries) GetServerByNameFromDB(ctx context.Context, name string) (*entity.Server, int, error) {
	row := q.db.QueryRowContext(ctx, getServerByName, name)
	var f entity.Server
	var status int
	err := row.Scan(&f.Id, &f.HostName, &status, &f.PublicIpAddress)
	log.Printf("Error from query is %v\n", err)
	return &f, status, err
}

func (q *Queries) GetServerByIdFromDB(ctx context.Context, id int) (*entity.Server, error) {
	row := q.db.QueryRowContext(ctx, getServerByName, id)
	var f entity.Server
	err := row.Scan(&f.Id, &f.HostName, &f.Status, &f.PublicIpAddress)
	return &f, err
}

func (q *Queries) ListServersFromDB(ctx context.Context) ([]*entity.Server, error) {
	rows, err := q.db.QueryContext(ctx, listServers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var fs []*entity.Server
	for rows.Next() {
		var f entity.Server
		if err := rows.Scan(&f.Id, &f.HostName, &f.Status, &f.PublicIpAddress); err != nil {
			return nil, err
		}
		fs = append(fs, &f)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return fs, err
}

type ListFruitsByRackIDsRow entity.Fruit

func (q *Queries) ListFruitsByRackIDs(ctx context.Context, rackIDs []int) ([]ListFruitsByRackIDsRow, error) {
	var retList []ListFruitsByRackIDsRow
	for _, rid := range rackIDs {
		res, err := q.GetFruit(ctx, rid)
		if err != nil {
			return nil, err
		}
		var tmp ListFruitsByRackIDsRow
		tmp.Id = res.Id
		tmp.Name = res.Name
		tmp.Quantity = res.Quantity
		retList = append(retList, tmp)
	}
	return retList, nil
}

func (q *Queries) GetServerStatusById(ctx context.Context, id int) (*string, error) {
	row := q.db.QueryRowContext(ctx, getServerStatusById, id)
	var status string
	err := row.Scan(&status)
	log.Printf("Status is %v, err is %v for Query: %s", status, err, getServerStatusById)
	return &status, err
}
