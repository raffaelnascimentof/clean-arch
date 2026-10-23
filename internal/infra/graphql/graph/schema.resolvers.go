package graph

import (
	"context"

	"github.com/devfullcycle/fc-clean-architecture/internal/infra/graphql/graph/model"
	"github.com/devfullcycle/fc-clean-architecture/internal/usecase"
)

// CreateOrder is the resolver for the createOrder mutation.
func (r *mutationResolver) CreateOrder(ctx context.Context, input model.CreateOrderInput) (*model.Order, error) {
	dto := usecase.CreateOrderInputDTO{
		Price: input.Price,
		Tax:   input.Tax,
	}
	if input.ID != nil {
		dto.ID = *input.ID
	}

	output, err := r.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}

	return &model.Order{
		ID:         output.ID,
		Price:      output.Price,
		Tax:        output.Tax,
		FinalPrice: output.FinalPrice,
	}, nil
}

// ListOrders is the resolver for the listOrders query.
func (r *queryResolver) ListOrders(ctx context.Context) ([]*model.Order, error) {
	orders, err := r.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	var result []*model.Order
	for _, o := range orders {
		result = append(result, &model.Order{
			ID:         o.ID,
			Price:      o.Price,
			Tax:        o.Tax,
			FinalPrice: o.FinalPrice,
		})
	}
	return result, nil
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
