package db

import (
	"context"
	"github.com/varunturlapati/vtgqlgen/pkg/entity"
	"log"
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
WHERE id = ?`
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
	Count   int                 `json:"count"`
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

