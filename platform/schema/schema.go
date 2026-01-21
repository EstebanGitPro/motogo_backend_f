package schema

import (
	"io"
	"os"
	"path/filepath"

	"github.com/EstebanGitPro/motogo-backend/tools/utils"
	"github.com/kaptinlin/jsonschema"
)

type Validators struct {
	FileReader                      FileReaderInterface
	RegisterValidator               *jsonschema.Schema
	MessageValidator                *jsonschema.Schema
	ResendVerificationValidator     *jsonschema.Schema
	PasswordResetValidator          *jsonschema.Schema
	UpdateProfileValidator          *jsonschema.Schema
	ResetPasswordWithTokenValidator *jsonschema.Schema
	ChangePasswordValidator         *jsonschema.Schema
	RegisterBranchValidator         *jsonschema.Schema // HU59
	ScheduleDetailValidator         *jsonschema.Schema // HU6-9 (schedule time slots)
	UpdateScheduleValidator         *jsonschema.Schema // HU31 (update schedule)
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

	message, err := validator.createSchema("message_schema.json")
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

	validator.RegisterValidator = register
	validator.MessageValidator = message
	validator.ResendVerificationValidator = resendVerification
	validator.PasswordResetValidator = passwordReset
	validator.UpdateProfileValidator = updateProfile
	validator.ResetPasswordWithTokenValidator = resetPasswordWithToken
	validator.ChangePasswordValidator = changePassword
	validator.RegisterBranchValidator = registerBranch
	validator.ScheduleDetailValidator = scheduleDetail
	validator.UpdateScheduleValidator = updateSchedule

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
