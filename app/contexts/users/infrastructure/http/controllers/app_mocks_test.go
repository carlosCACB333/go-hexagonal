package controllers_test

import (
	"context"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/commands"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/queries"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/value_objects"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/infrastructure/http/controllers"
	shared_ports "github.com/carloscacb333/go-hexagonal/app/shared/domain/ports"
	shared_middleware "github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockUserRepository struct{ mock.Mock }

func (m *MockUserRepository) Save(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, tenantID string, id uuid.UUID) (*entities.User, error) {
	args := m.Called(ctx, tenantID, id)
	if v := args.Get(0); v != nil {
		return v.(*entities.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, tenantID string, email value_objects.Email) (*entities.User, error) {
	args := m.Called(ctx, tenantID, email)
	if v := args.Get(0); v != nil {
		return v.(*entities.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, tenantID string, email value_objects.Email) (bool, error) {
	args := m.Called(ctx, tenantID, email)
	return args.Bool(0), args.Error(1)
}

type MockUserReadRepository struct{ mock.Mock }

func (m *MockUserReadRepository) FindByID(ctx context.Context, tenantID string, id uuid.UUID) (*entities.UserRead, error) {
	args := m.Called(ctx, tenantID, id)
	if v := args.Get(0); v != nil {
		return v.(*entities.UserRead), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserReadRepository) Upsert(ctx context.Context, dto *entities.UserRead) error {
	args := m.Called(ctx, dto)
	return args.Error(0)
}

type MockIdempotencyRepository struct{ mock.Mock }

func (m *MockIdempotencyRepository) IsProcessed(ctx context.Context, tenantID, key string) (bool, error) {
	args := m.Called(ctx, tenantID, key)
	return args.Bool(0), args.Error(1)
}
func (m *MockIdempotencyRepository) MarkAsProcessed(ctx context.Context, tenantID, key string) error {
	args := m.Called(ctx, tenantID, key)
	return args.Error(0)
}

type MockEventBus struct{ mock.Mock }

func (m *MockEventBus) Publish(ctx context.Context, event shared_ports.DomainEvent, correlationID string) error {
	args := m.Called(ctx, event, correlationID)
	return args.Error(0)
}
func (m *MockEventBus) Close() { m.Called() }

type MockHasher struct{ mock.Mock }

func (m *MockHasher) Hash(plainPassword string) (string, error) {
	args := m.Called(plainPassword)
	return args.String(0), args.Error(1)
}
func (m *MockHasher) Verify(hashedPassword, plainPassword string) bool {
	args := m.Called(hashedPassword, plainPassword)
	return args.Bool(0)
}

type deps struct {
	userRepo     *MockUserRepository
	userReadRepo *MockUserReadRepository
	idem         *MockIdempotencyRepository
	bus          *MockEventBus
	hasher       *MockHasher
}

func setupAppWithDeps(d deps) *fiber.App {

	createUseCase := commands.NewCreateUserUseCase(d.userRepo, d.idem, d.bus, d.hasher)
	getUseCase := queries.NewGetUserUseCase(d.userReadRepo)

	app := fiber.New(
		fiber.Config{ErrorHandler: shared_middleware.ErrorHandler(zap.NewNop())},
	)
	app.Use(shared_middleware.TenantMiddleware())
	app.Use(shared_middleware.CorrelationIDMiddleware())
	app.Post("/users", controllers.CreateUserController(createUseCase))
	app.Get("/users/:id", controllers.GetUserController(getUseCase))
	return app
}
