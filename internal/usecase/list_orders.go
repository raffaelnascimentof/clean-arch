package usecase

import (
	"github.com/devfullcycle/fc-clean-architecture/internal/domain/repository"
)

// OrderOutputDTO holds the data for a single order in the list response.
type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

// ListOrdersUseCase handles the business logic for listing all orders.
type ListOrdersUseCase struct {
	OrderRepository repository.OrderRepository
}

// NewListOrdersUseCase instantiates a ListOrdersUseCase.
func NewListOrdersUseCase(repo repository.OrderRepository) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: repo,
	}
}

// Execute retrieves all orders from the repository and returns them as DTOs.
func (l *ListOrdersUseCase) Execute() ([]OrderOutputDTO, error) {
	orders, err := l.OrderRepository.FindAll()
	if err != nil {
		return nil, err
	}

	var output []OrderOutputDTO
	for _, order := range orders {
		output = append(output, OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	// Return empty slice instead of nil for clean JSON/gRPC responses
	if output == nil {
		output = []OrderOutputDTO{}
	}

	return output, nil
}
