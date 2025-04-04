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

// Sample code for sanitizing a model response using the model armor.

package main

import (
	"context"
	"fmt"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/option"
)

// sanitizeModelResponse sanitizes a model response using the Model Armor API.
func sanitizeModelResponse(projectID, locationID, templateID, modelResponse string) (*modelarmorpb.SanitizeModelResponseResponse, error) {
	// [START modelarmor_sanitize_model_response]
	ctx := context.Background()

	// TODO(Developer): Uncomment and set these variables.
	// projectID := "YOUR_PROJECT_ID"
	// locationID := "us-central1"
	// templateID := "template_id"
	// modelResponse := "The model response data to sanitize"

	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx,
		option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationID)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// Initialize request argument(s)
	modelResponseData := &modelarmorpb.DataItem{
		DataItem: &modelarmorpb.DataItem_Text{
			Text: modelResponse,
		},
	}

	// Prepare request for sanitizing model response.
	req := &modelarmorpb.SanitizeModelResponseRequest{
		Name:              fmt.Sprintf("projects/%s/locations/%s/templates/%s", projectID, locationID, templateID),
		ModelResponseData: modelResponseData,
	}

	// Sanitize the model response.
	response, err := client.SanitizeModelResponse(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to sanitize model response: %v", err)
	}

	// Sanitization Result.
	fmt.Printf("Sanitization Result: %v\n", response)

	// [END modelarmor_sanitize_model_response]

	return response, nil
}
