package dataloaders

//go:generate go run github.com/vektah/dataloaden FruitLoader int *github.com/varunturlapati/vtgqlgen/pkg/entity.Fruit
//go:generate go run github.com/vektah/dataloaden RackLoader int *github.com/varunturlapati/vtgqlgen/pkg/entity.Rack

import (
	"context"
	"time"

	ds "github.com/varunturlapati/vtgqlgen/datasource"
	"github.com/varunturlapati/vtgqlgen/pkg/entity"
)

type contextKey string

const key = contextKey("dataloaders")

// Loaders holds references to the individual dataloaders.
type Loaders struct {
	// individual loaders will be defined here
	RackByFruitId      *RackLoader
	FruitByRackId      *FruitLoader
	RackFromAllSources *RackLoader
	RackByFilters      *RackLoader
}

func newLoaders(ctx context.Context, repo ds.Repository) *Loaders {
	return &Loaders{
		// individual loaders will be initialized here
		FruitByRackId:      newFruitByRackID(ctx, repo),
		RackByFruitId:      newRackByFruitID(ctx, repo),
		RackFromAllSources: newRackByID(ctx, repo),
		// RackByFilters: newRackByFilters(ctx, repo),
	}
}

func newRackByID(ctx context.Context, repo ds.Repository) *RackLoader {
	return NewRackLoader(RackLoaderConfig{
		Fetch: func(IDs []int) ([]*entity.Rack, []error) {
			res, err := repo.ListRacksByIDs(ctx, IDs)
			if err != nil {
				return nil, []error{err}
			}
			groupByRackID := make(map[int]*entity.Rack, len(IDs))
			for _, r := range res {
				groupByRackID[int(r.Id)] = &entity.Rack{
					Id:   r.Id,
					Name: r.Name,
					CustomFields: entity.CustomFields{
						RblxRackId:     r.CustomFields.RblxRackId,
						DesignRevision: r.CustomFields.DesignRevision,
					},
					Created: r.Created,
					Ipaddr:  r.Ipaddr,
					Live:    r.Live,
				}
			}
			result := make([]*entity.Rack, len(IDs))
			for i, fruitID := range IDs {
				result[i] = groupByRackID[fruitID]
			}
			return result, nil
		},
		Wait:     5 * time.Millisecond,
		MaxBatch: 10,
	})
}

// Retriever retrieves dataloaders from the request context.
type Retriever interface {
	Retrieve(context.Context) *Loaders
}

type retriever struct {
	key contextKey
}

func (r *retriever) Retrieve(ctx context.Context) *Loaders {
	return ctx.Value(r.key).(*Loaders)
}

// NewRetriever instantiates a new implementation of Retriever.
func NewRetriever() Retriever {
	return &retriever{key: key}
}

func newRackByFruitID(ctx context.Context, repo ds.Repository) *RackLoader {
	return NewRackLoader(RackLoaderConfig{
		Fetch: func(fruitIDs []int) ([]*entity.Rack, []error) {
			res, err := repo.ListRacksByFruitIDs(ctx, fruitIDs)
			if err != nil {
				return nil, []error{err}
			}
			groupByFruitID := make(map[int]*entity.Rack, len(fruitIDs))
			for _, r := range res {
				groupByFruitID[int(r.Id)] = &entity.Rack{
					Id:   int64(r.Id),
					Name: r.Name,
					CustomFields: entity.CustomFields{
						RblxRackId:     r.CustomFields.RblxRackId,
						DesignRevision: r.CustomFields.DesignRevision,
					},
					Created: r.Created,
				}
			}
			result := make([]*entity.Rack, len(fruitIDs))
			for i, fruitID := range fruitIDs {
				result[i] = groupByFruitID[fruitID]
			}
			return result, nil
		},
		Wait:     5 * time.Millisecond,
		MaxBatch: 10,
	})
}

/*
func newRackByFilters(ctx context.Context, repo ds.Repository) *RackLoader {
	return NewRackLoader(RackLoaderConfig{
		Fetch: func(rFilter *entity.RackFilter) ([]*entity.Rack, []error) {
			res, err := repo.ListRacksByFilters(ctx, rFilter)
			if err != nil {
				return nil, []error{err}
			}
			groupByFruitID := make(map[int]*entity.Rack, len(fruitIDs))
			for _, r := range res {
				groupByFruitID[int(r.Id)] = &entity.Rack{
					Id:   int64(r.Id),
					Name: r.Name,
					CustomFields: entity.CustomFields{
						RblxRackId:     r.CustomFields.RblxRackId,
						DesignRevision: r.CustomFields.DesignRevision,
					},
					Created: r.Created,
				}
			}
			result := make([]*entity.Rack, len(fruitIDs))
			for i, fruitID := range fruitIDs {
				result[i] = groupByFruitID[fruitID]
			}
			return result, nil
		},
		Wait:     5 * time.Millisecond,
		MaxBatch: 10,
	})
}
*/

func newFruitByRackID(ctx context.Context, repo ds.Repository) *FruitLoader {
	return NewFruitLoader(FruitLoaderConfig{
		Fetch: func(rackIDs []int) ([]*entity.Fruit, []error) {
			res, err := repo.ListFruitsByRackIDs(ctx, rackIDs)
			if err != nil {
				return nil, []error{err}
			}
			groupByRackID := make(map[int]*entity.Fruit, len(rackIDs))
			for _, r := range res {
				groupByRackID[r.Id] = &entity.Fruit{
					Id:       r.Id,
					Name:     r.Name,
					Quantity: r.Quantity,
				}
			}
			result := make([]*entity.Fruit, len(rackIDs))
			for i, fruitID := range rackIDs {
				result[i] = groupByRackID[fruitID]
			}
			return result, nil
		},
		Wait:     5 * time.Millisecond,
		MaxBatch: 10,
	})
}
