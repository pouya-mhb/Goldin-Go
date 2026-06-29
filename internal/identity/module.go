package identity

// Module contains the Identity bounded context dependencies.
//
// It is intentionally small for now. Application services, inbound handlers,
// and outbound adapters will be added here as the identity context grows.
type Module struct{}

// NewModule constructs the Identity module.
func NewModule() *Module {
	return &Module{}
}
