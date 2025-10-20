package events

import (
	"testing"
	"time"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/value_objects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockHasher is a mock implementation of the Hasher interface
type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) Verify(hashedPassword, plainPassword string) bool {
	args := m.Called(hashedPassword, plainPassword)
	return args.Bool(0)
}

type UserCreatedEventTestSuite struct {
	suite.Suite
	mockHasher *MockHasher
	testUser   *entities.User
}

func (suite *UserCreatedEventTestSuite) SetupTest() {
	suite.mockHasher = new(MockHasher)
	suite.mockHasher.On("Hash", "SecurePass123").Return("hashed_password", nil)

	email, _ := value_objects.NewEmail("john.doe@example.com")
	password, _ := value_objects.NewPassword(suite.mockHasher, "SecurePass123")
	displayName := "Johnny"

	suite.testUser, _ = entities.NewUser(
		"tenant-123",
		"John Doe",
		email,
		password,
		&displayName,
	)
}

func TestUserCreatedEventTestSuite(t *testing.T) {
	suite.Run(t, new(UserCreatedEventTestSuite))
}

func (suite *UserCreatedEventTestSuite) TestNewUserCreatedEvent_CreatesEventCorrectly() {
	// Act
	event := NewUserCreatedEvent(suite.testUser)

	// Assert
	assert.NotEmpty(suite.T(), event.EventID())
	assert.Equal(suite.T(), "user.created", event.EventType())
	assert.Equal(suite.T(), suite.testUser.ID.String(), event.AggregateID())
	assert.False(suite.T(), event.OccurredOn().IsZero())
}

func (suite *UserCreatedEventTestSuite) TestNewUserCreatedEvent_MapsUserDataCorrectly() {
	// Act
	event := NewUserCreatedEvent(suite.testUser)

	// Assert
	assert.Equal(suite.T(), suite.testUser.ID, event.Data.ID)
	assert.Equal(suite.T(), suite.testUser.TenantID, event.Data.TenantID)
	assert.Equal(suite.T(), suite.testUser.Name, event.Data.Name)
	assert.Equal(suite.T(), suite.testUser.Email.Value(), event.Data.Email)
	assert.Equal(suite.T(), suite.testUser.DisplayName, event.Data.DisplayName)
	assert.NotEmpty(suite.T(), event.Data.CreatedAt)
}

func (suite *UserCreatedEventTestSuite) TestNewUserCreatedEvent_FormatsCreatedAtCorrectly() {
	// Act
	event := NewUserCreatedEvent(suite.testUser)

	// Assert - Verify RFC3339 format
	_, err := time.Parse(time.RFC3339, event.Data.CreatedAt)
	assert.NoError(suite.T(), err)

}

func (suite *UserCreatedEventTestSuite) TestNewUserCreatedEvent_WithoutDisplayName() {
	// Arrange
	email, _ := value_objects.NewEmail("jane.doe@example.com")
	mockHasher := new(MockHasher)
	mockHasher.On("Hash", "AnotherPass456").Return("hashed_password_2", nil)
	password, _ := value_objects.NewPassword(mockHasher, "AnotherPass456")

	user, _ := entities.NewUser(
		"tenant-456",
		"Jane Doe",
		email,
		password,
		nil,
	)

	// Act
	event := NewUserCreatedEvent(user)

	// Assert
	assert.Nil(suite.T(), event.Data.DisplayName)
	assert.Equal(suite.T(), user.ID, event.Data.ID)
	assert.Equal(suite.T(), user.Email.Value(), event.Data.Email)
}

func (suite *UserCreatedEventTestSuite) TestNewUserCreatedEvent_EventIDIsUnique() {
	// Act
	event1 := NewUserCreatedEvent(suite.testUser)
	event2 := NewUserCreatedEvent(suite.testUser)

	// Assert - Each event should have a unique event ID
	assert.NotEqual(suite.T(), event1.EventID(), event2.EventID())
}

func (suite *UserCreatedEventTestSuite) TestNewUserCreatedEvent_AggregateIDMatchesUserID() {
	// Act
	event := NewUserCreatedEvent(suite.testUser)

	// Assert
	assert.Equal(suite.T(), suite.testUser.ID.String(), event.AggregateID())
	assert.Equal(suite.T(), suite.testUser.ID, event.Data.ID)
}

func (suite *UserCreatedEventTestSuite) TestNewUserCreatedEvent_EventTypeIsConstant() {
	// Act
	event := NewUserCreatedEvent(suite.testUser)

	// Assert
	assert.Equal(suite.T(), "user.created", event.EventType())
}
