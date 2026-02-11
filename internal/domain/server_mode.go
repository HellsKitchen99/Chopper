package domain

type ServerMode string

const (
	ReleaseMode ServerMode = "release"
	DebugMode   ServerMode = "debug"
	TestMode    ServerMode = "test"
)
