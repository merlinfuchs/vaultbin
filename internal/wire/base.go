package wire

type APIResponse[T any] struct {
	Success bool `json:"success" tstype:"true"`
	Data    T    `json:"data,omitempty"`
}

type APIErrorCode string
