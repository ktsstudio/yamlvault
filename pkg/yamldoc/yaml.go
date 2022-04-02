package yamldoc

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

type YamlDoc struct {
	data map[string]interface{}
}

type ValueTraverseFunc func(value interface{}) (interface{}, error)

func New(reader io.Reader) (*YamlDoc, error) {
	contents, err := ioutil.ReadAll(reader)
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

func NewFromFile(path string) (*YamlDoc, error) {
	if path == "-" {
		return New(os.Stdin)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return New(file)
}

func (d *YamlDoc) Save(writer io.Writer) error {
	out, err := yaml.Marshal(d.data)
	if err != nil {
		return err
	}

	if _, err := writer.Write(out); err != nil {
		return err
	}

	return nil
}

func (d *YamlDoc) SaveFile(path string) error {
	if path == "-" {
		return d.Save(os.Stdout)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return d.Save(file)
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
