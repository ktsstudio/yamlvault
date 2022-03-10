package yamldoc

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type YamlDoc struct {
	data map[string]interface{}
}

type ValueTraverseFunc func(value interface{}) (interface{}, error)

func New(path string) (*YamlDoc, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	if err := yaml.Unmarshal(contents, &data); err != nil {
		return nil, err
	}

	return &YamlDoc{
		data: data,
	}, nil
}

func (d *YamlDoc) Save(path string) error {
	out, err := yaml.Marshal(d.data)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, out, 0644); err != nil {
		return err
	}

	return nil
}

func (d *YamlDoc) TraverseValues(f ValueTraverseFunc) error {
	if len(d.data) == 0 {
		return nil
	}
	newData, err := d.dfs(d.data, f)
	if err != nil {
		return err
	}
	d.data = newData.(map[string]interface{})
	return nil
}

func (d *YamlDoc) traverseMap(input map[string]interface{}, f ValueTraverseFunc) (map[string]interface{}, error) {
	for k, v := range input {
		newVal, err := d.dfs(v, f)
		if err != nil {
			return nil, err
		}
		input[k] = newVal
	}
	return input, nil
}

func (d *YamlDoc) dfs(input interface{}, f ValueTraverseFunc) (interface{}, error) {
	switch input := input.(type) {
	case map[string]interface{}:
		return d.traverseMap(input, f)
	case map[interface{}]interface{}:
		stringMap := make(map[string]interface{})
		for k, v := range input {
			stringMap[fmt.Sprintf("%v", k)] = v
		}
		return d.traverseMap(stringMap, f)

	case []interface{}:
		for idx, v := range input {
			newVal, err := d.dfs(v, f)
			if err != nil {
				return nil, err
			}
			input[idx] = newVal
		}
		return input, nil

	default:
		return f(input)
	}
}
