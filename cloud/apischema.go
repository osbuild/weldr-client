package cloud

// APIResponse is returned by the cloudapi when there is an error
type APIResponse struct {
	Kind    string `json:"kind"`
	ID      string `json:"id"`
	Code    string `json:"code"`
	Details string `json:"details"`
	Reason  string `json:"reason"`
}
