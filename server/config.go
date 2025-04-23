package server

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Configuration is config storage of application.
type Configuration struct {
	store *viper.Viper
}

func newConfiguration() *Configuration {
	return &Configuration{
		store: viper.New(),
	}
}

// IsSet checks to see if the key has been set in any of the data locations.
func (c *Configuration) IsSet(key string) bool {
	return c.store.IsSet(key)
}

// Get can retrieve any value given the key to use.
func (c *Configuration) Get(key string) interface{} {
	return c.store.Get(key)
}

// GetString returns the value associated with the key as a string.
func (c *Configuration) GetString(key string) string {
	return c.store.GetString(key)
}

// GetBool returns the value associated with the key as a boolean.
func (c *Configuration) GetBool(key string) bool {
	return c.store.GetBool(key)
}

// GetInt returns the value associated with the key as an integer.
func (c *Configuration) GetInt(key string) int {
	return c.store.GetInt(key)
}

// GetInt32 returns the value associated with the key as an integer.
func (c *Configuration) GetInt32(key string) int32 {
	return c.store.GetInt32(key)
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Configuration) GetInt64(key string) int64 {
	return c.store.GetInt64(key)
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Configuration) GetFloat64(key string) float64 {
	return c.store.GetFloat64(key)
}

// GetTime returns the value associated with the key as time.
func (c *Configuration) GetTime(key string) time.Time {
	return c.store.GetTime(key)
}

// GetDuration returns the value associated with the key as a duration.
func (c *Configuration) GetDuration(key string) time.Duration {
	return c.store.GetDuration(key)
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Configuration) GetStringSlice(key string) []string {
	return c.store.GetStringSlice(key)
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Configuration) GetStringMap(key string) map[string]interface{} {
	return c.store.GetStringMap(key)
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Configuration) GetStringMapString(key string) map[string]string {
	return c.store.GetStringMapString(key)
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *Configuration) GetStringMapStringSlice(key string) map[string][]string {
	return c.store.GetStringMapStringSlice(key)
}

// GetSizeInBytes returns the size of the value associated with the given key
// in bytes.
func (c *Configuration) GetSizeInBytes(key string) uint {
	return c.store.GetSizeInBytes(key)
}

// Unmarshal maps current store into a struct.
func (c *Configuration) Unmarshal(obj interface{}) error {
	return c.store.Unmarshal(obj)
}

// UnmarshalKey maps current store into a struct.
func (c *Configuration) UnmarshalKey(key string, obj interface{}) error {
	if c.store.Get(key) == nil {
		return fmt.Errorf("config key: [%v] not exists", key)
	}
	return c.store.UnmarshalKey(key, obj)
}

func (c *Configuration) apply(path string) {
	c.store.SetConfigFile(path)
	c.store.MergeInConfig()
	c.store.SetConfigFile("")
}
