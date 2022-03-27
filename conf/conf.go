package conf

// IConfProvider
type IConfProvider interface {
	Add(string, interface{})
	Remove(string)
	Get(string)
}

// Conf
type Conf struct {
}
