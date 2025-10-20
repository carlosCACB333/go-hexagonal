package queries_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/queries"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/exceptions"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockUserReadRepository struct {
	mock.Mock
}

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

// Test Suite

type GetUserUseCaseSuite struct {
	suite.Suite
	repo *MockUserReadRepository
	uc   *queries.GetUserUseCase
	ctx  context.Context
}

func (s *GetUserUseCaseSuite) SetupTest() {
	s.repo = new(MockUserReadRepository)
	s.uc = queries.NewGetUserUseCase(s.repo)
	s.ctx = context.Background()
}

func TestGetUserUseCaseSuite(t *testing.T) {
	suite.Run(t, new(GetUserUseCaseSuite))
}

func (s *GetUserUseCaseSuite) TestExecute_Success() {
	// Arrange
	userID := uuid.New()
	tenantID := "tenant-123"
	displayName := "Johnny"
	expectedUser := entities.NewUserRead(
		userID,
		tenantID,
		"John Doe",
		"john.doe@example.com",
		&displayName,
		time.Now().Format(time.RFC3339),
	)

	query := queries.GetUserQuery{
		TenantID: tenantID,
		UserID:   userID,
	}

	s.repo.On("FindByID", mock.Anything, tenantID, userID).Return(expectedUser, nil).Once()

	// Act
	result, err := s.uc.Execute(s.ctx, query)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedUser.ID, result.ID)
	assert.Equal(s.T(), expectedUser.TenantID, result.TenantID)
	assert.Equal(s.T(), expectedUser.Name, result.Name)
	assert.Equal(s.T(), expectedUser.Email, result.Email)
	assert.Equal(s.T(), expectedUser.DisplayName, result.DisplayName)
	s.repo.AssertExpectations(s.T())
}

func (s *GetUserUseCaseSuite) TestExecute_UserNotFound() {
	// Arrange
	userID := uuid.New()
	tenantID := "tenant-456"

	query := queries.GetUserQuery{
		TenantID: tenantID,
		UserID:   userID,
	}

	s.repo.On("FindByID", mock.Anything, tenantID, userID).Return(nil, errors.New("not found")).Once()

	// Act
	result, err := s.uc.Execute(s.ctx, query)

	// Assert
	assert.Nil(s.T(), result)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), exceptions.ErrUserNotFound, err)
	s.repo.AssertExpectations(s.T())
}

func (s *GetUserUseCaseSuite) TestExecute_RepositoryError() {
	// Arrange
	userID := uuid.New()
	tenantID := "tenant-789"

	query := queries.GetUserQuery{
		TenantID: tenantID,
		UserID:   userID,
	}

	s.repo.On("FindByID", mock.Anything, tenantID, userID).Return(nil, errors.New("db connection error")).Once()

	// Act
	result, err := s.uc.Execute(s.ctx, query)

	// Assert
	assert.Nil(s.T(), result)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), exceptions.ErrUserNotFound, err)
	s.repo.AssertExpectations(s.T())
}

func (s *GetUserUseCaseSuite) TestExecute_WithoutDisplayName() {
	// Arrange
	userID := uuid.New()
	tenantID := "tenant-abc"
	expectedUser := entities.NewUserRead(
		userID,
		tenantID,
		"Jane Doe",
		"jane.doe@example.com",
		nil,
		time.Now().Format(time.RFC3339),
	)

	query := queries.GetUserQuery{
		TenantID: tenantID,
		UserID:   userID,
	}

	s.repo.On("FindByID", mock.Anything, tenantID, userID).Return(expectedUser, nil).Once()

	// Act
	result, err := s.uc.Execute(s.ctx, query)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Nil(s.T(), result.DisplayName)
	assert.Equal(s.T(), expectedUser.ID, result.ID)
	s.repo.AssertExpectations(s.T())
}

func (s *GetUserUseCaseSuite) TestExecute_DifferentTenants() {
	// Arrange
	userID1 := uuid.New()
	userID2 := uuid.New()
	tenant1 := "tenant-1"
	tenant2 := "tenant-2"

	user1 := entities.NewUserRead(userID1, tenant1, "User 1", "user1@example.com", nil, time.Now().Format(time.RFC3339))
	user2 := entities.NewUserRead(userID2, tenant2, "User 2", "user2@example.com", nil, time.Now().Format(time.RFC3339))

	s.repo.On("FindByID", mock.Anything, tenant1, userID1).Return(user1, nil).Once()
	s.repo.On("FindByID", mock.Anything, tenant2, userID2).Return(user2, nil).Once()

	// Act
	result1, err1 := s.uc.Execute(s.ctx, queries.GetUserQuery{TenantID: tenant1, UserID: userID1})
	result2, err2 := s.uc.Execute(s.ctx, queries.GetUserQuery{TenantID: tenant2, UserID: userID2})

	// Assert
	assert.NoError(s.T(), err1)
	assert.NoError(s.T(), err2)
	assert.Equal(s.T(), tenant1, result1.TenantID)
	assert.Equal(s.T(), tenant2, result2.TenantID)
	assert.NotEqual(s.T(), result1.ID, result2.ID)
	s.repo.AssertExpectations(s.T())
}
