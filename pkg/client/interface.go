package client

type Builder interface {
	New() (any, error)
}
