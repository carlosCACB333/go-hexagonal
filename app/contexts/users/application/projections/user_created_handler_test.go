package projections_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/projections"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/events"
	shared_events "github.com/carloscacb333/go-hexagonal/app/shared/domain/events"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

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

type UserCreatedHandlerSuite struct {
	suite.Suite
	repo      *MockUserReadRepository
	handler   *projections.UserCreatedHandler
	ctx       context.Context
	eventData *events.UserCreatedEvent
}

func (s *UserCreatedHandlerSuite) SetupTest() {
	s.repo = new(MockUserReadRepository)
	s.handler = projections.NewUserCreatedHandler(s.repo)
	s.ctx = context.Background()
	user := entities.NewUserRead(uuid.New(), "tenant-1", "John Doe", "john@example.com", nil, time.Now().Format(time.RFC3339))
	s.eventData = &events.UserCreatedEvent{
		BaseEvent: shared_events.NewBaseEvent(
			"user.created", user.ID.String(),
		),
		Data: *user,
	}
}

func TestUserCreatedHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserCreatedHandlerSuite))
}

func (s *UserCreatedHandlerSuite) TestExecute_Success() {
	// Arrange
	s.repo.On("Upsert", mock.Anything, &s.eventData.Data).Return(nil).Once()

	// Act
	err := s.handler.Handle(s.ctx, s.eventData)

	// Assert
	assert.NoError(s.T(), err)
	s.repo.AssertExpectations(s.T())
}

func (s *UserCreatedHandlerSuite) TestExecute_RepoErrorWrapped() {
	// Arrange
	s.repo.On("Upsert", mock.Anything, &s.eventData.Data).Return(errors.New("db error")).Once()

	// Act
	err := s.handler.Handle(s.ctx, s.eventData)

	// Assert
	if assert.Error(s.T(), err) {
		assert.Contains(s.T(), err.Error(), "failed to create user read model")
	}
	s.repo.AssertExpectations(s.T())
}

func (s *UserCreatedHandlerSuite) TestExecute_NilEvent_ReturnsError() {
	// Act
	err := s.handler.Handle(s.ctx, nil)

	// Assert
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to create user read model")

}
