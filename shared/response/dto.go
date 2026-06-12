package response

type Response[T any] struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}
