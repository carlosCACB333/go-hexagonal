package commands_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/commands"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
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

type CreateUserReadUseCaseSuite struct {
	suite.Suite
	repo *MockUserReadRepository
	uc   *commands.CreateUserReadUseCase
	ctx  context.Context
}

func (s *CreateUserReadUseCaseSuite) SetupTest() {
	s.repo = new(MockUserReadRepository)
	s.uc = commands.NewCreateUserReadUseCase(s.repo)
	s.ctx = context.Background()
}

func TestCreateUserReadUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CreateUserReadUseCaseSuite))
}

func (s *CreateUserReadUseCaseSuite) TestExecute_Success() {
	// Arrange
	display := "Johnny"
	user := entities.NewUserRead(uuid.New(), "tenant-1", "John Doe", "john@example.com", &display, time.Now().Format(time.RFC3339))
	s.repo.On("Upsert", mock.Anything, user).Return(nil).Once()

	// Act
	result, err := s.uc.Execute(s.ctx, user)

	// Assert
	assert.NoError(s.T(), err)
	assert.Same(s.T(), user, result)
	s.repo.AssertExpectations(s.T())
}

func (s *CreateUserReadUseCaseSuite) TestExecute_RepoErrorWrapped() {
	// Arrange
	user := entities.NewUserRead(uuid.New(), "tenant-2", "Jane Doe", "jane@example.com", nil, time.Now().Format(time.RFC3339))
	s.repo.On("Upsert", mock.Anything, user).Return(errors.New("db error")).Once()

	// Act
	result, err := s.uc.Execute(s.ctx, user)

	// Assert
	assert.Nil(s.T(), result)
	if assert.Error(s.T(), err) {
		assert.Contains(s.T(), err.Error(), "failed to create user read model")
	}
	s.repo.AssertExpectations(s.T())
}

func (s *CreateUserReadUseCaseSuite) TestExecute_NilUser_PassesThrough() {
	// Arrange
	s.repo.On("Upsert", mock.Anything, (*entities.UserRead)(nil)).Return(nil).Once()

	// Act
	result, err := s.uc.Execute(s.ctx, nil)

	// Assert
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), result)
	s.repo.AssertExpectations(s.T())
}
