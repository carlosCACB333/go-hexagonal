package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUserController_Success(t *testing.T) {
	d := deps{userRepo: new(MockUserRepository), idem: new(MockIdempotencyRepository), bus: new(MockEventBus), hasher: new(MockHasher)}
	app := setupAppWithDeps(d)

	// expectations
	d.idem.On("IsProcessed", mock.Anything, "tenant-1", "idem-1").Return(false, nil).Once()
	// ExistsByEmail expects a value_objects.Email, but we can match any and assert later
	d.userRepo.On("ExistsByEmail", mock.Anything, "tenant-1", mock.Anything).Return(false, nil).Once()
	d.hasher.On("Hash", "StrongPass1").Return("hashed_pwd", nil).Once()
	d.userRepo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()
	d.idem.On("MarkAsProcessed", mock.Anything, "tenant-1", "idem-1").Return(nil).Once()
	d.bus.On("Publish", mock.Anything, mock.Anything, "corr-1").Return(nil).Once()

	body := `{"name":"John Doe","email":"john@example.com","password":"StrongPass1","display_name":"Johnny"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", "tenant-1")
	req.Header.Set("X-Correlation-Id", "corr-1")
	req.Header.Set("X-Idempotency-Key", "idem-1")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var payload map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&payload)
	assert.NotEmpty(t, payload["user_id"]) // user ID generated inside use case
	assert.Equal(t, "User created successfully", payload["message"])

	d.userRepo.AssertExpectations(t)
	d.idem.AssertExpectations(t)
	d.bus.AssertExpectations(t)
	d.hasher.AssertExpectations(t)
}

func TestCreateUserController_InvalidJSON(t *testing.T) {
	d := deps{userRepo: new(MockUserRepository), idem: new(MockIdempotencyRepository), bus: new(MockEventBus), hasher: new(MockHasher)}
	app := setupAppWithDeps(d)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", "tenant-1")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var payload map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&payload)
	assert.Equal(t, "invalid request body", payload["message"])
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// No calls should be made to use case dependencies
	d.userRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
}

func TestCreateUserController_UseCaseReturnsApiError(t *testing.T) {
	d := deps{userRepo: new(MockUserRepository), idem: new(MockIdempotencyRepository), bus: new(MockEventBus), hasher: new(MockHasher)}
	app := setupAppWithDeps(d)

	// Simulate duplicate email via ExistsByEmail = true
	d.idem.On("IsProcessed", mock.Anything, "tenant-1", "").Return(false, nil).Once()
	d.userRepo.On("ExistsByEmail", mock.Anything, "tenant-1", mock.Anything).Return(true, nil).Once()

	body := `{"name":"John Doe","email":"john@example.com","password":"StrongPass1"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", "tenant-1")

	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
}

func TestCreateUserController_MissingTenantHeader(t *testing.T) {
	d := deps{userRepo: new(MockUserRepository), idem: new(MockIdempotencyRepository), bus: new(MockEventBus), hasher: new(MockHasher)}
	app := setupAppWithDeps(d)

	body := `{"name":"John Doe","email":"john@example.com","password":"StrongPass1"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	// No X-Tenant-Id header

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var payload map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&payload)

	assert.Equal(t, "X-Tenant-Id header is required", payload["error"])
	d.userRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
}

func TestCreateUserController_CorrelationGeneratedWhenMissing(t *testing.T) {
	d := deps{userRepo: new(MockUserRepository), idem: new(MockIdempotencyRepository), bus: new(MockEventBus), hasher: new(MockHasher)}
	app := setupAppWithDeps(d)

	d.userRepo.On("ExistsByEmail", mock.Anything, "tenant-1", mock.Anything).Return(false, nil).Once()
	d.hasher.On("Hash", "StrongPass1").Return("hashed_pwd", nil).Once()
	d.userRepo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()
	// MarkAsProcessed not called as no idempotency key
	d.bus.On("Publish", mock.Anything, mock.Anything, mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		corr := args.String(2)
		if corr == "" {
			t.Errorf("expected correlation id to be generated")
		}
	}).Once()

	body := `{"name":"Jane","email":"jane@example.com","password":"StrongPass1"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", "tenant-1")
	// No correlation id header

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	d.userRepo.AssertExpectations(t)
	// Assert idempotency repo was not called when no key is provided
	d.idem.AssertNotCalled(t, "IsProcessed", mock.Anything, mock.Anything, mock.Anything)
	d.idem.AssertNotCalled(t, "MarkAsProcessed", mock.Anything, mock.Anything, mock.Anything)
	d.bus.AssertExpectations(t)
	d.hasher.AssertExpectations(t)
}
