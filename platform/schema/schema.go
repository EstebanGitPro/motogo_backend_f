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

	register, err := validator.createSchema("register_person_schema.json")
	if err != nil {
		return nil, err
	}

	message, err := validator.createSchema("message_schema.json")
	if err != nil {
		return nil, err
	}

	createMessage, err := validator.createSchema("create_message_schema.json")
	if err != nil {
		return nil, err
	}

	resendVerification, err := validator.createSchema("resend_verification_email.json")
	if err != nil {
		return nil, err
	}

	passwordReset, err := validator.createSchema("password_reset_request.json")
	if err != nil {
		return nil, err
	}

	updateProfile, err := validator.createSchema("update_profile_schema.json")
	if err != nil {
		return nil, err
	}

	resetPasswordWithToken, err := validator.createSchema("reset_password_with_token_schema.json")
	if err != nil {
		return nil, err
	}

	changePassword, err := validator.createSchema("change_password_schema.json")
	if err != nil {
		return nil, err
	}

	registerBranch, err := validator.createSchema("register_branch_schema.json")
	if err != nil {
		return nil, err
	}

	scheduleDetail, err := validator.createSchema("schedule_detail_schema.json")
	if err != nil {
		return nil, err
	}

	updateSchedule, err := validator.createSchema("update_schedule_schema.json")
	if err != nil {
		return nil, err
	}

	scheduleException, err := validator.createSchema("schedule_exception_schema.json")
	if err != nil {
		return nil, err
	}

	updateScheduleException, err := validator.createSchema("update_schedule_exception_schema.json")
	if err != nil {
		return nil, err
	}

	franchise, err := validator.createSchema("franchise_schema.json")
	if err != nil {
		return nil, err
	}

	registerMotorcycle, err := validator.createSchema("register_motorcycle_schema.json")
	if err != nil {
		return nil, err
	}

	createEvidence, err := validator.createSchema("create_evidence_schema.json")
	if err != nil {
		return nil, err
	}

	completedService, err := validator.createSchema("completed_service_schema.json")
	if err != nil {
		return nil, err
	}

	diagnostic, err := validator.createSchema("diagnostic_schema.json")
	if err != nil {
		return nil, err
	}

	updateDiagnostic, err := validator.createSchema("update_diagnostic_schema.json")
	if err != nil {
		return nil, err
	}

	diagnosticSolution, err := validator.createSchema("diagnostic_solution_schema.json")
	if err != nil {
		return nil, err
	}

	updateScheduleDetail, err := validator.createSchema("update_schedule_detail_schema.json")
	if err != nil {
		return nil, err
	}

	diagnosticPermission, err := validator.createSchema("diagnostic_permission_schema.json")
	if err != nil {
		return nil, err
	}

	branchServices, err := validator.createSchema("branch_services_schema.json")
	if err != nil {
		return nil, err
	}

	franchiseBranch, err := validator.createSchema("franchise_branch_schema.json")
	if err != nil {
		return nil, err
	}

	updateStatus, err := validator.createSchema("update_status_schema.json")
	if err != nil {
		return nil, err
	}

	validator.RegisterValidator = register
	validator.MessageValidator = message
	validator.CreateMessageValidator = createMessage
	validator.ResendVerificationValidator = resendVerification
	validator.PasswordResetValidator = passwordReset
	validator.UpdateProfileValidator = updateProfile
	validator.ResetPasswordWithTokenValidator = resetPasswordWithToken
	validator.ChangePasswordValidator = changePassword
	validator.RegisterBranchValidator = registerBranch
	validator.ScheduleDetailValidator = scheduleDetail
	validator.UpdateScheduleValidator = updateSchedule
	validator.ScheduleExceptionValidator = scheduleException
	validator.UpdateScheduleExceptionValidator = updateScheduleException
	validator.FranchiseValidator = franchise
	validator.RegisterMotorcycleValidator = registerMotorcycle
	validator.CreateEvidenceValidator = createEvidence
	validator.CompletedServiceValidator = completedService
	validator.DiagnosticValidator = diagnostic
	validator.UpdateDiagnosticValidator = updateDiagnostic
	validator.DiagnosticSolutionValidator = diagnosticSolution
	validator.UpdateScheduleDetailValidator = updateScheduleDetail
	validator.DiagnosticPermissionValidator = diagnosticPermission
	validator.BranchServicesValidator = branchServices
	validator.FranchiseBranchValidator = franchiseBranch
	validator.UpdateStatusValidator = updateStatus

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
