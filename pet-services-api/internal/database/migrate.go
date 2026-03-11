package database

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type SchemaMigration struct {
	Version   string    `gorm:"primaryKey"`
	AppliedAt time.Time `gorm:"not null"`
}

func RunMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&SchemaMigration{}); err != nil {
		return fmt.Errorf("erro ao criar tabela de migrações: %w", err)
	}

	migrations := getMigrations()

	for _, migration := range migrations {
		var existing SchemaMigration
		result := db.Where("version = ?", migration.Version).First(&existing)

		if result.Error == nil {
			continue
		}

		if result.Error != gorm.ErrRecordNotFound {
			return fmt.Errorf("erro ao verificar migração %s: %w", migration.Version, result.Error)
		}

		slog.Info("[Migration] aplicando migração", "version", migration.Version, "description", migration.Description)

		err := db.Transaction(func(tx *gorm.DB) error {
			if err := migration.Up(tx); err != nil {
				return fmt.Errorf("erro ao executar migração: %w", err)
			}

			record := SchemaMigration{
				Version:   migration.Version,
				AppliedAt: time.Now(),
			}
			if err := tx.Create(&record).Error; err != nil {
				return fmt.Errorf("erro ao registrar migração: %w", err)
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("falha na migração %s: %w", migration.Version, err)
		}

		slog.Info("[Migration] migração aplicada com sucesso", "version", migration.Version)
	}

	slog.Info("[Migration] todas as migrações foram aplicadas")
	return nil
}

type Migration struct {
	Version     string
	Description string
	Up          func(*gorm.DB) error
}

func getMigrations() []Migration {
	return []Migration{
		{
			Version:     "20260110000000",
			Description: "Criar esquema inicial para pet-services",
			Up:          Migration20260110000000,
		},
		{
			Version:     "20260215000000",
			Description: "Criar tabela refresh_tokens",
			Up:          Migration20260215000000,
		},
		{
			Version:     "20260204000000",
			Description: "Criar tabela password_reset_tokens para redefinição de senha e verificação de email",
			Up:          Migration20260204000000,
		},
		{
			Version:     "20260213000000",
			Description: "Adicionar profile_complete ao usuário",
			Up:          Migration20260213000000,
		},
		{
			Version:     "20260218000000",
			Description: "Seed inicial da tabela species",
			Up:          Migration20260218000000,
		},
		{
			Version:     "20260218000001",
			Description: "Seed inicial da tabela categories",
			Up:          Migration20260218000001,
		},
		{
			Version:     "20260311000000",
			Description: "Seed inicial de usuários owner/provider e provider vinculado",
			Up:          Migration20260311000000,
		},
		{
			Version:     "20260311000001",
			Description: "Seed inicial de pets para usuário owner",
			Up:          Migration20260311000001,
		},
	}
}
