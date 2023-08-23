package component

type (
	options struct {
		name      string              // component name
		nameFunc  func(string) string // rename handler name
		schedName string              // schedName name
	}

	// Option used to customize handler
	Option func(options *options)
)

// WithName used to rename component name
func WithName(name string) Option {
	return func(opt *options) {
		opt.name = name
	}
}

// WithNameFunc override handler name by specific function
// such as: strings.ToUpper/strings.ToLower
func WithNameFunc(fn func(string) string) Option {
	return func(opt *options) {
		opt.nameFunc = fn
	}
}

// WithSchedulerName set the name of the service scheduler
func WithSchedulerName(name string) Option {
	return func(opt *options) {
		opt.schedName = name
	}
}
