package usecase

import (
	"github.com/devfullcycle/fc-clean-architecture/internal/domain/entity"
	"github.com/devfullcycle/fc-clean-architecture/internal/domain/repository"
	"github.com/devfullcycle/fc-clean-architecture/pkg/events"
)

// CreateOrderInputDTO holds the data required to create an order.
type CreateOrderInputDTO struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

// CreateOrderOutputDTO holds the data returned after creating an order.
type CreateOrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

// CreateOrderUseCase handles the business logic for creating orders.
type CreateOrderUseCase struct {
	OrderRepository repository.OrderRepository
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

// NewCreateOrderUseCase instantiates a CreateOrderUseCase.
func NewCreateOrderUseCase(
	repo repository.OrderRepository,
	orderCreated events.EventInterface,
	dispatcher events.EventDispatcherInterface,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: repo,
		OrderCreated:    orderCreated,
		EventDispatcher: dispatcher,
	}
}

// Execute runs the create order use case.
func (c *CreateOrderUseCase) Execute(input CreateOrderInputDTO) (CreateOrderOutputDTO, error) {
	order, err := entity.NewOrder(input.ID, input.Price, input.Tax)
	if err != nil {
		return CreateOrderOutputDTO{}, err
	}

	if err := c.OrderRepository.Save(order); err != nil {
		return CreateOrderOutputDTO{}, err
	}

	output := CreateOrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}

	c.OrderCreated.SetPayload(output)
	if err := c.EventDispatcher.Dispatch(c.OrderCreated); err != nil {
		return CreateOrderOutputDTO{}, err
	}

	return output, nil
}
