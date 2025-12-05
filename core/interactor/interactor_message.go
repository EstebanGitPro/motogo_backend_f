package interactor

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// MessageInteractor orchestrates message operations
type MessageInteractor struct {
	service input.MessageService
	logger  logger.Logger
}

// NewMessageInteractor creates a new message interactor
func NewMessageInteractor(service input.MessageService, log logger.Logger) *MessageInteractor {
	return &MessageInteractor{
		service: service,
		logger:  log,
	}
}

// CreateMessage creates a new system message with transaction handling
func (i *MessageInteractor) CreateMessage(ctx context.Context, message domain.Message) (result *domain.Message, err error) {
	i.logger.Info(logger.LogMessageCreate, message.ToLogger())

	// PASO 1: Validate message
	if err = i.service.ValidateMessage(ctx, message); err != nil {
		i.logger.Error(logger.LogMessageInteractorCreateStep1Error, "error", err)
		return
	}
	i.logger.Success(logger.LogMessageInteractorCreateStep1OK)

	// PASO 2: Begin transaction
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		i.logger.Error(logger.LogMessageInteractorCreateStep2Error, "error", err)
		return
	}
	i.logger.Success(logger.LogMessageInteractorCreateStep2OK)

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				i.logger.Error(logger.LogMessageInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				i.logger.Warn(logger.LogMessageInteractorRollbackOK)
			}
		}
	}()

	// PASO 3: Save message to DB
	if err = i.service.SaveMessageToDB(ctx, tx, message); err != nil {
		i.logger.Error(logger.LogMessageInteractorCreateStep3Error, "error", err)
		return
	}
	i.logger.Success(logger.LogMessageInteractorCreateStep3OK)

	// COMMIT: Confirm transaction
	if err = tx.Commit(); err != nil {
		i.logger.Error(logger.LogMessageInteractorCreateCommitErr, "error", err)
		return
	}
	i.logger.Success(logger.LogMessageInteractorCreateCommitOK)

	result = &message
	i.logger.Success(logger.LogMessageInteractorCreateComplete, message.ToLogger())

	err = nil // ensure defer does NOT execute rollback
	return
}

// UpdateMessage updates an existing system message with transaction handling
func (i *MessageInteractor) UpdateMessage(ctx context.Context, message domain.Message) (result *domain.Message, err error) {
	i.logger.Info(logger.LogMessageUpdate, message.ToLogger())

	// PASO 1: Validate message exists
	_, err = i.service.GetMessageByID(ctx, message.ID)
	if err != nil {
		i.logger.Error(logger.LogMessageInteractorUpdateStep1Error, "error", err)
		return nil, err
	}
	i.logger.Success(logger.LogMessageInteractorUpdateStep1OK)

	// PASO 2: Validate message data
	if err = i.service.ValidateMessage(ctx, message); err != nil {
		i.logger.Error(logger.LogMessageInteractorUpdateStep2Error, "error", err)
		return
	}
	i.logger.Success(logger.LogMessageInteractorUpdateStep2OK)

	// PASO 3: Begin transaction
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		i.logger.Error(logger.LogMessageInteractorUpdateStep3Error, "error", err)
		return
	}
	i.logger.Success(logger.LogMessageInteractorUpdateStep3OK)

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				i.logger.Error(logger.LogMessageInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				i.logger.Warn(logger.LogMessageInteractorRollbackOK)
			}
		}
	}()

	// PASO 4: Update message in DB
	if err = i.service.UpdateMessageInDB(ctx, tx, message); err != nil {
		i.logger.Error(logger.LogMessageInteractorUpdateStep4Error, "error", err)
		return
	}
	i.logger.Success(logger.LogMessageInteractorUpdateStep4OK)

	// COMMIT: Confirm transaction
	if err = tx.Commit(); err != nil {
		i.logger.Error(logger.LogMessageInteractorUpdateCommitErr, "error", err)
		return
	}
	i.logger.Success(logger.LogMessageInteractorUpdateCommitOK)

	result = &message
	i.logger.Success(logger.LogMessageInteractorUpdateComplete, message.ToLogger())

	// Refresh cache after update
	i.logger.Info(logger.LogMessageCacheRefresh)

	err = nil // ensure defer does NOT execute rollback
	return
}

// DeleteMessage deletes a system message with transaction handling
func (i *MessageInteractor) DeleteMessage(ctx context.Context, id string) (err error) {
	i.logger.Info(logger.LogMessageDelete, "id", id)

	// PASO 1: Validate message exists
	_, err = i.service.GetMessageByID(ctx, id)
	if err != nil {
		i.logger.Error(logger.LogMessageInteractorDeleteStep1Error, "error", err)
		return err
	}
	i.logger.Success(logger.LogMessageInteractorDeleteStep1OK)

	// PASO 2: Begin transaction
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		i.logger.Error(logger.LogMessageInteractorDeleteStep2Error, "error", err)
		return err
	}
	i.logger.Success(logger.LogMessageInteractorDeleteStep2OK)

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				i.logger.Error(logger.LogMessageInteractorRollbackError,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				i.logger.Warn(logger.LogMessageInteractorRollbackOK)
			}
		}
	}()

	// PASO 3: Delete message from DB
	if err = i.service.DeleteMessageFromDB(ctx, tx, id); err != nil {
		i.logger.Error(logger.LogMessageInteractorDeleteStep3Error, "error", err)
		return err
	}
	i.logger.Success(logger.LogMessageInteractorDeleteStep3OK)

	// COMMIT: Confirm transaction
	if err = tx.Commit(); err != nil {
		i.logger.Error(logger.LogMessageInteractorDeleteCommitErr, "error", err)
		return err
	}
	i.logger.Success(logger.LogMessageInteractorDeleteCommitOK)

	i.logger.Success(logger.LogMessageInteractorDeleteComplete, "id", id)

	// Refresh cache after delete
	i.logger.Info(logger.LogMessageCacheRefresh)

	err = nil // ensure defer does NOT execute rollback
	return
}

// GetMessageByID retrieves a message by ID (read-only, no transaction)
func (i *MessageInteractor) GetMessageByID(ctx context.Context, id string) (*domain.Message, error) {
	i.logger.Debug(logger.LogMessageGet, "id", id)

	message, err := i.service.GetMessageByID(ctx, id)
	if err != nil {
		i.logger.Error(logger.LogMessageGetError, "id", id, "error", err)
		return nil, err
	}

	i.logger.Debug(logger.LogMessageGetOK, message.ToLogger())
	return message, nil
}

// GetMessageByCode retrieves a message by code (read-only, no transaction)
func (i *MessageInteractor) GetMessageByCode(ctx context.Context, code string) (*domain.Message, error) {
	i.logger.Debug(logger.LogMessageGet, "code", code)

	message, err := i.service.GetMessageByCode(ctx, code)
	if err != nil {
		i.logger.Error(logger.LogMessageGetError, "code", code, "error", err)
		return nil, err
	}

	i.logger.Debug(logger.LogMessageGetOK, message.ToLogger())
	return message, nil
}

// ListMessages retrieves messages with optional filters (read-only, no transaction)
func (i *MessageInteractor) ListMessages(ctx context.Context, filters map[string]interface{}) ([]domain.Message, error) {
	i.logger.Debug(logger.LogMessageList, "filters", filters)

	messages, err := i.service.ListMessages(ctx, filters)
	if err != nil {
		i.logger.Error(logger.LogMessageListError, "error", err)
		return nil, err
	}

	i.logger.Debug(logger.LogMessageListOK, "count", len(messages))
	return messages, nil
}

// ListActiveMessages retrieves only active messages (read-only, no transaction)
func (i *MessageInteractor) ListActiveMessages(ctx context.Context) ([]domain.Message, error) {
	i.logger.Debug(logger.LogMessageList, "filter", "active_only")

	messages, err := i.service.ListActiveMessages(ctx)
	if err != nil {
		i.logger.Error(logger.LogMessageListError, "error", err)
		return nil, err
	}

	i.logger.Debug(logger.LogMessageListOK, "count", len(messages))
	return messages, nil
}
