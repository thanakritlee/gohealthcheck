package config

// Config can be implements by objects that can import
// some config in key-value format and able to get the
// config value given it's key.
type Config interface {
	GetConfig(string) string
}
