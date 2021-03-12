// Copyright © 2019 Banzai Cloud.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package anchore

// Check type for Anchore check.
type Check struct {
	LastEvaluation string `json:"last_evaluation"`
	PolicyID       string `json:"policy_id"`
	Status         string `json:"status"`
}

// Image type for Anchore image.
type Image struct {
	ImageDigest string `json:"imageDigest"`
	LastUpdated string `json:"last_updated"`
}

// SHAResult type for Anchore result.
type SHAResult struct {
	Status string
}
