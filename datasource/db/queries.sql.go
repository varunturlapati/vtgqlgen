package db

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/varunturlapati/vtgqlgen/pkg/entity"
)

const (
	getFruit = `-- name: GetFruit: one
SELECT id, name, quantity FROM fruits
WHERE id = ?
`
	listFruits = `-- name: ListFruits: many
SELECT id, name, quantity FROM fruits
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
	getServerByAttrsBaseQuery = `-- name: GetServerByAttrs: one
SELECT s.id, s.hostname, ss.name, s.publicipaddress FROM servers AS s, ServerStatus AS ss
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

func (q *Queries) ListFruits(ctx context.Context, fruitFilter *entity.FruitFilter, rackFilter *entity.RackFilter) ([]*entity.Fruit, error) {
	var query string
	var clauses []string
	var args []interface{}
	if fruitFilter != nil {
		if fruitFilter.Ids != nil {
			processClausesFromIdFilter(fruitFilter.Ids, &clauses, &args)
		}
		if fruitFilter.Names != nil {
			processClausesFromStringFilter(fruitFilter.Names, &clauses, &args)
		}
	}

	if len(clauses) > 0 {
		query = fmt.Sprintf("%s WHERE %s", listFruits, strings.Join(clauses, " AND "))
	} else {
		query = listFruits
	}

	// log.Printf("The final query after id and name filters is: %s\n", query)
	rows, err := q.db.QueryContext(ctx, query, args...)
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

func processClausesFromStringFilter(filter *entity.StringFilter, allClauses *[]string, allArgs *[]interface{}) {
	var clauses []string
	var args []interface{}

	if filter.StartsWith != "" {
		clauses = append(clauses, "name like ?")
		args = append(args, filter.StartsWith+"%")
	}
	if filter.EndsWith != "" {
		clauses = append(clauses, "name like ?")
		args = append(args, "%"+filter.EndsWith)
	}
	if filter.Contains != "" {
		clauses = append(clauses, "name like ?")
		args = append(args, "%"+filter.Contains+"%")
	}
	if filter.NotContain != "" {
		clauses = append(clauses, "name not like ?")
		args = append(args, "%"+filter.NotContain+"%")
	}
	if len(clauses) >= 1 {
		// qBytes.WriteString(strings.Join(clauses, " AND "))
		*allArgs = append(*allArgs, args...)
		*allClauses = append(*allClauses, clauses...)
	}
}

func processClausesFromIdFilter(filter *entity.IntFilter, allClauses *[]string, allArgs *[]interface{}) {
	var clauses []string
	var args []interface{}
	// var qBytes bytes.Buffer

	if filter.Ne != 0 {
		clauses = append(clauses, "id != ?")
		args = append(args, filter.Ne)
	}
	if filter.Ge != 0 {
		clauses = append(clauses, "id >= ?")
		args = append(args, filter.Ge)
	}
	if filter.Le != 0 {
		clauses = append(clauses, "id <= ?")
		args = append(args, filter.Le)
	}
	if filter.Gt != 0 {
		clauses = append(clauses, "id > ?")
		args = append(args, filter.Gt)
	}
	if filter.Lt != 0 {
		clauses = append(clauses, "id < ?")
		args = append(args, filter.Lt)
	}
	if len(clauses) >= 1 {
		// qBytes.WriteString(strings.Join(clauses, " AND "))
		*allArgs = append(*allArgs, args...)
		*allClauses = append(*allClauses, clauses...)
	}
	// retStr := qBytes.String()
	// log.Printf("IdFilter procesing yields the following query: %s with args: %v", retStr, args)
	// return retStr, args
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

func (q *Queries) GetServerByNameFromDB(ctx context.Context, name string) (*entity.Server, error) {
	row := q.db.QueryRowContext(ctx, getServerByName, name)
	var f entity.Server
	err := row.Scan(&f.Id, &f.HostName, &f.Status, &f.PublicIpAddress)
	log.Printf("Error from query is %v\n", err)
	return &f, err
}

func (q *Queries) GetServerByIdFromDB(ctx context.Context, id int) (*entity.Server, error) {
	row := q.db.QueryRowContext(ctx, getServerByName, id)
	var f entity.Server
	err := row.Scan(&f.Id, &f.HostName, &f.Status, &f.PublicIpAddress)
	return &f, err
}

func (q *Queries) GetServerByAttrsFromDB(ctx context.Context, attrs *entity.ServerAttrs) (*entity.Server, error) {
	var clauses []string
	var args []interface{}
	var query bytes.Buffer
	query.WriteString(getServerByAttrsBaseQuery)
	if attrs.HostName != "" {
		clauses = append(clauses, "s.Hostname = ?")
		args = append(args, attrs.HostName)
	}
	if attrs.Status != "" {
		clauses = append(clauses, "ss.Name = ?")
		args = append(args, attrs.Status)
	}
	if len(clauses) >= 1 {
		query.WriteString(fmt.Sprintf(" WHERE %s", strings.Join(clauses, " AND ")))
	} else {
		query.WriteString(" LIMIT 1")
	}
	/*
		finQ := fmt.Sprintf("%s", query.String())
		log.Printf("%s with args %v", finQ, args)
		row := q.db.QueryRowContext(ctx, finQ, args...)
	*/
	row := q.db.QueryRowContext(ctx, query.String(), args...)
	var s entity.Server
	err := row.Scan(&s.Id, &s.HostName, &s.Status, &s.PublicIpAddress)
	return &s, err
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
