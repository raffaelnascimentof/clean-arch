package graph

import (
	"github.com/devfullcycle/fc-clean-architecture/internal/usecase"
)

type Resolver struct {
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}
