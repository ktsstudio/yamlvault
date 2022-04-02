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

	if secret == nil {
		return nil, nil
	}

	var result map[string]interface{}
	if data, ok := secret.Data["data"]; ok {
		if data, ok := data.(map[string]interface{}); ok {
			result = data
		}
	} else {
		result = secret.Data
	}

	return result, nil
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
