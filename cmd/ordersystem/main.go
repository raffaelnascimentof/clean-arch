package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/devfullcycle/fc-clean-architecture/configs"
	"github.com/devfullcycle/fc-clean-architecture/internal/domain/event"
	"github.com/devfullcycle/fc-clean-architecture/internal/infra/database"
	"github.com/devfullcycle/fc-clean-architecture/internal/infra/graphql/graph"
	"github.com/devfullcycle/fc-clean-architecture/internal/infra/grpc/pb"
	"github.com/devfullcycle/fc-clean-architecture/internal/infra/grpc/service"
	"github.com/devfullcycle/fc-clean-architecture/internal/infra/web"
	"github.com/devfullcycle/fc-clean-architecture/internal/infra/web/webserver"
	"github.com/devfullcycle/fc-clean-architecture/internal/usecase"
	"github.com/devfullcycle/fc-clean-architecture/pkg/events"
	"github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := configs.LoadConfig()

	// Connect to MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&multiStatements=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open(cfg.DBDriver, dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("Connected to database successfully")

	// Run migrations
	if err := runMigrations(db, cfg.DBName); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Dependency setup
	orderRepository := database.NewOrderRepository(db)
	orderCreatedEvent := event.NewOrderCreated()
	eventDispatcher := events.NewEventDispatcher()

	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository, orderCreatedEvent, eventDispatcher)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)

	// Start HTTP server
	orderHandler := web.NewOrderHandler(createOrderUseCase, listOrdersUseCase)
	ws := webserver.NewWebServer(cfg.WebPort)
	ws.AddHandler("POST", "/order", orderHandler.Create)
	ws.AddHandler("GET", "/order", orderHandler.List)

	go func() {
		log.Printf("Starting REST server on port %s\n", cfg.WebPort)
		ws.Start()
	}()

	// Start gRPC server
	grpcServer := grpc.NewServer()
	orderGrpcService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderGrpcService)
	reflection.Register(grpcServer)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
		if err != nil {
			log.Fatalf("failed to listen on gRPC port: %v", err)
		}
		log.Printf("Starting gRPC server on port %s\n", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Start GraphQL server
	gqlSrv := graphql_handler.NewDefaultServer(
		graph.NewExecutableSchema(graph.Config{
			Resolvers: &graph.Resolver{
				CreateOrderUseCase: *createOrderUseCase,
				ListOrdersUseCase:  *listOrdersUseCase,
			},
		}),
	)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", gqlSrv)

	log.Printf("Starting GraphQL server on port %s\n", cfg.GraphQLPort)
	log.Printf("GraphQL playground available at http://localhost%s\n", cfg.GraphQLPort)
	if err := http.ListenAndServe(cfg.GraphQLPort, nil); err != nil {
		log.Fatalf("failed to serve GraphQL: %v", err)
	}
}

func runMigrations(db *sql.DB, dbName string) error {
	driver, err := migratemysql.WithInstance(db, &migratemysql.Config{
		DatabaseName: dbName,
	})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://sql/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}
