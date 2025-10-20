package entities_test

import (
	"testing"
	"time"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/value_objects"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

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

type UserTestSuite struct {
	suite.Suite
	mockHasher   *MockHasher
	testTenantID string
	testName     string
	testEmail    value_objects.Email
	testPassword value_objects.Password
}

func (suite *UserTestSuite) SetupTest() {
	suite.mockHasher = new(MockHasher)
	suite.testTenantID = "tenant-123"
	suite.testName = "John Doe"

	var err error
	suite.testEmail, err = value_objects.NewEmail("john.doe@example.com")
	assert.NoError(suite.T(), err)

	suite.mockHasher.On("Hash", "SecurePass123").Return("hashed_password", nil)
	suite.testPassword, err = value_objects.NewPassword(suite.mockHasher, "SecurePass123")
	assert.NoError(suite.T(), err)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (suite *UserTestSuite) TestNewUser_WithDisplayName() {
	// Arrange
	displayName := "Johnny"

	// Act
	user, err := entities.NewUser(
		suite.testTenantID,
		suite.testName,
		suite.testEmail,
		suite.testPassword,
		&displayName,
	)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.NotEqual(suite.T(), uuid.Nil, user.ID)
	assert.Equal(suite.T(), suite.testTenantID, user.TenantID)
	assert.Equal(suite.T(), suite.testName, user.Name)
	assert.Equal(suite.T(), suite.testEmail, user.Email)
	assert.Equal(suite.T(), suite.testPassword, user.Password)
	assert.NotNil(suite.T(), user.DisplayName)
	assert.Equal(suite.T(), displayName, *user.DisplayName)
	assert.False(suite.T(), user.CreatedAt.IsZero())
	assert.False(suite.T(), user.UpdatedAt.IsZero())
	assert.True(suite.T(), user.CreatedAt.Equal(user.UpdatedAt))
}

func (suite *UserTestSuite) TestNewUser_WithoutDisplayName() {
	// Arrange & Act
	user, err := entities.NewUser(
		suite.testTenantID,
		suite.testName,
		suite.testEmail,
		suite.testPassword,
		nil,
	)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.NotEqual(suite.T(), uuid.Nil, user.ID)
	assert.Equal(suite.T(), suite.testTenantID, user.TenantID)
	assert.Equal(suite.T(), suite.testName, user.Name)
	assert.Equal(suite.T(), suite.testEmail, user.Email)
	assert.Equal(suite.T(), suite.testPassword, user.Password)
	assert.Nil(suite.T(), user.DisplayName)
	assert.False(suite.T(), user.CreatedAt.IsZero())
	assert.False(suite.T(), user.UpdatedAt.IsZero())
}

func (suite *UserTestSuite) TestNewUser_GeneratesUniqueIDs() {
	// Act
	user1, err1 := entities.NewUser(suite.testTenantID, suite.testName, suite.testEmail, suite.testPassword, nil)
	user2, err2 := entities.NewUser(suite.testTenantID, suite.testName, suite.testEmail, suite.testPassword, nil)

	// Assert
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.NotEqual(suite.T(), user1.ID, user2.ID)
}

func (suite *UserTestSuite) TestNewUser_SetsTimestamps() {
	// Arrange
	beforeCreation := time.Now()

	// Act
	time.Sleep(1 * time.Millisecond) // Pequeña pausa para asegurar diferencia de tiempo
	user, err := entities.NewUser(suite.testTenantID, suite.testName, suite.testEmail, suite.testPassword, nil)
	time.Sleep(1 * time.Millisecond)
	afterCreation := time.Now()

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), user.CreatedAt.After(beforeCreation))
	assert.True(suite.T(), user.CreatedAt.Before(afterCreation))
	assert.True(suite.T(), user.UpdatedAt.After(beforeCreation))
	assert.True(suite.T(), user.UpdatedAt.Before(afterCreation))
}

func TestNewUser_WithEmptyTenantID(t *testing.T) {
	// Arrange
	mockHasher := new(MockHasher)
	email, _ := value_objects.NewEmail("test@example.com")
	mockHasher.On("Hash", "SecurePass123").Return("hashed_password", nil)
	password, _ := value_objects.NewPassword(mockHasher, "SecurePass123")

	// Act
	user, err := entities.NewUser("", "Test User", email, password, nil)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "", user.TenantID)
}

func TestNewUser_WithEmptyName(t *testing.T) {
	// Arrange
	mockHasher := new(MockHasher)
	email, _ := value_objects.NewEmail("test@example.com")
	mockHasher.On("Hash", "SecurePass123").Return("hashed_password", nil)
	password, _ := value_objects.NewPassword(mockHasher, "SecurePass123")

	// Act
	user, err := entities.NewUser("tenant-123", "", email, password, nil)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "", user.Name)
}

func TestNewUser_WithLongDisplayName(t *testing.T) {
	// Arrange
	mockHasher := new(MockHasher)
	email, _ := value_objects.NewEmail("test@example.com")
	mockHasher.On("Hash", "SecurePass123").Return("hashed_password", nil)
	password, _ := value_objects.NewPassword(mockHasher, "SecurePass123")
	longDisplayName := "This is a very long display name that might be used in some edge cases to test the system limits"

	// Act
	user, err := entities.NewUser("tenant-123", "Test User", email, password, &longDisplayName)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, longDisplayName, *user.DisplayName)
}

func TestUser_StructFields(t *testing.T) {
	// Arrange
	mockHasher := new(MockHasher)
	id := uuid.New()
	tenantID := "tenant-456"
	name := "Jane Doe"
	email, _ := value_objects.NewEmail("jane.doe@example.com")
	mockHasher.On("Hash", "AnotherPass456").Return("hashed_password_2", nil)
	password, _ := value_objects.NewPassword(mockHasher, "AnotherPass456")
	displayName := "Janie"
	createdAt := time.Now()
	updatedAt := time.Now()

	// Act
	user := &entities.User{
		ID:          id,
		TenantID:    tenantID,
		Name:        name,
		Email:       email,
		Password:    password,
		DisplayName: &displayName,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	// Assert
	assert.Equal(t, id, user.ID)
	assert.Equal(t, tenantID, user.TenantID)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, password, user.Password)
	assert.Equal(t, displayName, *user.DisplayName)
	assert.Equal(t, createdAt, user.CreatedAt)
	assert.Equal(t, updatedAt, user.UpdatedAt)
}
