package vaulter

import (
	"fmt"
	vault "github.com/hashicorp/vault/api"
)

type Vaulter struct {
	client *vault.Client
}

func New() (*Vaulter, error) {
	config := vault.DefaultConfig()
	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Vaulter{
		client: client,
	}, nil
}

func (v *Vaulter) RetrieveData(path string) (map[string]interface{}, error) {
	secret, err := v.client.Logical().Read(path)
	if err != nil {
		return nil, err
	}

	if secret.Data == nil || secret.Data["data"] == nil {
		return nil, nil
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data type assertion failed: %T %#v", secret.Data["data"], secret.Data["data"])
	}

	return data, nil
}

func (v *Vaulter) RetrieveStringKey(path, key string) (string, error) {
	data, err := v.RetrieveData(path)
	if err != nil {
		return "", err
	}

	value, ok := data[key].(string)
	if !ok {
		return "", fmt.Errorf("value type assertion failed: %T %#v", data[key], data[key])
	}

	return value, nil
}
