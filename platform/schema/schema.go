package schema

import (
	"io"
	"os"
	"path/filepath"

	"github.com/EstebanGitPro/motogo-backend/tools/utils"
	"github.com/kaptinlin/jsonschema"
)

type Validators struct {
	FileReader        FileReaderInterface
	RegisterValidator *jsonschema.Schema
}

type FileReaderInterface interface {
	ReadJsonSchema(resourcePath string) ([]byte, error)
}

type DefaultFileReader struct{}

func (f *DefaultFileReader) ReadJsonSchema(resourcePath string) ([]byte, error) {

	root, err := utils.FindModuleRoot()

	if err != nil {
		return nil, err
	}

	data, err := os.Open(filepath.Join(root, "platform/schema/json_schemas", resourcePath))
	if err != nil {
		return nil, err
	}
	defer data.Close()

	return io.ReadAll(data)

}

func NewValidator(fileReader FileReaderInterface) (*Validators, error) {
	validator := &Validators{
		FileReader: fileReader,
	}

	register, err := validator.createSchema("register_person_schema.json")
	if err != nil {
		return nil, err
	}

	validator.RegisterValidator = register

	return validator, nil

}

func (v *Validators) createSchema(resourcePath string) (*jsonschema.Schema, error) {
	compiler := jsonschema.NewCompiler()
	compiler.AssertFormat = true
	schemaJSON, err := v.FileReader.ReadJsonSchema(resourcePath)
	if err != nil {
		return nil, NewSchemaError("MOD_V_VAL_ERR_00003", "Error leyendo esquema JSON: "+err.Error())
	}

	if schemaJSON == nil {
		return nil, NewSchemaError("MOD_V_VAL_ERR_00004", "El esquema JSON está vacío o es nulo")
	}

	schema, err := compiler.Compile(schemaJSON)
	if err != nil {
		return nil, NewSchemaError("MOD_V_VAL_ERR_00005", "Error compilando esquema: "+err.Error())
	}

	return schema, nil
}

// ValidateRegister validates data against the register person schema
func (v *Validators) ValidateRegister(data interface{}) error {
	if v.RegisterValidator == nil {
		return NewSchemaError("MOD_V_VAL_ERR_00004", "Validador de registro no inicializado")
	}

	result := v.RegisterValidator.Validate(data)
	if !result.IsValid() {
		// Collect all validation errors
		var errorMessages []string
		for _, err := range result.Errors {
			errorMessages = append(errorMessages, err.Message)
		}
		
		if len(errorMessages) > 0 {
			return NewSchemaError("MOD_V_VAL_ERR_00006", "Errores de validación: "+errorMessages[0])
		}
	}

	return nil
}
