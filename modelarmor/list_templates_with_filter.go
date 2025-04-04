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

// Sample code for listing model armor templates with filters.

package modelarmor

import (
	"context"
	"fmt"
	"io"
	"strings"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// listModelArmorTemplatesWithFilter lists all model armor templates in the specified project and location with a filter.
func listModelArmorTemplatesWithFilter(w io.Writer, projectID, locationID, templateID string) ([]string, error) {
	// [START modelarmor_list_templates_with_filter]
	ctx := context.Background()

	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx,
		option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationID)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// Preparing the parent path
	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)

	// Get the list of templates
	req := &modelarmorpb.ListTemplatesRequest{
		Parent: parent,
		Filter: fmt.Sprintf(`name="%s/templates/%s"`, parent, templateID),
	}

	it := client.ListTemplates(ctx, req)
	var templateNames []string

	for {
		template, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate templates: %v", err)
		}
		templateNames = append(templateNames, template.Name)
	}

	// Print templates name using fmt.Fprintf with the io.Writer
	fmt.Fprintf(w, "Templates Found: %s\n", strings.Join(templateNames, ", "))
	// [END modelarmor_list_templates_with_filter]

	return templateNames, nil
}
