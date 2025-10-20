package entities_test

import (
	"testing"
	"time"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserReadTestSuite struct {
	suite.Suite
	testID        uuid.UUID
	testTenantID  string
	testName      string
	testEmail     string
	testCreatedAt string
}

func (suite *UserReadTestSuite) SetupTest() {
	suite.testID = uuid.New()
	suite.testTenantID = "tenant-123"
	suite.testName = "John Doe"
	suite.testEmail = "john.doe@example.com"
	suite.testCreatedAt = time.Now().Format(time.RFC3339)
}

func TestUserReadTestSuite(t *testing.T) {
	suite.Run(t, new(UserReadTestSuite))
}

func (suite *UserReadTestSuite) TestNewUserRead_WithDisplayName() {
	// Arrange
	displayName := "Johnny"

	// Act
	userRead := entities.NewUserRead(
		suite.testID,
		suite.testTenantID,
		suite.testName,
		suite.testEmail,
		&displayName,
		suite.testCreatedAt,
	)

	// Assert
	assert.NotNil(suite.T(), userRead)
	assert.Equal(suite.T(), suite.testID, userRead.ID)
	assert.Equal(suite.T(), suite.testTenantID, userRead.TenantID)
	assert.Equal(suite.T(), suite.testName, userRead.Name)
	assert.Equal(suite.T(), suite.testEmail, userRead.Email)
	assert.NotNil(suite.T(), userRead.DisplayName)
	assert.Equal(suite.T(), displayName, *userRead.DisplayName)
	assert.Equal(suite.T(), suite.testCreatedAt, userRead.CreatedAt)
}

func (suite *UserReadTestSuite) TestNewUserRead_WithoutDisplayName() {
	// Arrange & Act
	userRead := entities.NewUserRead(
		suite.testID,
		suite.testTenantID,
		suite.testName,
		suite.testEmail,
		nil,
		suite.testCreatedAt,
	)

	// Assert
	assert.NotNil(suite.T(), userRead)
	assert.Equal(suite.T(), suite.testID, userRead.ID)
	assert.Equal(suite.T(), suite.testTenantID, userRead.TenantID)
	assert.Equal(suite.T(), suite.testName, userRead.Name)
	assert.Equal(suite.T(), suite.testEmail, userRead.Email)
	assert.Nil(suite.T(), userRead.DisplayName)
	assert.Equal(suite.T(), suite.testCreatedAt, userRead.CreatedAt)
}

func (suite *UserReadTestSuite) TestUserRead_StructureFields() {
	// Arrange
	displayName := "Test Display"

	// Act
	userRead := &entities.UserRead{
		ID:          suite.testID,
		TenantID:    suite.testTenantID,
		Name:        suite.testName,
		Email:       suite.testEmail,
		DisplayName: &displayName,
		CreatedAt:   suite.testCreatedAt,
	}

	// Assert
	assert.Equal(suite.T(), suite.testID, userRead.ID)
	assert.Equal(suite.T(), suite.testTenantID, userRead.TenantID)
	assert.Equal(suite.T(), suite.testName, userRead.Name)
	assert.Equal(suite.T(), suite.testEmail, userRead.Email)
	assert.Equal(suite.T(), displayName, *userRead.DisplayName)
	assert.Equal(suite.T(), suite.testCreatedAt, userRead.CreatedAt)
}

func TestNewUserRead_EmptyValues(t *testing.T) {
	// Arrange
	id := uuid.New()
	emptyTenantID := ""
	emptyName := ""
	emptyEmail := ""
	emptyCreatedAt := ""

	// Act
	userRead := entities.NewUserRead(id, emptyTenantID, emptyName, emptyEmail, nil, emptyCreatedAt)

	// Assert
	assert.NotNil(t, userRead)
	assert.Equal(t, id, userRead.ID)
	assert.Equal(t, emptyTenantID, userRead.TenantID)
	assert.Equal(t, emptyName, userRead.Name)
	assert.Equal(t, emptyEmail, userRead.Email)
	assert.Nil(t, userRead.DisplayName)
	assert.Equal(t, emptyCreatedAt, userRead.CreatedAt)
}

func TestNewUserRead_WithLongDisplayName(t *testing.T) {
	// Arrange
	id := uuid.New()
	tenantID := "tenant-456"
	name := "Jane Doe"
	email := "jane.doe@example.com"
	longDisplayName := "This is a very long display name that might be used in some edge cases"
	createdAt := time.Now().Format(time.RFC3339)

	// Act
	userRead := entities.NewUserRead(id, tenantID, name, email, &longDisplayName, createdAt)

	// Assert
	assert.NotNil(t, userRead)
	assert.Equal(t, id, userRead.ID)
	assert.Equal(t, longDisplayName, *userRead.DisplayName)
}
