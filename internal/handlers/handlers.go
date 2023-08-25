package handlers

type APIHandler struct{}

func NewAPIHandler() (h *APIHandler, err error) {
	hh := &APIHandler{}

	return hh, err
}

type BaseBodyDater interface {
	IsValid() bool
}
