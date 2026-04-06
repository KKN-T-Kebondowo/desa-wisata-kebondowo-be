package e2e

import (
	"kebondowo/models"
	"testing"
)

// --- Health / Auth Flow ---

func TestHealthCheck(t *testing.T) {
	// The test server does not have a /healthchecker route (not added in setupServer).
	// This tests that an unknown route returns 404.
	resp := doRequest("GET", "/api/healthchecker", nil, "")
	if resp.StatusCode != 404 {
		t.Errorf("expected 404 for unregistered healthcheck route, got %d", resp.StatusCode)
	}
}

func TestAuthLogin(t *testing.T) {
	body := map[string]string{
		"username": "admin",
		"password": "admin1234",
	}
	resp := doRequest("POST", "/api/auth/login", body, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}

	result := parseJSON(resp.Body)
	if _, ok := result["access_token"]; !ok {
		t.Error("response missing access_token")
	}
	if _, ok := result["refresh_token"]; !ok {
		t.Error("response missing refresh_token")
	}
}

func TestAuthLoginWrongPassword(t *testing.T) {
	body := map[string]string{
		"username": "admin",
		"password": "wrongpassword",
	}
	resp := doRequest("POST", "/api/auth/login", body, "")
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for wrong password, got %d", resp.StatusCode)
	}
}

func TestAuthLoginNonExistentUser(t *testing.T) {
	body := map[string]string{
		"username": "nonexistent",
		"password": "anything",
	}
	resp := doRequest("POST", "/api/auth/login", body, "")
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for non-existent user, got %d", resp.StatusCode)
	}
}

func TestAuthLoginMissingFields(t *testing.T) {
	body := map[string]string{
		"username": "admin",
	}
	resp := doRequest("POST", "/api/auth/login", body, "")
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for missing password, got %d", resp.StatusCode)
	}
}

func TestRegisterRequiresAuth(t *testing.T) {
	body := map[string]interface{}{
		"username": "newuser",
		"password": "password123",
		"roleid":   2,
	}
	resp := doRequest("POST", "/api/auth/register", body, "")
	if resp.StatusCode != 401 {
		t.Errorf("expected 401 for unauthenticated register, got %d: %s", resp.StatusCode, resp.Body)
	}
}

func TestRegisterWithAuth(t *testing.T) {
	// Get user role ID
	var userRole models.Role
	testDB.Where("name = ?", "user").First(&userRole)

	body := map[string]interface{}{
		"username": "testuser_register",
		"password": "password123",
		"roleid":   userRole.ID,
	}
	resp := doRequest("POST", "/api/auth/register", body, authToken)
	if resp.StatusCode != 201 {
		t.Fatalf("expected 201, got %d: %s", resp.StatusCode, resp.Body)
	}
	result := parseJSON(resp.Body)
	user := result["user"].(map[string]interface{})
	if user["username"] != "testuser_register" {
		t.Errorf("expected username 'testuser_register', got %v", user["username"])
	}

	// Cleanup
	testDB.Exec("DELETE FROM users WHERE username = 'testuser_register'")
}

func TestRefreshToken(t *testing.T) {
	// First login to get a refresh token
	loginBody := map[string]string{
		"username": "admin",
		"password": "admin1234",
	}
	loginResp := doRequest("POST", "/api/auth/login", loginBody, "")
	loginResult := parseJSON(loginResp.Body)
	refreshToken := loginResult["refresh_token"].(string)

	// Use refresh token
	body := map[string]string{
		"refresh_token": refreshToken,
	}
	resp := doRequest("POST", "/api/auth/refresh", body, "")
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}
	result := parseJSON(resp.Body)
	if _, ok := result["access_token"]; !ok {
		t.Error("response missing access_token after refresh")
	}
}

func TestRefreshTokenInvalid(t *testing.T) {
	body := map[string]string{
		"refresh_token": "invalid-token",
	}
	resp := doRequest("POST", "/api/auth/refresh", body, "")
	if resp.StatusCode != 403 {
		t.Errorf("expected 403 for invalid refresh token, got %d", resp.StatusCode)
	}
}

func TestGetMe(t *testing.T) {
	resp := doRequest("GET", "/api/users/me", nil, authToken)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}
	result := parseJSON(resp.Body)
	user := result["user"].(map[string]interface{})
	if user["username"] != "admin" {
		t.Errorf("expected username 'admin', got %v", user["username"])
	}
}

func TestGetMeNoAuth(t *testing.T) {
	resp := doRequest("GET", "/api/users/me", nil, "")
	if resp.StatusCode != 401 {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestLogout(t *testing.T) {
	resp := doRequest("GET", "/api/auth/logout", nil, authToken)
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d: %s", resp.StatusCode, resp.Body)
	}
	result := parseJSON(resp.Body)
	if result["status"] != "success" {
		t.Errorf("expected status 'success', got %v", result["status"])
	}
}
