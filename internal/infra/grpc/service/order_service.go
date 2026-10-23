package service

import (
	"context"

	"github.com/devfullcycle/fc-clean-architecture/internal/infra/grpc/pb"
	"github.com/devfullcycle/fc-clean-architecture/internal/usecase"
)

// OrderService implements the gRPC OrderServiceServer interface.
type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}

// NewOrderService creates a new gRPC OrderService.
func NewOrderService(
	createOrderUseCase usecase.CreateOrderUseCase,
	listOrdersUseCase usecase.ListOrdersUseCase,
) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

// CreateOrder handles the gRPC CreateOrder RPC.
func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.CreateOrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}

	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

// ListOrders handles the gRPC ListOrders RPC.
func (s *OrderService) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := s.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	var pbOrders []*pb.OrderResponse
	for _, o := range orders {
		pbOrders = append(pbOrders, &pb.OrderResponse{
			Id:         o.ID,
			Price:      float32(o.Price),
			Tax:        float32(o.Tax),
			FinalPrice: float32(o.FinalPrice),
		})
	}

	return &pb.ListOrdersResponse{Orders: pbOrders}, nil
}
