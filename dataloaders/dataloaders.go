package dataloaders

//go:generate go run github.com/vektah/dataloaden FruitLoader int *github.com/varunturlapati/vtgqlgen/pkg/entity.Fruit
//go:generate go run github.com/vektah/dataloaden RackLoader int *github.com/varunturlapati/vtgqlgen/pkg/entity.Rack

import (
	"context"
	"github.com/varunturlapati/vtgqlgen/pkg/entity"
	"time"

	ds "github.com/varunturlapati/vtgqlgen/datasource"
)

type contextKey string

const key = contextKey("dataloaders")

type Loaders struct {
	// define individual loaders here
	RackByFruitID *RackLoader
}

func newLoaders(ctx context.Context, repo ds.Repository) *Loaders {
	return &Loaders{
		RackByFruitID: newRackByFruitID(ctx, repo),
	}
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
		Fetch:    func(fruitIDs []int) ([]*entity.ServerRack, []error) {
			res, err := repo.ListRacksByFruitIDs(ctx, fruitIDs)
			if err != nil {
				return nil, []error{err}
			}
			groupByFruitID := make(map[int]*entity.ServerRack, len(fruitIDs))
			for _, r := range res {
				groupByFruitID[r.Id] = &entity.ServerRack{
					Id:           int64(r.Id),
					Name:         r.Name,
					CustomFields: entity.CustomFields{
						RblxRackId: r.CustomFields.RblxRackId,
						DesignRevision: r.CustomFields.DesignRevision,
					},
					Created:      r.Created,
				}
			}
			result := make([]*entity.ServerRack, len(fruitIDs))
			for i, fruitID := range fruitIDs {
				result[i] = groupByFruitID[fruitID]
			}
			return result, nil
		},
		Wait:     5 * time.Millisecond,
		MaxBatch: 10,
	})
}