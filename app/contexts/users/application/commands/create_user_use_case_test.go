package commands_test

import (
	"context"
	"errors"
	"testing"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/commands"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
	user_exceptions "github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/exceptions"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/value_objects"
	shared_ports "github.com/carloscacb333/go-hexagonal/app/shared/domain/ports"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mocks

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

// Test Suite

type CreateUserUseCaseSuite struct {
	suite.Suite
	repo   *MockUserRepository
	idem   *MockIdempotencyRepository
	event  *MockEventBus
	hasher *MockHasher
	uc     *commands.CreateUserUseCase
	ctx    context.Context
}

func (s *CreateUserUseCaseSuite) SetupTest() {
	s.repo = new(MockUserRepository)
	s.idem = new(MockIdempotencyRepository)
	s.event = new(MockEventBus)
	s.hasher = new(MockHasher)
	s.uc = commands.NewCreateUserUseCase(s.repo, s.idem, s.event, s.hasher)
	s.ctx = context.Background()
}

func TestCreateUserUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CreateUserUseCaseSuite))
}

func (s *CreateUserUseCaseSuite) TestExecute_Success() {
	cmd := commands.CreateUserCommand{
		TenantID:       "tenant-1",
		IdempotencyKey: "key-1",
		CorrelationID:  "corr-123",
		Name:           "John Doe",
		Email:          "john@example.com",
		Password:       "StrongPass1",
	}

	email, _ := value_objects.NewEmail(cmd.Email)
	// Idempotency check
	s.idem.On("IsProcessed", mock.Anything, cmd.TenantID, cmd.IdempotencyKey).Return(false, nil).Once()
	// Email unique check
	s.repo.On("ExistsByEmail", mock.Anything, cmd.TenantID, email).Return(false, nil).Once()
	// Password hashing
	s.hasher.On("Hash", cmd.Password).Return("hashed_pwd", nil).Once()

	// Capture saved user
	var savedUser *entities.User
	s.repo.On("Save", mock.Anything, mock.AnythingOfType("*entities.User")).Return(nil).Run(func(args mock.Arguments) {
		savedUser = args.Get(1).(*entities.User)
	}).Once()

	// Event publish
	s.event.On("Publish", mock.Anything, mock.Anything, cmd.CorrelationID).Return(nil).Run(func(args mock.Arguments) {
		evt := args.Get(1).(shared_ports.DomainEvent)
		s.Require().NotNil(savedUser)
		s.Equal("user.created", evt.EventType())
		s.Equal(savedUser.ID.String(), evt.AggregateID())
	}).Once()

	// Mark as processed
	s.idem.On("MarkAsProcessed", mock.Anything, cmd.TenantID, cmd.IdempotencyKey).Return(nil).Once()

	resp, err := s.uc.Execute(s.ctx, cmd)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), savedUser.ID, resp.UserID)

	s.repo.AssertExpectations(s.T())
	s.idem.AssertExpectations(s.T())
	s.event.AssertExpectations(s.T())
	s.hasher.AssertExpectations(s.T())
}

func (s *CreateUserUseCaseSuite) TestExecute_IdempotencyAlreadyProcessed() {
	cmd := commands.CreateUserCommand{TenantID: "tenant-1", IdempotencyKey: "dup-key", Email: "john@example.com", Password: "StrongPass1"}
	s.idem.On("IsProcessed", mock.Anything, cmd.TenantID, cmd.IdempotencyKey).Return(true, nil).Once()

	resp, err := s.uc.Execute(s.ctx, cmd)
	assert.Nil(s.T(), resp)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "already processed")

	s.repo.AssertNotCalled(s.T(), "Save", mock.Anything, mock.Anything)
	s.event.AssertNotCalled(s.T(), "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func (s *CreateUserUseCaseSuite) TestExecute_IdempotencyCheckError() {
	cmd := commands.CreateUserCommand{TenantID: "tenant-1", IdempotencyKey: "key-err", Email: "john@example.com", Password: "StrongPass1"}
	s.idem.On("IsProcessed", mock.Anything, cmd.TenantID, cmd.IdempotencyKey).Return(false, errors.New("db down")).Once()

	resp, err := s.uc.Execute(s.ctx, cmd)
	assert.Nil(s.T(), resp)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to check idempotency")

	s.repo.AssertNotCalled(s.T(), "Save", mock.Anything, mock.Anything)
	s.event.AssertNotCalled(s.T(), "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func (s *CreateUserUseCaseSuite) TestExecute_InvalidEmail() {
	cmd := commands.CreateUserCommand{TenantID: "t", Email: "not-an-email", Password: "StrongPass1"}

	resp, err := s.uc.Execute(s.ctx, cmd)
	assert.Nil(s.T(), resp)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), user_exceptions.ErrInvalidEmail, err)

	s.repo.AssertNotCalled(s.T(), "Save", mock.Anything, mock.Anything)
}

func (s *CreateUserUseCaseSuite) TestExecute_ExistsByEmailError() {
	cmd := commands.CreateUserCommand{TenantID: "t1", Email: "john@example.com", Password: "StrongPass1"}
	email, _ := value_objects.NewEmail(cmd.Email)

	s.idem.On("IsProcessed", mock.Anything, cmd.TenantID, "").Return(false, nil).Once()
	s.repo.On("ExistsByEmail", mock.Anything, cmd.TenantID, email).Return(false, errors.New("db error")).Once()

	resp, err := s.uc.Execute(s.ctx, cmd)
	assert.Nil(s.T(), resp)
	assert.Error(s.T(), err)
}

func (s *CreateUserUseCaseSuite) TestExecute_DuplicateEmail() {
	cmd := commands.CreateUserCommand{TenantID: "t1", Email: "john@example.com", Password: "StrongPass1"}
	email, _ := value_objects.NewEmail(cmd.Email)

	s.idem.On("IsProcessed", mock.Anything, cmd.TenantID, "").Return(false, nil).Once()
	s.repo.On("ExistsByEmail", mock.Anything, cmd.TenantID, email).Return(true, nil).Once()

	resp, err := s.uc.Execute(s.ctx, cmd)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), user_exceptions.ErrDuplicateEmail, err)
}

func (s *CreateUserUseCaseSuite) TestExecute_WeakPassword() {
	cmd := commands.CreateUserCommand{TenantID: "t1", Email: "john@example.com", Password: "weak"}

	s.idem.On("IsProcessed", mock.Anything, cmd.TenantID, "").Return(false, nil).Once()
	email, _ := value_objects.NewEmail(cmd.Email)
	s.repo.On("ExistsByEmail", mock.Anything, cmd.TenantID, email).Return(false, nil).Once()

	resp, err := s.uc.Execute(s.ctx, cmd)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), user_exceptions.ErrWeakPassword, err)
	// Ensure hasher.Hash was not called
	s.hasher.AssertNotCalled(s.T(), "Hash", mock.Anything)
}

func (s *CreateUserUseCaseSuite) TestExecute_HasherError() {
	cmd := commands.CreateUserCommand{TenantID: "t1", Email: "john@example.com", Password: "StrongPass1"}

	s.idem.On("IsProcessed", mock.Anything, cmd.TenantID, "").Return(false, nil).Once()
	email, _ := value_objects.NewEmail(cmd.Email)
	s.repo.On("ExistsByEmail", mock.Anything, cmd.TenantID, email).Return(false, nil).Once()
	s.hasher.On("Hash", cmd.Password).Return("", errors.New("hash failed")).Once()

	resp, err := s.uc.Execute(s.ctx, cmd)
	assert.Nil(s.T(), resp)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to hash password")
}

func (s *CreateUserUseCaseSuite) TestExecute_SaveError() {
	cmd := commands.CreateUserCommand{TenantID: "t1", Email: "john@example.com", Password: "StrongPass1"}

	s.idem.On("IsProcessed", mock.Anything, cmd.TenantID, "").Return(false, nil).Once()
	email, _ := value_objects.NewEmail(cmd.Email)
	s.repo.On("ExistsByEmail", mock.Anything, cmd.TenantID, email).Return(false, nil).Once()
	s.hasher.On("Hash", cmd.Password).Return("hashed", nil).Once()
	s.repo.On("Save", mock.Anything, mock.AnythingOfType("*entities.User")).Return(errors.New("save err")).Once()

	resp, err := s.uc.Execute(s.ctx, cmd)
	assert.Nil(s.T(), resp)
	assert.Error(s.T(), err)
}

func (s *CreateUserUseCaseSuite) TestExecute_MarkAsProcessedError() {
	cmd := commands.CreateUserCommand{TenantID: "t1", IdempotencyKey: "key-1", Email: "john@example.com", Password: "StrongPass1"}

	s.idem.On("IsProcessed", mock.Anything, cmd.TenantID, cmd.IdempotencyKey).Return(false, nil).Once()
	email, _ := value_objects.NewEmail(cmd.Email)
	s.repo.On("ExistsByEmail", mock.Anything, cmd.TenantID, email).Return(false, nil).Once()
	s.hasher.On("Hash", cmd.Password).Return("hashed", nil).Once()
	s.repo.On("Save", mock.Anything, mock.AnythingOfType("*entities.User")).Return(nil).Once()
	s.idem.On("MarkAsProcessed", mock.Anything, cmd.TenantID, cmd.IdempotencyKey).Return(errors.New("mark err")).Once()

	resp, err := s.uc.Execute(s.ctx, cmd)
	assert.Nil(s.T(), resp)
	assert.Error(s.T(), err)
	// Publish should not be called if marking as processed fails
	s.event.AssertNotCalled(s.T(), "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func (s *CreateUserUseCaseSuite) TestExecute_PublishError() {
	cmd := commands.CreateUserCommand{TenantID: "t1", IdempotencyKey: "key-1", CorrelationID: "corr-x", Email: "john@example.com", Password: "StrongPass1"}

	s.idem.On("IsProcessed", mock.Anything, cmd.TenantID, cmd.IdempotencyKey).Return(false, nil).Once()
	email, _ := value_objects.NewEmail(cmd.Email)
	s.repo.On("ExistsByEmail", mock.Anything, cmd.TenantID, email).Return(false, nil).Once()
	s.hasher.On("Hash", cmd.Password).Return("hashed", nil).Once()
	s.repo.On("Save", mock.Anything, mock.AnythingOfType("*entities.User")).Return(nil).Once()
	s.idem.On("MarkAsProcessed", mock.Anything, cmd.TenantID, cmd.IdempotencyKey).Return(nil).Once()
	s.event.On("Publish", mock.Anything, mock.Anything, cmd.CorrelationID).Return(errors.New("bus down")).Once()

	resp, err := s.uc.Execute(s.ctx, cmd)
	assert.Nil(s.T(), resp)
	assert.Error(s.T(), err)
}
