package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserController_Success(t *testing.T) {
	d := deps{userReadRepo: new(MockUserReadRepository)}
	app := setupAppWithDeps(d)

	// expectations
	userID := uuid.New()
	tenantID := "tenant-123"
	display := "Johnny"

	expectedUser := entities.NewUserRead(
		userID,
		tenantID,
		"John Doe",
		"john.doe@example.com",
		&display,
		time.Now().Format(time.RFC3339),
	)

	d.userReadRepo.On("FindByID", mock.Anything, tenantID, userID).Return(expectedUser, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)
	req.Header.Set("X-Tenant-Id", tenantID)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var payload map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&payload)

	assert.Equal(t, userID.String(), payload["id"])
	assert.Equal(t, "John Doe", payload["name"])
	assert.Equal(t, "john.doe@example.com", payload["email"])
	d.userReadRepo.AssertExpectations(t)
}

func TestGetUserController_InvalidUUID(t *testing.T) {
	d := deps{userReadRepo: new(MockUserReadRepository)}
	app := setupAppWithDeps(d)

	tenantID := "tenant-123"

	req := httptest.NewRequest(http.MethodGet, "/users/invalid-uuid", nil)
	req.Header.Set("X-Tenant-Id", tenantID)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Expect Bad Request due to invalid UUID
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var payload map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&payload)

	// The error handler returns ApiError JSON with fields code/message
	assert.Equal(t, float64(http.StatusBadRequest), payload["code"]) // numbers decode as float64
	assert.Equal(t, "invalid user id", payload["message"])

	d.userReadRepo.AssertExpectations(t)
}

func TestGetUserController_NotFound(t *testing.T) {
	d := deps{userReadRepo: new(MockUserReadRepository)}
	app := setupAppWithDeps(d)

	// expectations
	tenantID := "tenant-123"
	userID := uuid.New()

	// Simulate repository signaling not found via error
	d.userReadRepo.On("FindByID", mock.Anything, tenantID, userID).Return(nil, assert.AnError).Once()

	req := httptest.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)
	req.Header.Set("X-Tenant-Id", tenantID)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Expect Not Found
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var payload map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&payload)

	assert.Equal(t, float64(http.StatusNotFound), payload["code"]) // numbers decode as float64
	assert.Equal(t, "user not found", payload["message"])

	d.userReadRepo.AssertExpectations(t)
}
