// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestCasdoorContainerConnectivity tests basic connectivity to Casdoor running in a container.
// This test verifies that the Docker container setup works correctly for offline testing.
func TestCasdoorContainerConnectivity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping container test in short mode")
	}

	ctx := context.Background()

	// Start the Casdoor all-in-one container.
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

	// Always cleanup.
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	})

	// Get container endpoint.
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "8000")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	endpoint := fmt.Sprintf("http://%s:%s", host, port.Port())
	t.Logf("Casdoor endpoint: %s", endpoint)

	// Wait a bit for full initialization.
	time.Sleep(3 * time.Second)

	// Test 1: Verify port is accessible using docker CLI.
	t.Run("docker_port_check", func(t *testing.T) {
		inspect, err := container.Inspect(ctx)
		if err != nil {
			t.Fatalf("Failed to inspect container: %v", err)
		}
		containerID := inspect.ID
		t.Logf("Container ID: %s", containerID[:12])

		cmd := exec.Command("docker", "port", containerID[:12])
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("docker port command failed: %v, output: %s", err, string(output))
		}
		t.Logf("Docker port mapping:\n%s", string(output))

		if !strings.Contains(string(output), "8000/tcp") {
			t.Error("Expected port 8000/tcp to be mapped")
		}
	})

	// Test 2: HTTP health check.
	t.Run("http_health_check", func(t *testing.T) {
		resp, err := http.Get(endpoint + "/api/health")
		if err != nil {
			t.Fatalf("Health check request failed: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		body, _ := io.ReadAll(resp.Body)
		t.Logf("Health check response: status=%d, body=%s", resp.StatusCode, string(body))

		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	// Test 3: Fetch credentials via API (login with admin).
	t.Run("fetch_credentials_via_api", func(t *testing.T) {
		config, err := fetchConfigViaAPI(endpoint, defaultAdminPassword)
		if err != nil {
			t.Fatalf("Failed to fetch credentials: %v", err)
		}

		t.Logf("Fetched credentials:")
		t.Logf("  Endpoint: %s", config.Endpoint)
		t.Logf("  ClientID: %s", config.ClientID)
		if len(config.ClientSecret) > 8 {
			t.Logf("  ClientSecret: %s... (truncated)", config.ClientSecret[:8])
		} else {
			t.Logf("  ClientSecret: %s", config.ClientSecret)
		}
		t.Logf("  Organization: %s", config.OrganizationName)
		t.Logf("  Application: %s", config.ApplicationName)
		t.Logf("  Certificate: %d bytes", len(config.Certificate))
	})

	// Test 4: SDK Client CRUD operations.
	t.Run("sdk_client_crud", func(t *testing.T) {
		config, err := fetchConfigViaAPI(endpoint, defaultAdminPassword)
		if err != nil {
			t.Fatalf("Failed to fetch credentials: %v", err)
		}

		client := casdoorsdk.NewClient(
			config.Endpoint,
			config.ClientID,
			config.ClientSecret,
			config.Certificate,
			config.OrganizationName,
			config.ApplicationName,
		)

		// List organizations.
		orgs, err := client.GetOrganizations()
		if err != nil {
			t.Fatalf("GetOrganizations failed: %v", err)
		}
		t.Logf("Found %d organizations", len(orgs))
		for _, org := range orgs {
			t.Logf("  - %s/%s", org.Owner, org.Name)
		}

		// Create a test organization.
		rName := "tf-test-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
		newOrg := &casdoorsdk.Organization{
			Owner:        "admin",
			Name:         rName,
			DisplayName:  "Container Test Org",
			PasswordType: "bcrypt",
		}

		success, err := client.AddOrganization(newOrg)
		if err != nil {
			t.Fatalf("AddOrganization failed: %v", err)
		}
		if !success {
			t.Fatal("AddOrganization returned false")
		}
		t.Logf("Created organization: %s", rName)

		// Cleanup - delete the test organization.
		t.Cleanup(func() {
			_, err := client.DeleteOrganization(&casdoorsdk.Organization{Owner: "admin", Name: rName})
			if err != nil {
				t.Logf("Cleanup: failed to delete organization %s: %v", rName, err)
			} else {
				t.Logf("Cleanup: deleted organization %s", rName)
			}
		})

		// Read it back.
		org, err := client.GetOrganization(rName)
		if err != nil {
			t.Fatalf("GetOrganization failed: %v", err)
		}
		if org == nil {
			t.Fatal("GetOrganization returned nil")
		}
		t.Logf("Read organization: %s/%s (%s)", org.Owner, org.Name, org.DisplayName)

		// Update.
		org.DisplayName = "Updated Container Test Org"
		success, err = client.UpdateOrganization(org)
		if err != nil {
			t.Fatalf("UpdateOrganization failed: %v", err)
		}
		if !success {
			t.Fatal("UpdateOrganization returned false")
		}
		t.Logf("Updated organization: %s", rName)

		// Verify update.
		org, err = client.GetOrganization(rName)
		if err != nil {
			t.Fatalf("GetOrganization after update failed: %v", err)
		}
		if org.DisplayName != "Updated Container Test Org" {
			t.Errorf("Expected DisplayName 'Updated Container Test Org', got '%s'", org.DisplayName)
		}
	})
}
