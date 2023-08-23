package component

// Component is the interface that represent a component.
type Component interface {
	Init()
	AfterInit()
	BeforeShutdown()
	Shutdown()
}
