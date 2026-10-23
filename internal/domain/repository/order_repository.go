package repository

import "github.com/devfullcycle/fc-clean-architecture/internal/domain/entity"

// OrderRepository defines the contract for order persistence operations.
type OrderRepository interface {
	Save(order *entity.Order) error
	GetTotal() (int, error)
	FindAll() ([]entity.Order, error)
}
