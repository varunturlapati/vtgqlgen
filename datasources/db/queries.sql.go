package db

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	Count   int          `json:"count"`
	Results []ServerRack `json:"results"`
}

func (q *Queries) GetFruit(ctx context.Context, id int) (Fruit, error) {
	row := q.db.QueryRowContext(ctx, getFruit, id)
	var f Fruit
	err := row.Scan(&f.Id, &f.Name, &f.Quantity)
	return f, err
}

func (q *Queries) ListFruits(ctx context.Context) ([]Fruit, error) {
	rows, err := q.db.QueryContext(ctx, listFruits)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var fs []Fruit
	for rows.Next() {
		var f Fruit
		if err := rows.Scan(&f.Id, &f.Name, &f.Quantity); err != nil {
			return nil, err
		}
		fs = append(fs, f)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return fs, err
}

func (q *Queries) GetDetail(ctx context.Context, fruitName string) (Detail, error) {
	log.Println(fruitName)
	row := q.db.QueryRowContext(ctx, getDetail, fruitName)
	var f Detail
	err := row.Scan(&f.Id, &f.Name, &f.Color, &f.Taste)
	log.Printf("Queries.GetDetail result: f is %+v\nerr is %v\n", f, err)
	return f, err
}

func (q *Queries) ListDetails(ctx context.Context) ([]Detail, error) {
	rows, err := q.db.QueryContext(ctx, listDetails)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var fs []Detail
	for rows.Next() {
		var f Detail
		if err := rows.Scan(&f.Id, &f.Name, &f.Color, &f.Taste); err != nil {
			return nil, err
		}
		fs = append(fs, f)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return fs, err
}

func (q *Queries) GetLevel(ctx context.Context, color string) (Level, error) {
	row := q.db.QueryRowContext(ctx, getLevel, color)
	var f Level
	err := row.Scan(&f.Color, &f.Level)
	return f, err
}

func (q *Queries) ListLevels(ctx context.Context) ([]Level, error) {
	rows, err := q.db.QueryContext(ctx, listLevels)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var fs []Level
	for rows.Next() {
		var f Level
		if err := rows.Scan(&f.Color, &f.Level); err != nil {
			return nil, err
		}
		fs = append(fs, f)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return fs, err
}

func (q *Queries) GetRack(ctx context.Context, id int) (ServerRack, error) {
	var rack ServerRack
	resp, err := http.Get(fmt.Sprintf("http://localhost:8000/api/dcim/racks/%v/", id))
	if err != nil {
		return rack, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rack, err
	}
	err = json.Unmarshal(body, &rack)
	if err != nil {
		return rack, err
	}
	return rack, nil
}

func (q *Queries) ListRacks(ctx context.Context) ([]ServerRack, error) {
	var res Result
	var racks []ServerRack
	resp, err := http.Get("http://localhost:8000/api/dcim/racks/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	for _, elem := range res.Results {
		tmpRack := ServerRack{
			Id:   elem.Id,
			Name: elem.Name,
			CustomFields: CustomFields{
				RblxRackId:     elem.CustomFields.RblxRackId,
				DesignRevision: elem.CustomFields.DesignRevision,
				CageId:         elem.CustomFields.CageId,
			},
			Created:      elem.Created,
		}
		racks = append(racks, tmpRack)
	}
	return racks, nil
}