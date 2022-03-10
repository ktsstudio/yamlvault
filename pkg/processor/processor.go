package processor

import (
	"fmt"
	"github.com/ktsstudio/yamlvault/pkg/vaulter"
	"github.com/ktsstudio/yamlvault/pkg/yamldoc"
	"log"
	"strings"
)

type Processor struct {
	Vaulter   *vaulter.Vaulter
	pathCache map[string]map[string]interface{}
}

func New(v *vaulter.Vaulter) (*Processor, error) {
	return &Processor{
		Vaulter:   v,
		pathCache: make(map[string]map[string]interface{}),
	}, nil
}

func (p *Processor) Process(doc *yamldoc.YamlDoc) error {
	err := doc.TraverseValues(func(value interface{}) (interface{}, error) {
		if value, ok := value.(string); ok {
			if !strings.HasPrefix(value, "vault:") {
				return value, nil
			}
			parts := strings.Split(value, ":")
			if len(parts) != 3 {
				return nil, fmt.Errorf("value must be of format vault:<oath>:<key>")
			}
			path, key := parts[1], parts[2]
			return p.getValueFromVault(path, key)
		}
		return value, nil
	})

	if err != nil {
		return err
	}

	log.Printf("done processing")

	return nil
}

func (p *Processor) getValueFromVault(path, key string) (string, error) {
	data, exists := p.pathCache[path]
	if !exists {
		var err error
		data, err = p.Vaulter.RetrieveData(path)
		if err != nil {
			return "", err
		}
		log.Printf("fetched `%s' from Vault\n", path)
		p.pathCache[path] = data
	}

	if data == nil {
		return "", fmt.Errorf("vault path `%s' not found", path)
	}

	value, exists := data[key]
	if !exists {
		return "", fmt.Errorf("vault path `%s' does not contain key `%s'", path, key)
	}

	stringValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("value type assertion failed: `%T' `%#v'", data[key], data[key])
	}

	return stringValue, nil
}
