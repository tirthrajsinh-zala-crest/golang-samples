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

// Sample code for getting a model armor template.

package main

import (
	"context"
	"fmt"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/option"
)

// getModelArmorTemplate gets a Model Armor template.
func getModelArmorTemplate(projectID, location, templateID string) (*modelarmorpb.Template, error) {
	// [START modelarmor_get_template]
	ctx := context.Background()

	// TODO(Developer): Uncomment and set these variables.
	// projectID := "YOUR_PROJECT_ID"
	// location := "us-central1"
	// templateID := "template_id"

	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx,
		option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", location)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// Initialize request arguments.
	req := &modelarmorpb.GetTemplateRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/templates/%s", projectID, location, templateID),
	}

	// Get the template.
	response, err := client.GetTemplate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %v", err)
	}

	// Print the template name.
	fmt.Printf("Retrieved template: %s\n", response.Name)

	// [END modelarmor_get_template]

	return response, nil
}
