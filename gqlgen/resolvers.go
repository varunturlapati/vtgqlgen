package gqlgen

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"
	"github.com/varunturlapati/vtgqlgen/dataloaders"
	"github.com/varunturlapati/vtgqlgen/datasource"
	"github.com/varunturlapati/vtgqlgen/pkg/entity"
	"log"
)

type Resolver struct {
	Repository  datasource.Repository
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
	log.Printf("Calling fruitResolver.Rack() with dataloader")
	return r.DataLoaders.Retrieve(ctx).RackByFruitId.Load(obj.Id)
}

func (r *queryResolver) Fruits(ctx context.Context, fFilter *entity.FruitFilter, rFilter *entity.RackFilter) ([]entity.Fruit, error) {
	fruitPtrs, err := r.Repository.ListFruits(ctx, fFilter, rFilter)
	if err != nil {
		return nil, err
	}
	// netboxFruits, err := r.Repository.ListRacks(ctx, rFilter)
	var fruits []entity.Fruit
	for _, p := range fruitPtrs {
		fruits = append(fruits, *p)
	}
	return fruits, nil
}

func (r *queryResolver) Fruit(ctx context.Context, id int) (*entity.Fruit, error) {
	return r.Repository.GetFruit(ctx, id)
}

func (r *queryResolver) Racks(ctx context.Context, rFilter *entity.RackFilter) ([]entity.Rack, error) {
	rackPtrs, err := r.Repository.ListRacks(ctx, rFilter)
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

func (r *queryResolver) Servers(ctx context.Context) ([]entity.Server, error) {
	srvPtrs, err := r.Repository.ListServers(ctx)
	if err != nil {
		return nil, err
	}
	var servers []entity.Server
	for _, p := range srvPtrs {
		servers = append(servers, *p)
	}
	return servers, nil
}

func (r *queryResolver) Server(ctx context.Context, name *string, id *int) (*entity.Server, error) {
	if name == nil && id == nil {
		return nil, nil
	}
	if name != nil && *name != "" {
		return r.Repository.GetServerByName(ctx, *name)
	}
	if *id != 0 {
		return r.Repository.GetServerById(ctx, *id)
	}
	return nil, nil
}

func (r *rackResolver) Fruit(ctx context.Context, obj *entity.Rack) (*entity.Fruit, error) {
	// return r.Repository.GetFruit(ctx, int(obj.Id))
	return r.DataLoaders.Retrieve(ctx).FruitByRackId.Load(int(obj.Id))
}

// Fruit returns FruitResolver implementation.
func (r *Resolver) Fruit() FruitResolver { return &fruitResolver{r} }

// Rack returns RackResolver implementation.
func (r *Resolver) Rack() RackResolver { return &rackResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Mutation returns an implementation of the MutationResolver interface.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

func (r *mutationResolver) CreateFruit(ctx context.Context, data FruitInput) (*entity.Fruit, error) {
	fruit, err := r.Repository.CreateFruit(ctx, &entity.CreateFruitParams{
		Name:     data.Name,
		Quantity: data.Quantity,
	})
	if err != nil {
		return nil, err
	}
	return fruit, nil
}

func (r *mutationResolver) UpdateFruit(ctx context.Context, id int, data FruitInput) (*entity.Fruit, error) {
	fruit, err := r.Repository.UpdateFruit(ctx, &entity.UpdateFruitParams{
		Id:       id,
		Name:     data.Name,
		Quantity: data.Quantity,
	})
	if err != nil {
		return nil, err
	}
	return fruit, nil
}

func (r *mutationResolver) DeleteFruit(ctx context.Context, id int) (*entity.Fruit, error) {
	fruit, err := r.Repository.DeleteFruit(ctx, id)
	if err != nil {
		return nil, err
	}
	return fruit, nil
}

type fruitResolver struct{ *Resolver }
type rackResolver struct{ *Resolver }
type serverResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

func (r *queryResolver) ServerByAttrs(ctx context.Context, attrs *entity.ServerAttrs) (*entity.Server, error) {
	if attrs == nil {
		return nil, nil
	}
	return r.Repository.GetServerByAttrs(ctx, attrs)
}

type mutationResolver struct{ *Resolver }
