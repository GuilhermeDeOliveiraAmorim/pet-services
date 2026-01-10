package provider

import (
	"context"
	"fmt"
	"pet-services-api/internal/application/exceptions"
	"pet-services-api/internal/application/logging"
	"pet-services-api/internal/domain/provider"

	"github.com/google/uuid"
)

type MinioUploader interface {
	Upload(ctx context.Context, objectName string, fileData []byte, contentType string) (string, error)
}

// Ajusta o caso de uso para receber o serviço Minio
type AddServicePhotoUseCase struct {
	logger       logging.LoggerService
	providerRepo provider.Repository
	minio        MinioUploader
}

func NewAddServicePhotoUseCase(logger logging.LoggerService, providerRepo provider.Repository, minio MinioUploader) *AddServicePhotoUseCase {
	return &AddServicePhotoUseCase{
		logger:       logger,
		providerRepo: providerRepo,
		minio:        minio,
	}
}

// Agora o input recebe o arquivo (bytes) e contentType
// Se quiser manter compatibilidade, pode usar URL como fallback
// Aqui, prioriza upload de arquivo
type AddServicePhotoInput struct {
	ServiceID   uuid.UUID
	FileName    string // nome do arquivo para Minio
	FileData    []byte // dados do arquivo
	ContentType string // tipo do arquivo
	URL         string // opcional: se já tiver URL, não faz upload
}

func (uc *AddServicePhotoUseCase) Execute(ctx context.Context, input AddServicePhotoInput) (*provider.Provider, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "AddServicePhotoUseCase",
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.ServiceID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "AddServicePhotoUseCase",
			Message: "ServiceID é obrigatório",
			Error:   fmt.Errorf("serviceID é obrigatório"),
		})
		return nil, []exceptions.ProblemDetails{
			{
				Type:   exceptions.RFC400,
				Title:  "ServiceID é obrigatório",
				Status: exceptions.RFC400_CODE,
				Detail: "O ID do serviço é obrigatório.",
			},
		}
	}

	// Buscar o provider através do ID do serviço
	providerEntity, err := uc.providerRepo.FindProviderByServiceID(ctx, input.ServiceID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    "AddServicePhotoUseCase",
			Message: "Serviço não encontrado",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{
			{
				Type:   exceptions.RFC404,
				Title:  "Serviço não encontrado",
				Status: exceptions.RFC404_CODE,
				Detail: "O serviço informado não foi encontrado.",
			},
		}
	}

	var photoURL string
	if len(input.FileData) > 0 && input.FileName != "" && input.ContentType != "" {
		// Verificar se o serviço Minio está disponível
		if uc.minio == nil {
			uc.logger.Log(logging.Logger{
				Context: ctx,
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.USECASES,
				Code:    exceptions.RFC503_CODE,
				From:    "AddServicePhotoUseCase",
				Message: "Serviço de upload não está disponível",
				Error:   fmt.Errorf("minio service not configured"),
			})
			return nil, []exceptions.ProblemDetails{
				{
					Type:   exceptions.RFC503,
					Title:  "Serviço de upload não disponível",
					Status: exceptions.RFC503_CODE,
					Detail: "O serviço de upload de arquivos não está configurado. Use uma URL externa ou configure o serviço Minio.",
				},
			}
		}

		photoURL, err = uc.minio.Upload(ctx, input.FileName, input.FileData, input.ContentType)
		if err != nil {
			uc.logger.Log(logging.Logger{
				Context: ctx,
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.USECASES,
				Code:    exceptions.RFC500_CODE,
				From:    "AddServicePhotoUseCase",
				Message: "Erro ao fazer upload para Minio",
				Error:   err,
			})
			return nil, []exceptions.ProblemDetails{
				{
					Type:   exceptions.RFC500,
					Title:  "Erro ao fazer upload para Minio",
					Status: exceptions.RFC500_CODE,
					Detail: err.Error(),
				},
			}
		}
	} else if input.URL != "" {
		photoURL = input.URL
	} else {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "AddServicePhotoUseCase",
			Message: "Arquivo ou URL da foto obrigatórios",
			Error:   fmt.Errorf("arquivo ou URL da foto obrigatórios"),
		})
		return nil, []exceptions.ProblemDetails{
			{
				Type:   exceptions.RFC400,
				Title:  "Arquivo ou URL da foto obrigatórios",
				Status: exceptions.RFC400_CODE,
				Detail: "É necessário enviar o arquivo da foto ou uma URL válida.",
			},
		}
	}

	// Adicionar foto ao serviço específico
	if err := providerEntity.AddServicePhoto(input.ServiceID, photoURL); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "AddServicePhotoUseCase",
			Message: "Erro ao adicionar foto ao serviço",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{
			{
				Type:   exceptions.RFC400,
				Title:  "Erro ao adicionar foto ao serviço",
				Status: exceptions.RFC400_CODE,
				Detail: err.Error(),
			},
		}
	}

	if err := uc.providerRepo.Update(ctx, providerEntity); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    "AddServicePhotoUseCase",
			Message: "Erro ao persistir serviço com foto",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{
			{
				Type:   exceptions.RFC500,
				Title:  "Erro ao persistir serviço com foto",
				Status: exceptions.RFC500_CODE,
				Detail: err.Error(),
			},
		}
	}

	return providerEntity, nil
}
