package http

// Registrar is implemented by an HTTP interface layer that maps concrete
// method/path pairs to fully-qualified strategy actions.
//
// Example:
//
//	routes.Handle("GET", "/messages/{id}", "app::messages:get")
type Registrar interface {
	Handle(method string, path string, action string) error
}
