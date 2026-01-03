package request

import (
    "context"
    "fmt"
    "strings"

    "github.com/google/uuid"

    domainProvider "github.com/guilherme/pet-services-api/internal/domain/provider"
    domainRequest "github.com/guilherme/pet-services-api/internal/domain/request"
)

// AcceptRequestUseCase permite ao prestador aceitar uma solicitação.
type AcceptRequestUseCase struct {
    requestRepo  domainRequest.Repository
    providerRepo domainProvider.Repository
}

func NewAcceptRequestUseCase(requestRepo domainRequest.Repository, providerRepo domainProvider.Repository) *AcceptRequestUseCase {
    return &AcceptRequestUseCase{requestRepo: requestRepo, providerRepo: providerRepo}
}

// AcceptRequestInput dados para aceitar uma solicitação.
type AcceptRequestInput struct {
    RequestID  uuid.UUID
    ProviderID uuid.UUID
}

func (uc *AcceptRequestUseCase) Execute(ctx context.Context, input AcceptRequestInput) error {
    if input.RequestID == uuid.Nil {
        return fmt.Errorf("requestID é obrigatório")
    }
    if input.ProviderID == uuid.Nil {
        return domainProvider.ErrProviderNotFound
    }

    req, err := uc.requestRepo.FindByID(ctx, input.RequestID)
    if err != nil {
        return err
    }

    if req.ProviderID != input.ProviderID {
        return domainProvider.ErrProviderNotFound
    }

    provider, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
    if err != nil {
        return domainProvider.ErrProviderNotFound
    }
    if !provider.IsActive {
        return domainProvider.ErrProviderNotActive
    }

    if err := req.Accept(); err != nil {
        return err
    }

    if err := uc.requestRepo.Update(ctx, req); err != nil {
        return fmt.Errorf("falha ao aceitar solicitação: %w", err)
    }

    return nil
}

// RejectRequestUseCase permite ao prestador rejeitar uma solicitação.
type RejectRequestUseCase struct {
    requestRepo  domainRequest.Repository
    providerRepo domainProvider.Repository
}

func NewRejectRequestUseCase(requestRepo domainRequest.Repository, providerRepo domainProvider.Repository) *RejectRequestUseCase {
    return &RejectRequestUseCase{requestRepo: requestRepo, providerRepo: providerRepo}
}

// RejectRequestInput dados para rejeitar uma solicitação.
type RejectRequestInput struct {
    RequestID  uuid.UUID
    ProviderID uuid.UUID
    Reason     string
}

func (uc *RejectRequestUseCase) Execute(ctx context.Context, input RejectRequestInput) error {
    if input.RequestID == uuid.Nil {
        return fmt.Errorf("requestID é obrigatório")
    }
    if input.ProviderID == uuid.Nil {
        return domainProvider.ErrProviderNotFound
    }

    req, err := uc.requestRepo.FindByID(ctx, input.RequestID)
    if err != nil {
        return err
    }

    if req.ProviderID != input.ProviderID {
        return domainProvider.ErrProviderNotFound
    }

    provider, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
    if err != nil {
        return domainProvider.ErrProviderNotFound
    }
    if !provider.IsActive {
        return domainProvider.ErrProviderNotActive
    }

    reason := strings.TrimSpace(input.Reason)
    if reason == "" {
        return fmt.Errorf("motivo da rejeição é obrigatório")
    }

    if err := req.Reject(reason); err != nil {
        return err
    }

    if err := uc.requestRepo.Update(ctx, req); err != nil {
        return fmt.Errorf("falha ao rejeitar solicitação: %w", err)
    }

    return nil
}

// CompleteRequestUseCase permite ao prestador concluir uma solicitação aceita.
type CompleteRequestUseCase struct {
    requestRepo  domainRequest.Repository
    providerRepo domainProvider.Repository
}

func NewCompleteRequestUseCase(requestRepo domainRequest.Repository, providerRepo domainProvider.Repository) *CompleteRequestUseCase {
    return &CompleteRequestUseCase{requestRepo: requestRepo, providerRepo: providerRepo}
}

// CompleteRequestInput dados para concluir uma solicitação.
type CompleteRequestInput struct {
    RequestID  uuid.UUID
    ProviderID uuid.UUID
}

func (uc *CompleteRequestUseCase) Execute(ctx context.Context, input CompleteRequestInput) error {
    if input.RequestID == uuid.Nil {
        return fmt.Errorf("requestID é obrigatório")
    }
    if input.ProviderID == uuid.Nil {
        return domainProvider.ErrProviderNotFound
    }

    req, err := uc.requestRepo.FindByID(ctx, input.RequestID)
    if err != nil {
        return err
    }

    if req.ProviderID != input.ProviderID {
        return domainProvider.ErrProviderNotFound
    }

    if _, err := uc.providerRepo.FindByID(ctx, input.ProviderID); err != nil {
        return domainProvider.ErrProviderNotFound
    }

    if err := req.Complete(); err != nil {
        return err
    }

    if err := uc.requestRepo.Update(ctx, req); err != nil {
        return fmt.Errorf("falha ao concluir solicitação: %w", err)
    }

    return nil
}

// CancelRequestUseCase permite ao dono cancelar uma solicitação.
type CancelRequestUseCase struct {
    requestRepo domainRequest.Repository
}

func NewCancelRequestUseCase(requestRepo domainRequest.Repository) *CancelRequestUseCase {
    return &CancelRequestUseCase{requestRepo: requestRepo}
}

// CancelRequestInput dados para cancelar uma solicitação.
type CancelRequestInput struct {
    RequestID uuid.UUID
    OwnerID   uuid.UUID
}

func (uc *CancelRequestUseCase) Execute(ctx context.Context, input CancelRequestInput) error {
    if input.RequestID == uuid.Nil {
        return fmt.Errorf("requestID é obrigatório")
    }
    if input.OwnerID == uuid.Nil {
        return fmt.Errorf("ownerID é obrigatório")
    }

    req, err := uc.requestRepo.FindByID(ctx, input.RequestID)
    if err != nil {
        return err
    }

    if req.OwnerID != input.OwnerID {
        return fmt.Errorf("não autorizado a cancelar esta solicitação")
    }

    if err := req.Cancel(); err != nil {
        return err
    }

    if err := uc.requestRepo.Update(ctx, req); err != nil {
        return fmt.Errorf("falha ao cancelar solicitação: %w", err)
    }

    return nil
}
