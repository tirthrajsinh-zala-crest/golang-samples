// Copyright 2024 Google LLC
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

package snippets

// [START compute_consistency_group_add_disk]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

// addDiskConsistencyGroup adds a disk to a consistency group for a project in a given region.
func addDiskConsistencyGroup(w io.Writer, projectID, region, groupName, diskName string) error {
	// projectID := "your_project_id"
	// region := "europe-west4"
	// diskName := "your_disk_name"
	// groupName := "your_group_name"

	ctx := context.Background()
	disksClient, err := compute.NewRegionDisksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewResourcePoliciesRESTClient: %w", err)
	}
	defer disksClient.Close()

	consistencyGroupUrl := fmt.Sprintf("projects/%s/regions/%s/resourcePolicies/%s", projectID, region, groupName)

	req := &computepb.AddResourcePoliciesRegionDiskRequest{
		Project: projectID,
		Disk:    diskName,
		RegionDisksAddResourcePoliciesRequestResource: &computepb.RegionDisksAddResourcePoliciesRequest{
			ResourcePolicies: []string{consistencyGroupUrl},
		},
		Region: region,
	}

	op, err := disksClient.AddResourcePolicies(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to add disk: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Disk added\n")

	return nil
}

// [END compute_consistency_group_add_disk]
