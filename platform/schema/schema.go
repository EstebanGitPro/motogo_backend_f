package schema

import (
	"io"
	"os"
	"path/filepath"

	"github.com/EstebanGitPro/motogo-backend/tools/utils"
	"github.com/kaptinlin/jsonschema"
)

type Validators struct {
	FileReader                       FileReaderInterface
	RegisterValidator                *jsonschema.Schema
	MessageValidator                 *jsonschema.Schema
	CreateMessageValidator           *jsonschema.Schema
	ResendVerificationValidator      *jsonschema.Schema
	PasswordResetValidator           *jsonschema.Schema
	UpdateProfileValidator           *jsonschema.Schema
	ResetPasswordWithTokenValidator  *jsonschema.Schema
	ChangePasswordValidator          *jsonschema.Schema
	RegisterBranchValidator          *jsonschema.Schema // HU59
	ScheduleDetailValidator          *jsonschema.Schema // HU6-9 (schedule time slots)
	UpdateScheduleValidator          *jsonschema.Schema // HU31 (update schedule)
	ScheduleExceptionValidator       *jsonschema.Schema // HU20 (create exception)
	UpdateScheduleExceptionValidator *jsonschema.Schema // HU21 (update exception)
	FranchiseValidator               *jsonschema.Schema // HU26-29 (franchises)
	RegisterMotorcycleValidator      *jsonschema.Schema // HU43 (register motorcycle)
	UpdateMotorcycleValidator        *jsonschema.Schema // HU44 (update motorcycle)
	CreateEvidenceValidator          *jsonschema.Schema // HU16 (create evidence)
	CompletedServiceValidator        *jsonschema.Schema // HU64 (completed services)
	DiagnosticValidator              *jsonschema.Schema // HU11 (create diagnostic)
	UpdateDiagnosticValidator        *jsonschema.Schema // HU12 (update diagnostic)
	DiagnosticSolutionValidator      *jsonschema.Schema // PATCH diagnostic solution
	UpdateScheduleDetailValidator    *jsonschema.Schema // HU7 (update schedule detail)
	DiagnosticPermissionValidator    *jsonschema.Schema // diagnostic permission
	BranchServicesValidator          *jsonschema.Schema // branch services association
	FranchiseBranchValidator         *jsonschema.Schema // franchise branch association
	UpdateStatusValidator            *jsonschema.Schema // HU74 (update status)
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
	defer func() { _ = data.Close() }() // Close error intentionally ignored

	return io.ReadAll(data)
}

func NewValidator(fileReader FileReaderInterface) (*Validators, error) {
	validator := &Validators{
		FileReader: fileReader,
	}

	type schemaEntry struct {
		field **jsonschema.Schema
		name  string
	}

	schemas := []schemaEntry{
		{&validator.RegisterValidator, "register_person_schema.json"},
		{&validator.MessageValidator, "message_schema.json"},
		{&validator.CreateMessageValidator, "create_message_schema.json"},
		{&validator.ResendVerificationValidator, "resend_verification_email.json"},
		{&validator.PasswordResetValidator, "password_reset_request.json"},
		{&validator.UpdateProfileValidator, "update_profile_schema.json"},
		{&validator.ResetPasswordWithTokenValidator, "reset_password_with_token_schema.json"},
		{&validator.ChangePasswordValidator, "change_password_schema.json"},
		{&validator.RegisterBranchValidator, "register_branch_schema.json"},
		{&validator.ScheduleDetailValidator, "schedule_detail_schema.json"},
		{&validator.UpdateScheduleValidator, "update_schedule_schema.json"},
		{&validator.ScheduleExceptionValidator, "schedule_exception_schema.json"},
		{&validator.UpdateScheduleExceptionValidator, "update_schedule_exception_schema.json"},
		{&validator.FranchiseValidator, "franchise_schema.json"},
		{&validator.RegisterMotorcycleValidator, "register_motorcycle_schema.json"},
		{&validator.UpdateMotorcycleValidator, "update_motorcycle_schema.json"},
		{&validator.CreateEvidenceValidator, "create_evidence_schema.json"},
		{&validator.CompletedServiceValidator, "completed_service_schema.json"},
		{&validator.DiagnosticValidator, "diagnostic_schema.json"},
		{&validator.UpdateDiagnosticValidator, "update_diagnostic_schema.json"},
		{&validator.DiagnosticSolutionValidator, "diagnostic_solution_schema.json"},
		{&validator.UpdateScheduleDetailValidator, "update_schedule_detail_schema.json"},
		{&validator.DiagnosticPermissionValidator, "diagnostic_permission_schema.json"},
		{&validator.BranchServicesValidator, "branch_services_schema.json"},
		{&validator.FranchiseBranchValidator, "franchise_branch_schema.json"},
		{&validator.UpdateStatusValidator, "update_status_schema.json"},
	}

	for _, s := range schemas {
		sch, err := validator.createSchema(s.name)
		if err != nil {
			return nil, err
		}
		*s.field = sch
	}

	return validator, nil
}

func (v *Validators) createSchema(resourcePath string) (*jsonschema.Schema, error) {
	compiler := jsonschema.NewCompiler()
	compiler.AssertFormat = true
	schemaJSON, err := v.FileReader.ReadJsonSchema(resourcePath)
	if err != nil {
		return nil, ErrSchemaReadFailed
	}

	if schemaJSON == nil {
		return nil, ErrSchemaEmpty
	}

	schema, err := compiler.Compile(schemaJSON)
	if err != nil {
		return nil, ErrSchemaCompileFailed
	}

	return schema, nil
}

// ValidateRegister validates data against the register person schema
func (v *Validators) ValidateRegister(data interface{}) error {
	if v.RegisterValidator == nil {
		return ErrSchemaEmpty
	}

	result := v.RegisterValidator.Validate(data)
	if !result.IsValid() {
		return ErrValidationFailed
	}

	return nil
}

// ValidateMessage validates data against the message schema
func (v *Validators) ValidateMessage(data interface{}) error {
	if v.MessageValidator == nil {
		return ErrSchemaEmpty
	}

	result := v.MessageValidator.Validate(data)
	if !result.IsValid() {
		return ErrValidationFailed
	}

	return nil
}

// ValidateRegisterMotorcycle validates data against the register motorcycle schema (HU43)
func (v *Validators) ValidateRegisterMotorcycle(data interface{}) error {
	if v.RegisterMotorcycleValidator == nil {
		return ErrSchemaEmpty
	}

	result := v.RegisterMotorcycleValidator.Validate(data)
	if !result.IsValid() {
		return ErrValidationFailed
	}

	return nil
}

// ValidateUpdateMotorcycle validates data against the update motorcycle schema (HU44)
func (v *Validators) ValidateUpdateMotorcycle(data interface{}) error {
	if v.UpdateMotorcycleValidator == nil {
		return ErrSchemaEmpty
	}

	result := v.UpdateMotorcycleValidator.Validate(data)
	if !result.IsValid() {
		return ErrValidationFailed
	}

	return nil
}
