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

// Sample code for updating the model armor folder settings of a folder.

package modelarmor

import (
	"context"
	"fmt"
	"io"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/option"
)

// updateFolderFloorSettings updates floor settings of a folder.
func updateFolderFloorSettings(w io.Writer, folderID, locationID string) (*modelarmorpb.FloorSetting, error) {
	// [START modelarmor_update_folder_floor_settings]
	ctx := context.Background()

	// TODO(Developer): Uncomment and set these variables.
	// folderID := "YOUR_FOLDER_ID"
	// locationID := "us-central1"

	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx,
		option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", locationID)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// Prepare folder floor settings path/name
	floorSettingsName := fmt.Sprintf("folders/%s/locations/global/floorSetting", folderID)

	// Prepare the floor setting update
	enableEnforcement := true
	floorSetting := &modelarmorpb.FloorSetting{
		Name: floorSettingsName,
		FilterConfig: &modelarmorpb.FilterConfig{
			RaiSettings: &modelarmorpb.RaiFilterSettings{
				RaiFilters: []*modelarmorpb.RaiFilterSettings_RaiFilter{
					{
						FilterType:      modelarmorpb.RaiFilterType_HATE_SPEECH,
						ConfidenceLevel: modelarmorpb.DetectionConfidenceLevel_HIGH,
					},
				},
			},
		},
		EnableFloorSettingEnforcement: &enableEnforcement,
	}

	// Prepare request for updating the floor setting.
	req := &modelarmorpb.UpdateFloorSettingRequest{
		FloorSetting: floorSetting,
	}

	// Update the floor setting.
	response, err := client.UpdateFloorSetting(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update floor setting: %v", err)
	}

	// Print the updated config
	fmt.Fprintf(w, "Updated Floor Setting: %v\n", response)

	// [END modelarmor_update_folder_floor_settings]

	return response, nil
}
