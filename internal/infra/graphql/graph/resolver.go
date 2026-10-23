package graph

import "github.com/devfullcycle/fc-clean-architecture/internal/usecase"

// Resolver holds use cases injected into GraphQL resolvers.
type Resolver struct {
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}
