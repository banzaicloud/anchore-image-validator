package anchore

type Check struct {
	LastEvaluation string `json:"last_evaluation"`
	PolicyId       string `json:"policy_id"`
	Status         string `json:"status"`
}

type Image struct {
	ImageDigest string `json:"imageDigest"`
}

type SHAResult struct {
	Status string
}
