package idn

// Provider is an interface that defines the methods that a identity provider must implement.
type Provider interface {
	ID() string
}

type Hasher interface {
	Hash(string) string
}
