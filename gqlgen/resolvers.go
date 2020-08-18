package gqlgen

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	"github.com/varunturlapati/vtgqlgen/dataloaders"
	"github.com/varunturlapati/vtgqlgen/datasource"
	"github.com/varunturlapati/vtgqlgen/pkg/entity"
)

type Resolver struct{
	Repository datasource.Repository
	DataLoaders dataloaders.Retriever
}

func (r *fruitResolver) Detail(ctx context.Context, obj *entity.Fruit) (*entity.Detail, error) {
	return r.Repository.GetDetail(ctx, obj.Name)
}

func (r *fruitResolver) Level(ctx context.Context, obj *entity.Fruit) (*entity.Level, error) {
	return r.Repository.GetLevel(ctx, obj.Name)
}

func (r *fruitResolver) Rack(ctx context.Context, obj *entity.Fruit) (*entity.Rack, error) {
	//return r.Repository.GetRack(ctx, obj.Id)
	return r.DataLoaders.Retrieve(ctx).RackByFruitId.Load(obj.Id)
}

func (r *queryResolver) Fruits(ctx context.Context) ([]entity.Fruit, error) {
	fruitPtrs, err := r.Repository.ListFruits(ctx)
	if err != nil {
		return nil, err
	}
	var fruits []entity.Fruit
	for _, p := range fruitPtrs {
		fruits = append(fruits, *p)
	}
	return fruits, nil
}

func (r *queryResolver) Fruit(ctx context.Context, id int) (*entity.Fruit, error) {
	return r.Repository.GetFruit(ctx, id)
}

func (r *queryResolver) Racks(ctx context.Context) ([]entity.Rack, error) {
	rackPtrs, err := r.Repository.ListRacks(ctx)
	if err != nil {
		return nil, err
	}
	var racks []entity.Rack
	for _, p := range rackPtrs {
		racks = append(racks, *p)
	}
	return racks, nil
}

func (r *queryResolver) Rack(ctx context.Context, id int) (*entity.Rack, error) {
	return r.Repository.GetRack(ctx, id)
}

func (r *rackResolver) Fruit(ctx context.Context, obj *entity.Rack) (*entity.Fruit, error) {
	// return r.Repository.GetFruit(ctx, int(obj.Id))
	return r.DataLoaders.Retrieve(ctx).FruitByRackId.Load(int(obj.Id))
}

// Fruit returns FruitResolver implementation.
func (r *Resolver) Fruit() FruitResolver { return &fruitResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Rack returns RackResolver implementation.
func (r *Resolver) Rack() RackResolver { return &rackResolver{r} }

type fruitResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type rackResolver struct{ *Resolver }
