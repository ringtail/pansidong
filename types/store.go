package types

type ListOptions struct {
	Limit int
}

type RefreshOptions struct {
	Force bool
}

type InsertOptions struct {
	Force bool
}

type Store interface {
	Next(options *ListOptions) ([]*ProxyIP, error)
	List(options *ListOptions) ([]*ProxyIP, error)
	Get(key string) (*ProxyIP, error)
	Expire(key string) error
}

type MemoryStore interface {
	Store
	Refresh(ips []*ProxyIP, options *RefreshOptions) error
}

type BackendStore interface {
	Store
	Insert(ips []*ProxyIP, options *InsertOptions) error
}


