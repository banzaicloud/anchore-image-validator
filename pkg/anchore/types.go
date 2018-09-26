package anchore

// Check type for Anchore check
type Check struct {
	LastEvaluation string `json:"last_evaluation"`
	PolicyId       string `json:"policy_id"`
	Status         string `json:"status"`
}

// Image type for Anchore image
type Image struct {
	ImageDigest string `json:"imageDigest"`
}

// SHAResult type for Anchore result
type SHAResult struct {
	Status string
}
