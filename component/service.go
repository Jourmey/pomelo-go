package component

type (
	//Handler represents a message.Message's handler's meta information.
	Handler struct {
	}

	// Service implements a specific service, some of it's methods will be
	// called when the correspond events is occurred.
	Service struct {
		Name    string // name of service
		Comp    Component
		Options options // options
	}
)

func NewService(comp Component, opts []Option) *Service {
	s := &Service{
		Comp: comp,
	}

	// apply options
	for i := range opts {
		opt := opts[i]
		opt(&s.Options)
	}

	s.Name = s.Options.name

	return s
}
