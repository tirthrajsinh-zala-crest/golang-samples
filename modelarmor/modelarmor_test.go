// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package modelarmor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

// testOrganizationID returns the organization ID.
func testOrganizationID(t *testing.T) string {
	orgID := "951890214235"
	return orgID
}

// testFolderID returns the folder ID.
func testFolderID(t *testing.T) string {
	folderID := "695279264361"
	return folderID
}

// testLocation retrieves the GOLANG_SAMPLES_LOCATION environment variable
// used to determine the region for running the test.
// Skips the test if the environment variable is not set.
func testLocation(t *testing.T) string {
	t.Helper()

	v := os.Getenv("GOLANG_SAMPLES_LOCATION")
	if v == "" {
		t.Skip("testLocation: missing GOLANG_SAMPLES_LOCATION")
	}

	return v
}

// testClient initializes and returns a new Model Armor API client and context
// targeting the endpoint based on the specified location.
func testClient(t *testing.T) (*modelarmor.Client, context.Context) {
	t.Helper()

	ctx := context.Background()
	locationId := testLocation(t)
	// Create option for Model Armor client.
	opts := option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationId))
	// Create a new client using the regional endpoint
	client, err := modelarmor.NewClient(ctx, opts)
	if err != nil {
		t.Fatalf("testClient: failed to create client: %v", err)
	}

	return client, ctx
}

// testCleanupTemplate deletes the specified Model Armor template if it exists,
// ignoring the error if the template is already deleted.
func testCleanupTemplate(t *testing.T, templateName string) {
	t.Helper()

	client, ctx := testClient(t)
	if err := client.DeleteTemplate(ctx, &modelarmorpb.DeleteTemplateRequest{Name: templateName}); err != nil {
		// Ignore NotFound errors (template may already be deleted)
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			t.Fatalf("testCleanupTemplate: failed to delete template: %v", err)
		}
	}
}

func testDisableFloorSettings(floorSettingName string, locationID string) error {
	ctx := context.Background()

	client, err := modelarmor.NewClient(ctx,
		option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationID)),
	)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	disable := false

	floorSetting := &modelarmorpb.FloorSetting{
		Name:                          floorSettingName,
		EnableFloorSettingEnforcement: &disable,
	}

	req := &modelarmorpb.UpdateFloorSettingRequest{
		FloorSetting: floorSetting,
	}

	_, err = client.UpdateFloorSetting(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to disable floor setting: %w", err)
	}
	return nil
}

// TestCreateModelArmorTemplate verifies the creation of a Model Armor template.
// It ensures the output contains a confirmation message after creation.
func TestCreateModelArmorTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())
	templateName := fmt.Sprintf("projects/%s/locations/%s/templates/%s", tc.ProjectID, testLocation(t), templateID)
	var b bytes.Buffer
	if err := createModelArmorTemplate(&b, tc.ProjectID, testLocation(t), templateID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupTemplate(t, templateName)

	if got, want := b.String(), "Created template:"; !strings.Contains(got, want) {
		t.Errorf("createModelArmorTemplate: expected %q to contain %q", got, want)
	}
}

// TestDeleteModelArmorTemplate verifies the deletion of a Model Armor template.
// It ensures the output contains a confirmation message after deletion.
func TestDeleteModelArmorTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)

	templateID := fmt.Sprintf("test-model-armor-%s", uuid.New().String())

	var buf bytes.Buffer
	// Create template first to ensure it exists for deletion
	if err := createModelArmorTemplate(&buf, tc.ProjectID, testLocation(t), templateID); err != nil {
		t.Fatal(err)
	}

	// Attempt to delete the template
	if err := deleteModelArmorTemplate(&buf, tc.ProjectID, testLocation(t), templateID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Successfully deleted Model Armor template:"; !strings.Contains(got, want) {
		t.Errorf("deleteModelArmorTemplate: expected %q to contain %q", got, want)
	}
}

// TestUpdateFolderFloorSettings tests updating floor settings for a specific folder.
// It verifies that the output buffer contains a confirmation message indicating a successful update.
func TestUpdateFolderFloorSettings(t *testing.T) {
	folderID := testFolderID(t)
	locationID := testLocation(t)
	var buf bytes.Buffer
	if err := updateFolderFloorSettings(&buf, folderID, locationID); err != nil {
		t.Fatal(err)
	}
	// Prepare folder floor setting path/name
	floorSettingName := fmt.Sprintf("folders/%s/locations/global/floorSetting", folderID)
	
	defer testDisableFloorSettings(floorSettingName, locationID)

	if got, want := buf.String(), "Updated folder floor setting: "; !strings.Contains(got, want) {
		t.Errorf("updateFolderFloorSettings: expected %q to contain %q", got, want)
	}
}

// TestUpdateOrganizationFloorSettings tests updating floor settings for a specific organization.
// It ensures the output buffer includes a success message confirming the update.
func TestUpdateOrganizationFloorSettings(t *testing.T) {
	organizationID := testOrganizationID(t)
	locationID := testLocation(t)
	var buf bytes.Buffer
	if err := updateOrganizationFloorSettings(&buf, organizationID, locationID); err != nil {
		t.Fatal(err)
	}

	// Prepare organization floor setting path/name
	floorSettingsName := fmt.Sprintf("organizations/%s/locations/global/floorSetting", organizationID)

	defer testDisableFloorSettings(floorSettingsName, locationID)

	if got, want := buf.String(), "Updated org floor setting: "; !strings.Contains(got, want) {
		t.Errorf("updateOrganizationFloorSettings: expected %q to contain %q", got, want)
	}
}

// TestUpdateProjectFloorSettings tests updating floor settings for a specific project.
// It checks that the resulting output includes the expected confirmation message.
func TestUpdateProjectFloorSettings(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)
	var buf bytes.Buffer
	if err := updateProjectFloorSettings(&buf, tc.ProjectID, locationID); err != nil {
		t.Fatal(err)
	}

	// Prepare project floor settings path/name
	floorSettingsName := fmt.Sprintf("projects/%s/locations/global/floorSetting", tc.ProjectID)

	defer testDisableFloorSettings(floorSettingsName, locationID)

	if got, want := buf.String(), "Updated project floor setting: "; !strings.Contains(got, want) {
		t.Errorf("updateProjectFloorSettings: expected %q to contain %q", got, want)
	}
}
