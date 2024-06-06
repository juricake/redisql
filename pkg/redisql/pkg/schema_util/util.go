package schema_util

import (
	"encoding/json"
	"github.com/alecthomas/jsonschema"
	"github.com/xeipuuv/gojsonschema"
)

func Encode(definition interface{}) ([]byte, error) {
	schemaJSON, err := json.MarshalIndent(jsonschema.Reflect(definition), "", "  ")
	if err != nil {
		return nil, err
	}

	return schemaJSON, nil
}

func Validate(data interface{}, schema string) error {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	valueLoader := gojsonschema.NewGoLoader(data)
	_, err := gojsonschema.Validate(schemaLoader, valueLoader)
	if err != nil {
		return err
	}

	return nil
}
