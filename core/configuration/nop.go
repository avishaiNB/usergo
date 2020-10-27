package configuration

type nopConfigurationClient struct{}

// NewNopConfigurationClient returns a log.Logger that doesn't do anything.
// Should be used `for testing only
func NewNopConfigurationClient() Client {
	return nopConfigurationClient{}
}

func (nopConfigurationClient nopConfigurationClient) Get(key string, defaultValue interface{}) interface{} {
	return nil
}
