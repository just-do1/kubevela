package driver

import (
	"errors"

	"github.com/oam-dev/kubevela/references/appfile/api"
)

// ConfigMapDriverName is local storage driver name
const ConfigMapDriverName = "ConfigMap"

// ConfigMap Storage
type ConfigMap struct {
	api.Driver
}

// NewConfigMapStorage get storage client of ConfigMap type
func NewConfigMapStorage() *ConfigMap {
	return &ConfigMap{}
}

// Name of local storage
func (c *ConfigMap) Name() string {
	return ConfigMapDriverName
}

// List applications from configmap storage
func (c *ConfigMap) List(envName string) ([]*api.Application, error) {
	// TODO support configmap storage
	return nil, errors.New("not implement")
}

// Save applications from configmap storage
func (c *ConfigMap) Save(app *api.Application, envName string) error {
	// TODO support configmap storage
	return errors.New("not implement")
}

// Delete applications from configmap storage
func (c *ConfigMap) Delete(envName, appName string) error {
	// TODO support configmap storage
	return errors.New("not implement")
}

// Get applications from configmap storage
func (c *ConfigMap) Get(envName, appName string) (*api.Application, error) {
	// TODO support configmap storage
	return nil, errors.New("not implement")
}
