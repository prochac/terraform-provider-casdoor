// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Test credentials from the Casdoor SDK demo server.
// Source: https://github.com/casdoor/casdoor-go-sdk/blob/master/casdoorsdk/test_util.go
const (
	testCasdoorEndpoint     = "https://demo.casdoor.com"
	testClientID            = "294b09fbc17f95daf2fe"
	testClientSecret        = "dd8982f7046ccba1bbd7851d5c1ece4e52bf039d"
	testCasdoorOrganization = "casbin"
	testCasdoorApplication  = "app-vue-python-example"
	testJwtPublicKey        = `-----BEGIN CERTIFICATE-----
MIIE+TCCAuGgAwIBAgIDAeJAMA0GCSqGSIb3DQEBCwUAMDYxHTAbBgNVBAoTFENh
c2Rvb3IgT3JnYW5pemF0aW9uMRUwEwYDVQQDEwxDYXNkb29yIENlcnQwHhcNMjEx
MDE1MDgxMTUyWhcNNDExMDE1MDgxMTUyWjA2MR0wGwYDVQQKExRDYXNkb29yIE9y
Z2FuaXphdGlvbjEVMBMGA1UEAxMMQ2FzZG9vciBDZXJ0MIICIjANBgkqhkiG9w0B
AQEFAAOCAg8AMIICCgKCAgEAsInpb5E1/ym0f1RfSDSSE8IR7y+lw+RJjI74e5ej
rq4b8zMYk7HeHCyZr/hmNEwEVXnhXu1P0mBeQ5ypp/QGo8vgEmjAETNmzkI1NjOQ
CjCYwUrasO/f/MnA1bZL8MjIoxRFG9/ONrMjf1xv0QAGnWG/d3aIR3O1E7gn9P9O
rrGjgePl1d+wOiIojyHCvXjHCNrtMwJGJWZrFcxMPAqjPyivBPW4s1K1WilMDhR2
v8a+VVFX8A0O7ZVwvKcFVzFZi1YMB36WGXhvYFZ98AIFn8yCUdJe/awjHH3p5KoO
O6J07NKhiUvC6uf/lKR0aCiRmCvoIVRX6V9hpQ1qZKNPWZLzDqpb4Ns3lYbQl0Hp
VWEFcKl8v4h0kd7zq/6pQ9xq5Cq4VkPYkBGnvqMz7BmhJB7dqWPNgYS6HZlChGGk
Nuc6c6X+6VFjpxY6Qv5qmfMtMfaVz02fUk9pq4YNE8gwNJQ5F1OKKJBLQZ7BlpJW
8ABHv8+L38IW5ZmZPwLrxUdbB5hL7IVQIyJBxHDQb8f+t/BQwKX2RlJwpHbq6t9E
cAbH07w4n2X/uE7mVMHDxBM6F2JZcfL1YzXb6dU5yP6sxlEt/ezEqpqP3svSuB0M
S10CWLl2E3XLLxMV6cUo/0TD0iQT7ZC9u9qMI9NKfB2BL6u5wP7Gxzpvbq6AS8AG
g0MCAwEAAaMQMA4wDAYDVR0TAQH/BAIwADANBgkqhkiG9w0BAQsFAAOCAgEAMvXJ
gwTsXcZIykFhv2YI71TcsFPU+kA5XoALRPnnAlmMRpKymKzC1M4xwxGPw0/mxDE4
N/ImHPaIbJLd4n8fUj9hL1sDLYX0czWjKGe24SiKK+8Lh5mYrJdJPxZX/yN8s9bZ
k39N4K8Y77ZJd2P/o/AqXIGjoMQv0vqMI2xPNbHhpg7sKaQFYCW/K2i4DJE8S2v+
zHV7kLe4D8gQoRCbKAbVi+PrgqBd7kNSS+xJd6A0BhU1bN4W/qGZhHdB7xhPr5gE
e0pVgAuY/fGRhH8TBTvhhjsDJLJJBJLp8zFt9M5FEQVbg7fVNFXy6hVyoP/J5W5r
4tVNpHHvK6DnWyKV4r6a0IhJeNLPJ3TFthA9kQIqJn9xBdJBbT6BOqVEsL1YKLWv
MJSCjfM2V/sLMnR7X1i5a5gYaEnED3L8Sq5N9K3N/7RvwVjLs7l3L0OMJX+pMKp2
Fa2G6eH2IH/L8fU1s5wH4dI26kPJNqE42IJS6mHWV2aQplHvNXNLNXJcQh3/pqz9
KFBQY1dI0zKy3G0C/vPhq4/lX9GvAY+FNnT+p2qLGnwP9E4dH1N1VjpjEmYQ3aIL
5+v2oLkLF9e8l8a7w+VGZujMPEpHfQC8/ZM4xP/9vrNqFP5PjJs6EpMN3fh7XpDo
/xRfXkKF5tT5h/FpLI1O3P8A0fRs5JNx7DAbLdQ=
-----END CERTIFICATE-----`

	// Default admin password for casdoor-all-in-one container.
	defaultAdminPassword = "123"
)

// CasdoorTestConfig holds the configuration for connecting to a Casdoor server.
type CasdoorTestConfig struct {
	Endpoint         string
	ClientID         string
	ClientSecret     string
	Certificate      string
	OrganizationName string
	ApplicationName  string
}

// TestEnv holds the test environment including optional container.
type TestEnv struct {
	Config    CasdoorTestConfig
	Container testcontainers.Container
}

// useLocalContainer returns true if tests should use a local Docker container
// instead of the demo server. Set CASDOOR_TEST_LOCAL=1 to enable.
func useLocalContainer() bool {
	return os.Getenv("CASDOOR_TEST_LOCAL") == "1"
}

// setupTestEnv sets up the test environment. If CASDOOR_TEST_LOCAL=1, it starts
// a local Casdoor container. Otherwise, it uses the demo server.
func setupTestEnv(ctx context.Context, t *testing.T) *TestEnv {
	t.Helper()

	if useLocalContainer() {
		return setupLocalContainer(ctx, t)
	}

	return &TestEnv{
		Config: getDemoConfig(),
	}
}

// cleanupTestEnv cleans up the test environment. Call this in t.Cleanup().
func cleanupTestEnv(ctx context.Context, t *testing.T, env *TestEnv) {
	t.Helper()

	if env.Container != nil {
		if err := env.Container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}
}

// setupLocalContainer starts a local Casdoor container and returns the test environment.
func setupLocalContainer(ctx context.Context, t *testing.T) *TestEnv {
	t.Helper()

	req := testcontainers.ContainerRequest{
		Image:        "casbin/casdoor-all-in-one:latest",
		ExposedPorts: []string{"8000/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForHTTP("/api/health").WithPort("8000/tcp").WithStatusCodeMatcher(func(status int) bool {
				return status == 200
			}),
			wait.ForLog("http server Running on"),
		).WithDeadline(180 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start casdoor container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "8000")
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("Failed to get container port: %v", err)
	}

	endpoint := fmt.Sprintf("http://%s:%s", host, port.Port())

	// Wait for Casdoor to be fully ready.
	time.Sleep(2 * time.Second)

	config, err := fetchConfigViaAPI(endpoint, defaultAdminPassword)
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("Failed to fetch config via API: %v", err)
	}

	t.Logf("Using local container at %s", endpoint)

	return &TestEnv{
		Config:    *config,
		Container: container,
	}
}

// fetchConfigViaAPI logs in with admin credentials and fetches application config via API.
func fetchConfigViaAPI(endpoint, adminPassword string) (*CasdoorTestConfig, error) {
	// Create HTTP client with cookie jar for session management.
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	// Step 1: Login with admin credentials.
	loginPayload := map[string]string{
		"application":  "app-built-in",
		"organization": "built-in",
		"username":     "admin",
		"password":     adminPassword,
		"type":         "login",
	}
	loginBody, _ := json.Marshal(loginPayload)

	loginResp, err := client.Post(
		endpoint+"/api/login",
		"application/json",
		strings.NewReader(string(loginBody)),
	)
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	defer func() { _ = loginResp.Body.Close() }()

	if loginResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed with status %d", loginResp.StatusCode)
	}

	var loginResult struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	}
	if err := json.NewDecoder(loginResp.Body).Decode(&loginResult); err != nil {
		return nil, fmt.Errorf("failed to decode login response: %w", err)
	}
	if loginResult.Status != "ok" {
		return nil, fmt.Errorf("login failed: %s", loginResult.Msg)
	}

	// Step 2: Get application details.
	appResp, err := client.Get(endpoint + "/api/get-application?id=admin/app-built-in")
	if err != nil {
		return nil, fmt.Errorf("get application request failed: %w", err)
	}
	defer func() { _ = appResp.Body.Close() }()

	var appResult struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
		Data   struct {
			ClientID     string `json:"clientId"`
			ClientSecret string `json:"clientSecret"`
			Cert         string `json:"cert"`
		} `json:"data"`
	}
	if err := json.NewDecoder(appResp.Body).Decode(&appResult); err != nil {
		return nil, fmt.Errorf("failed to decode application response: %w", err)
	}

	// Step 3: Get certificate.
	certName := appResult.Data.Cert
	if certName == "" {
		certName = "cert-built-in"
	}

	certResp, err := client.Get(endpoint + "/api/get-cert?id=admin/" + certName)
	if err != nil {
		return nil, fmt.Errorf("get cert request failed: %w", err)
	}
	defer func() { _ = certResp.Body.Close() }()

	var certResult struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
		Data   struct {
			Certificate string `json:"certificate"`
		} `json:"data"`
	}
	if err := json.NewDecoder(certResp.Body).Decode(&certResult); err != nil {
		return nil, fmt.Errorf("failed to decode cert response: %w", err)
	}

	certificate := certResult.Data.Certificate
	if certificate == "" {
		// Use a placeholder if certificate is not available.
		certificate = "placeholder-cert"
	}

	return &CasdoorTestConfig{
		Endpoint:         endpoint,
		ClientID:         appResult.Data.ClientID,
		ClientSecret:     appResult.Data.ClientSecret,
		Certificate:      certificate,
		OrganizationName: "built-in",
		ApplicationName:  "app-built-in",
	}, nil
}

// getDemoConfig returns the configuration for the Casdoor demo server.
func getDemoConfig() CasdoorTestConfig {
	return CasdoorTestConfig{
		Endpoint:         testCasdoorEndpoint,
		ClientID:         testClientID,
		ClientSecret:     testClientSecret,
		Certificate:      testJwtPublicKey,
		OrganizationName: testCasdoorOrganization,
		ApplicationName:  testCasdoorApplication,
	}
}

// setupTestConfig returns test configuration (for backwards compatibility).
func setupTestConfig(t *testing.T) CasdoorTestConfig {
	t.Parallel()
	t.Helper()
	ctx := context.Background()
	env := setupTestEnv(ctx, t)
	t.Cleanup(func() {
		cleanupTestEnv(ctx, t, env)
	})
	return env.Config
}

func testAccProtoV6ProviderFactories(_ CasdoorTestConfig) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"casdoor": providerserver.NewProtocol6WithError(New("test")()),
	}
}

func testAccProviderConfig(config CasdoorTestConfig) string {
	return fmt.Sprintf(`
provider "casdoor" {
  endpoint          = %q
  client_id         = %q
  client_secret     = %q
  certificate       = %q
  organization_name = %q
  application_name  = %q
}
`, config.Endpoint, config.ClientID, config.ClientSecret, config.Certificate, config.OrganizationName, config.ApplicationName)
}
