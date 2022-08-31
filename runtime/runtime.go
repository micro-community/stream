package runtime

// Service for monibuca
type Service interface {
	Init() error
	Run(configFile string) error
	String() string
}
