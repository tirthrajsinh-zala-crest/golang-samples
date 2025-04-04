// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     [https://www.apache.org/licenses/LICENSE-2.0](https://www.apache.org/licenses/LICENSE-2.0)
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Sample code for sanitizing user prompt with model armor.

package main

import (
	"context"
	"fmt"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/option"
)

// sanitizeUserPrompt sanitizes a user prompt using the Model Armor API.
func sanitizeUserPrompt(projectID, locationID, templateID, userPrompt string) (*modelarmorpb.SanitizeUserPromptResponse, error) {
	// [START modelarmor_sanitize_user_prompt]
	ctx := context.Background()

	// TODO(Developer): Uncomment and set these variables.
	// projectID := "YOUR_PROJECT_ID"
	// locationID := "us-central1"
	// templateID := "template_id"
	// userPrompt := "Prompt entered by the user"

	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx,
		option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationID)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// Initialize request argument(s)
	userPromptData := &modelarmorpb.DataItem{
		DataItem: &modelarmorpb.DataItem_Text{
			Text: userPrompt,
		},
	}

	// Prepare request for sanitizing the defined prompt.
	req := &modelarmorpb.SanitizeUserPromptRequest{
		Name:           fmt.Sprintf("projects/%s/locations/%s/templates/%s", projectID, locationID, templateID),
		UserPromptData: userPromptData,
	}

	// Sanitize the user prompt.
	response, err := client.SanitizeUserPrompt(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to sanitize user prompt: %v", err)
	}

	// Sanitization Result.
	fmt.Printf("Sanitization Result: %v\n", response)

	// [END modelarmor_sanitize_user_prompt]

	return response, nil
}
