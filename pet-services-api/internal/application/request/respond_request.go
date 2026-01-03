package request

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainProvider "github.com/guilherme/pet-services-api/internal/domain/provider"
	domainRequest "github.com/guilherme/pet-services-api/internal/domain/request"
)

// AcceptRequestUseCase permite ao prestador aceitar uma solicitação.
type AcceptRequestUseCase struct {
	requestRepo  domainRequest.Repository
	providerRepo domainProvider.Repository
	logger       *slog.Logger
}

func NewAcceptRequestUseCase(requestRepo domainRequest.Repository, providerRepo domainProvider.Repository, logger *slog.Logger) *AcceptRequestUseCase {
	return &AcceptRequestUseCase{requestRepo: requestRepo, providerRepo: providerRepo, logger: logging.EnsureLogger(logger)}
}

// AcceptRequestInput dados para aceitar uma solicitação.
type AcceptRequestInput struct {
	RequestID  uuid.UUID
	ProviderID uuid.UUID
}

func (uc *AcceptRequestUseCase) Execute(ctx context.Context, input AcceptRequestInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "AcceptRequestUseCase", slog.String("request_id", input.RequestID.String()), slog.String("provider_id", input.ProviderID.String()))(&err)

	if input.RequestID == uuid.Nil {
		err = fmt.Errorf("requestID é obrigatório")
		return err
	}
	if input.ProviderID == uuid.Nil {
		err = domainProvider.ErrProviderNotFound
		return err
	}

	req, err := uc.requestRepo.FindByID(ctx, input.RequestID)
	if err != nil {
		return err
	}

	if req.ProviderID != input.ProviderID {
		err = domainProvider.ErrProviderNotFound
		return err
	}

	provider, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		err = domainProvider.ErrProviderNotFound
		return err
	}
	if !provider.IsActive {
		err = domainProvider.ErrProviderNotActive
		return err
	}

	if err := req.Accept(); err != nil {
		return err
	}

	if err := uc.requestRepo.Update(ctx, req); err != nil {
		err = fmt.Errorf("falha ao aceitar solicitação: %w", err)
		return err
	}

	return nil
}

// RejectRequestUseCase permite ao prestador rejeitar uma solicitação.
type RejectRequestUseCase struct {
	requestRepo  domainRequest.Repository
	providerRepo domainProvider.Repository
	logger       *slog.Logger
}

func NewRejectRequestUseCase(requestRepo domainRequest.Repository, providerRepo domainProvider.Repository, logger *slog.Logger) *RejectRequestUseCase {
	return &RejectRequestUseCase{requestRepo: requestRepo, providerRepo: providerRepo, logger: logging.EnsureLogger(logger)}
}

// RejectRequestInput dados para rejeitar uma solicitação.
type RejectRequestInput struct {
	RequestID  uuid.UUID
	ProviderID uuid.UUID
	Reason     string
}

func (uc *RejectRequestUseCase) Execute(ctx context.Context, input RejectRequestInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "RejectRequestUseCase", slog.String("request_id", input.RequestID.String()), slog.String("provider_id", input.ProviderID.String()))(&err)

	if input.RequestID == uuid.Nil {
		err = fmt.Errorf("requestID é obrigatório")
		return err
	}
	if input.ProviderID == uuid.Nil {
		err = domainProvider.ErrProviderNotFound
		return err
	}

	req, err := uc.requestRepo.FindByID(ctx, input.RequestID)
	if err != nil {
		return err
	}

	if req.ProviderID != input.ProviderID {
		err = domainProvider.ErrProviderNotFound
		return err
	}

	provider, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		err = domainProvider.ErrProviderNotFound
		return err
	}
	if !provider.IsActive {
		err = domainProvider.ErrProviderNotActive
		return err
	}

	reason := strings.TrimSpace(input.Reason)
	if reason == "" {
		err = fmt.Errorf("motivo da rejeição é obrigatório")
		return err
	}

	if err := req.Reject(reason); err != nil {
		return err
	}

	if err := uc.requestRepo.Update(ctx, req); err != nil {
		err = fmt.Errorf("falha ao rejeitar solicitação: %w", err)
		return err
	}

	return nil
}

// CompleteRequestUseCase permite ao prestador concluir uma solicitação aceita.
type CompleteRequestUseCase struct {
	requestRepo  domainRequest.Repository
	providerRepo domainProvider.Repository
	logger       *slog.Logger
}

func NewCompleteRequestUseCase(requestRepo domainRequest.Repository, providerRepo domainProvider.Repository, logger *slog.Logger) *CompleteRequestUseCase {
	return &CompleteRequestUseCase{requestRepo: requestRepo, providerRepo: providerRepo, logger: logging.EnsureLogger(logger)}
}

// CompleteRequestInput dados para concluir uma solicitação.
type CompleteRequestInput struct {
	RequestID  uuid.UUID
	ProviderID uuid.UUID
}

func (uc *CompleteRequestUseCase) Execute(ctx context.Context, input CompleteRequestInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "CompleteRequestUseCase", slog.String("request_id", input.RequestID.String()), slog.String("provider_id", input.ProviderID.String()))(&err)

	if input.RequestID == uuid.Nil {
		err = fmt.Errorf("requestID é obrigatório")
		return err
	}
	if input.ProviderID == uuid.Nil {
		err = domainProvider.ErrProviderNotFound
		return err
	}

	req, err := uc.requestRepo.FindByID(ctx, input.RequestID)
	if err != nil {
		return err
	}

	if req.ProviderID != input.ProviderID {
		err = domainProvider.ErrProviderNotFound
		return err
	}

	if _, err := uc.providerRepo.FindByID(ctx, input.ProviderID); err != nil {
		err = domainProvider.ErrProviderNotFound
		return err
	}

	if err := req.Complete(); err != nil {
		return err
	}

	if err := uc.requestRepo.Update(ctx, req); err != nil {
		err = fmt.Errorf("falha ao concluir solicitação: %w", err)
		return err
	}

	return nil
}

// CancelRequestUseCase permite ao dono cancelar uma solicitação.
type CancelRequestUseCase struct {
	requestRepo domainRequest.Repository
	logger      *slog.Logger
}

func NewCancelRequestUseCase(requestRepo domainRequest.Repository, logger *slog.Logger) *CancelRequestUseCase {
	return &CancelRequestUseCase{requestRepo: requestRepo, logger: logging.EnsureLogger(logger)}
}

// CancelRequestInput dados para cancelar uma solicitação.
type CancelRequestInput struct {
	RequestID uuid.UUID
	OwnerID   uuid.UUID
}

func (uc *CancelRequestUseCase) Execute(ctx context.Context, input CancelRequestInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "CancelRequestUseCase", slog.String("request_id", input.RequestID.String()), slog.String("owner_id", input.OwnerID.String()))(&err)

	if input.RequestID == uuid.Nil {
		err = fmt.Errorf("requestID é obrigatório")
		return err
	}
	if input.OwnerID == uuid.Nil {
		err = fmt.Errorf("ownerID é obrigatório")
		return err
	}

	req, err := uc.requestRepo.FindByID(ctx, input.RequestID)
	if err != nil {
		return err
	}

	if req.OwnerID != input.OwnerID {
		err = fmt.Errorf("não autorizado a cancelar esta solicitação")
		return err
	}

	if err := req.Cancel(); err != nil {
		return err
	}

	if err := uc.requestRepo.Update(ctx, req); err != nil {
		err = fmt.Errorf("falha ao cancelar solicitação: %w", err)
		return err
	}

	return nil
}
