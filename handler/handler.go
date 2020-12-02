package handler

import "sync"

// Handler struct.
type Handler struct {
	sync.Mutex
	opts *Options
}

// New handler struct.
func New(opt ...Option) *Handler {
	opts := NewOptions(opt...)

	return &Handler{opts: &opts}
}

// Options struct
type Options struct {
	ConfigFile   string
	GodocURL     string
	MappingRules []Rules
}

// Option type
type Option func(*Options)

// NewOptions returns a new set of options for the http handler.
func NewOptions(options ...Option) Options {
	opts := Options{
		ConfigFile:   defaultConfigFile,
		MappingRules: make([]Rules, 0),
	}

	for _, o := range options {
		o(&opts)
	}

	return opts
}

// WithConfigFile sets the config file path.
func WithConfigFile(configFile string) Option {
	return func(o *Options) {
		o.ConfigFile = configFile
	}
}

// WithGodocURL sets the redirect URL for godoc.
func WithGodocURL(godocURL string) Option {
	return func(o *Options) {
		o.GodocURL = godocURL
	}
}

const (
	defaultConfigFile = "/tmp/config.json"
)
