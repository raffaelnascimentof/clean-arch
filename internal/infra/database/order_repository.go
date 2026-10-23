package database

import (
	"database/sql"

	"github.com/devfullcycle/fc-clean-architecture/internal/domain/entity"
	"github.com/devfullcycle/fc-clean-architecture/internal/domain/repository"
)

// OrderRepository is a MySQL-backed implementation of repository.OrderRepository.
type OrderRepository struct {
	DB *sql.DB
}

// NewOrderRepository creates a new OrderRepository.
func NewOrderRepository(db *sql.DB) repository.OrderRepository {
	return &OrderRepository{DB: db}
}

// Save inserts a new order into the database.
func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.DB.Prepare(
		"INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	return err
}

// GetTotal returns the total count of orders in the database.
func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&total)
	return total, err
}

// FindAll retrieves all orders from the database.
func (r *OrderRepository) FindAll() ([]entity.Order, error) {
	rows, err := r.DB.Query("SELECT id, price, tax, final_price FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var o entity.Order
		if err := rows.Scan(&o.ID, &o.Price, &o.Tax, &o.FinalPrice); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
