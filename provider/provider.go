package provider

// provider
type provider struct {
	Name      string
	providers map[string]interface{}
}

// defaultProvider
var defaultProvider = &provider{
	Name:      "default",
	providers: map[string]interface{}{},
}

func GetProvider(key string) (interface{}, bool) {
	return defaultProvider.GetProvicer(key)
}
func AddProvider(key string, value interface{}) {
	defaultProvider.AddProvider(key, value)
}

func RemoveProvider(key string) {
	defaultProvider.RemoveProvider(key)
}

func NewProvider(name string) *provider {
	return &provider{Name: name}
}

func (provider *provider) AddProvider(key string, value interface{}) {
	provider.providers[key] = value
}

func (provider *provider) RemoveProvider(key string) {
	delete(provider.providers, key)
}

func (provider *provider) GetProvicer(key string) (interface{}, bool) {
	v, ok := provider.providers[key]
	return v, ok
}
