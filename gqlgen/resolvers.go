package gqlgen

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"
	"log"

	gql "github.com/graph-gophers/graphql-go"

	"github.com/varunturlapati/vtgqlgen/dataloaders"
	ds "github.com/varunturlapati/vtgqlgen/datasource"
	"github.com/varunturlapati/vtgqlgen/pkg/entity"
)

type Fruit struct {
	Id       gql.ID
	Name     string
	Quantity int32
	Color    string
	Level    string
	// X        Detail
	Rack entity.ServerRack
}

type Resolver struct {
	Repository ds.Repository
	DataLoaders dataloaders.Retriever
}

func (r *queryResolver) Fruits(ctx context.Context) ([]entity.Fruit, error) {
	fruits, err := r.Repository.ListFruits(ctx)
	if err != nil {
		return nil, err
	}
	return fruits, nil
	/*
		details, err := r.Repository.ListDetails(ctx)
		if err != nil {
			return nil, err
		}
	*/
	// panic("not implemented")
}

func (r *queryResolver) Fruit(ctx context.Context, id int) (*entity.Fruit, error) {
	fObj, err := r.Repository.GetFruit(ctx, id)
	if err != nil {
		return nil, err
	}
	return &fObj, nil
	// panic("not implemented")
}

func (r *Resolver) Fruit() FruitResolver { return &fruitResolver{r} }

// func (r *Resolver) Detail() DetailResolver { return &detailResolver{r}}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

type fruitResolver struct {
	*Resolver
}

type detailResolver struct {
	*Resolver
}

type levelResolver struct {
	*Resolver
}

type rackResolver struct {
	*Resolver
}

func (r *fruitResolver) Detail(ctx context.Context, obj *entity.Fruit) (*entity.Detail, error) {
	//d := &detailResolver{r.Resolver}
	//detail, err := d.Repository.GetDetail(ctx, obj.Name)
	detail, err := r.Repository.GetDetail(ctx, obj.Name)
	log.Println(err)
	return &detail, err
}

func (r *fruitResolver) Level(ctx context.Context, obj *entity.Fruit) (*entity.Level, error) {
	l := &levelResolver{r.Resolver}
	level, err := l.Repository.GetLevel(ctx, obj.Level)
	return &level, err
}

func (r *fruitResolver) ID(ctx context.Context, obj *entity.Fruit) (int, error) {
	return obj.Id, nil
}

func (r *fruitResolver) Color(ctx context.Context, obj *entity.ServerRack) (*detailResolver, error) {
	return &detailResolver{r.Resolver}, nil
}

func (r *detailResolver) Color(ctx context.Context, obj *entity.Detail) (string, error) {
	return obj.Color, nil
}

func (r *detailResolver) Taste(ctx context.Context, obj *entity.Detail) (string, error) {
	return obj.Taste, nil
}

func (r *fruitResolver) Rack(ctx context.Context, obj *entity.Fruit) (*entity.ServerRack, error) {
	/*
	a := &rackResolver{r.Resolver}
	rack, err := a.Repository.GetRack(ctx, obj.Id)
	retRack := &Rack{
		ID:           int(rack.Id),
		Name:         rack.Name,
		Created:      &rack.Created,
		CustomFields: &rack.CustomFields,
	}
	return retRack, err
	*/
	return r.DataLoaders.Retrieve(ctx).RackByFruitID.Load(obj.Id)
}
