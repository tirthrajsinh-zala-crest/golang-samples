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

package parametermanager

// [START parametermanager_create_param]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
)

// createParam creates a new parameter with the format type "unformatted" in Parameter Manager.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// parameterID: The ID of the parameter to be created.
//
// The function returns an error if the parameter creation fails.
func createParam(w io.Writer, projectID, parameterID string) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Parameter Manager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the create parameter.
	parent := fmt.Sprintf("projects/%s/locations/global", projectID)

	// Build the request to create a new parameter
	req := &parametermanagerpb.CreateParameterRequest{
		Parent:      parent,
		ParameterId: parameterID,
		Parameter:   &parametermanagerpb.Parameter{},
	}

	// Call the API to create the parameter.
	parameter, err := client.CreateParameter(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create parameter: %w", err)
	}

	fmt.Fprintf(w, "Created parameter: %s\n", parameter.Name)
	return nil
}

// [END parametermanager_create_param]
