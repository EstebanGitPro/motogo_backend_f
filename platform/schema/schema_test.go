package schema_test

import (
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/platform/schema"
	"github.com/stretchr/testify/assert"
)

// MockFileReader implements schema.FileReader for testing.
type MockFileReader struct {
	schemas map[string][]byte
	err     error
}

func (m *MockFileReader) ReadJsonSchema(resourcePath string) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	data, ok := m.schemas[resourcePath]
	if !ok {
		return nil, errors.New("schema not found: " + resourcePath)
	}
	return data, nil
}

// validMinimalSchema is the simplest valid JSON Schema.
var validMinimalSchema = []byte(`{"type": "object"}`)

// allSchemaFiles lists every JSON schema file used by NewValidator.
var allSchemaFiles = []string{
	"register_person_schema.json",
	"message_schema.json",
	"create_message_schema.json",
	"resend_verification_email.json",
	"password_reset_request.json",
	"update_profile_schema.json",
	"reset_password_with_token_schema.json",
	"change_password_schema.json",
	"register_branch_schema.json",
	"schedule_detail_schema.json",
	"update_schedule_schema.json",
	"schedule_exception_schema.json",
	"update_schedule_exception_schema.json",
	"franchise_schema.json",
	"register_motorcycle_schema.json",
	"update_motorcycle_schema.json",
	"create_evidence_schema.json",
	"completed_service_schema.json",
	"diagnostic_schema.json",
	"update_diagnostic_schema.json",
	"diagnostic_solution_schema.json",
	"update_schedule_detail_schema.json",
	"diagnostic_permission_schema.json",
	"branch_services_schema.json",
	"franchise_branch_schema.json",
	"update_status_schema.json",
}

// buildMockReader creates a MockFileReader that returns valid schemas for all files.
func buildMockReader() *MockFileReader {
	schemas := make(map[string][]byte, len(allSchemaFiles))
	for _, f := range allSchemaFiles {
		schemas[f] = validMinimalSchema
	}
	return &MockFileReader{schemas: schemas}
}

// TestNewValidator_Success validates that NewValidator correctly initializes
// all schema validators when all schema files are available.
func TestNewValidator_Success(t *testing.T) {
	reader := buildMockReader()

	v, err := schema.NewValidator(reader)

	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.NotNil(t, v.RegisterValidator)
	assert.NotNil(t, v.MessageValidator)
	assert.NotNil(t, v.RegisterBranchValidator)
	assert.NotNil(t, v.DiagnosticValidator)
}

// TestNewValidator_ReadError validates that NewValidator returns an error
// when a schema file cannot be read.
func TestNewValidator_ReadError(t *testing.T) {
	reader := &MockFileReader{
		err: errors.New("disk read failure"),
	}

	v, err := schema.NewValidator(reader)

	assert.Error(t, err)
	assert.Nil(t, v)
}

// TestNewValidator_EmptySchema validates that NewValidator returns ErrSchemaEmpty
// when a schema file returns nil bytes.
func TestNewValidator_EmptySchema(t *testing.T) {
	reader := buildMockReader()
	// Override one schema to return nil (empty)
	reader.schemas["register_person_schema.json"] = nil

	v, err := schema.NewValidator(reader)

	assert.Error(t, err)
	assert.Nil(t, v)
}

// TestNewValidator_InvalidSchemaJSON validates that NewValidator returns
// ErrSchemaCompileFailed when a schema JSON is malformed.
func TestNewValidator_InvalidSchemaJSON(t *testing.T) {
	reader := buildMockReader()
	// Override one schema to contain invalid JSON
	reader.schemas["register_person_schema.json"] = []byte(`{not valid json}`)

	v, err := schema.NewValidator(reader)

	assert.Error(t, err)
	assert.Nil(t, v)
}

// ============================================
// ValidateRegister Tests
// ============================================

func TestValidateRegister_ValidData_NoError(t *testing.T) {
	reader := buildMockReader()
	v, err := schema.NewValidator(reader)
	assert.NoError(t, err)

	// Minimal valid object for {"type": "object"} schema
	err = v.ValidateRegister(map[string]interface{}{})
	assert.NoError(t, err)
}

func TestValidateRegister_InvalidData_Error(t *testing.T) {
	reader := buildMockReader()
	v, err := schema.NewValidator(reader)
	assert.NoError(t, err)

	// A string is not a valid object
	err = v.ValidateRegister("not an object")
	assert.Error(t, err)
}

func TestValidateRegister_NilValidator_Error(t *testing.T) {
	v := &schema.Validators{}
	err := v.ValidateRegister(map[string]interface{}{})
	assert.Error(t, err)
}

// ============================================
// ValidateMessage Tests
// ============================================

func TestValidateMessage_ValidData_NoError(t *testing.T) {
	reader := buildMockReader()
	v, err := schema.NewValidator(reader)
	assert.NoError(t, err)

	err = v.ValidateMessage(map[string]interface{}{})
	assert.NoError(t, err)
}

func TestValidateMessage_InvalidData_Error(t *testing.T) {
	reader := buildMockReader()
	v, err := schema.NewValidator(reader)
	assert.NoError(t, err)

	err = v.ValidateMessage("not an object")
	assert.Error(t, err)
}

func TestValidateMessage_NilValidator_Error(t *testing.T) {
	v := &schema.Validators{}
	err := v.ValidateMessage(map[string]interface{}{})
	assert.Error(t, err)
}

// ============================================
// ValidateRegisterMotorcycle Tests
// ============================================

func TestValidateRegisterMotorcycle_ValidData_NoError(t *testing.T) {
	reader := buildMockReader()
	v, err := schema.NewValidator(reader)
	assert.NoError(t, err)

	err = v.ValidateRegisterMotorcycle(map[string]interface{}{})
	assert.NoError(t, err)
}

func TestValidateRegisterMotorcycle_InvalidData_Error(t *testing.T) {
	reader := buildMockReader()
	v, err := schema.NewValidator(reader)
	assert.NoError(t, err)

	err = v.ValidateRegisterMotorcycle(42)
	assert.Error(t, err)
}

func TestValidateRegisterMotorcycle_NilValidator_Error(t *testing.T) {
	v := &schema.Validators{}
	err := v.ValidateRegisterMotorcycle(map[string]interface{}{})
	assert.Error(t, err)
}

// ============================================
// ValidateUpdateMotorcycle Tests
// ============================================

func TestValidateUpdateMotorcycle_ValidData_NoError(t *testing.T) {
	reader := buildMockReader()
	v, err := schema.NewValidator(reader)
	assert.NoError(t, err)

	err = v.ValidateUpdateMotorcycle(map[string]interface{}{})
	assert.NoError(t, err)
}

func TestValidateUpdateMotorcycle_InvalidData_Error(t *testing.T) {
	reader := buildMockReader()
	v, err := schema.NewValidator(reader)
	assert.NoError(t, err)

	err = v.ValidateUpdateMotorcycle(true)
	assert.Error(t, err)
}

func TestValidateUpdateMotorcycle_NilValidator_Error(t *testing.T) {
	v := &schema.Validators{}
	err := v.ValidateUpdateMotorcycle(map[string]interface{}{})
	assert.Error(t, err)
}

// ============================================
// DefaultFileReader Tests
// ============================================

func TestDefaultFileReader_ReadJsonSchema_Success(t *testing.T) {
	reader := &schema.DefaultFileReader{}

	// Read an actual schema file that exists in the project
	data, err := reader.ReadJsonSchema("register_person_schema.json")

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.NotEmpty(t, data)
}

func TestDefaultFileReader_ReadJsonSchema_NotFound(t *testing.T) {
	reader := &schema.DefaultFileReader{}

	data, err := reader.ReadJsonSchema("nonexistent_schema.json")

	assert.Error(t, err)
	assert.Nil(t, data)
}
