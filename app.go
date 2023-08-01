package kart

type AppInfo interface {
	ID() string
	Name() string
	Version() string
}
