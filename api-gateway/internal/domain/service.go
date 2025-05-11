package domain

type ProxyService interface {
	Proxy(path string) (string, error)
}
