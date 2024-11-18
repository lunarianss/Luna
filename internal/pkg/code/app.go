package code

const (
	// ErrAppMapMode - 500: Error occurred while attempt to index from appTemplate using mode.
	ErrAppMapMode int = iota + 110101
	// ErrAppNotFoundRelatedConfig - 500: Error occurred while attempt to find app related config.
	ErrAppNotFoundRelatedConfig
)
