package database

import (
	"errors"
	"sort"
	"time"

	"gorm.io/gorm"
)

// SchemaMigration armazena versões já aplicadas.
type SchemaMigration struct {
	Version   string    `gorm:"primaryKey;size:14"`
	AppliedAt time.Time `gorm:"not null"`
}

// Migration descreve uma alteração de schema.
type Migration struct {
	Version     string
	Description string
	Up          func(*gorm.DB) error
}

// RunMigrations executa migrações pendentes em ordem.
func RunMigrations(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	// Garante tabela de controle.
	if err := db.AutoMigrate(&SchemaMigration{}); err != nil {
		return err
	}

	migrations := getMigrations()
	sort.SliceStable(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	applied := map[string]bool{}
	var rows []SchemaMigration
	if err := db.Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		applied[row.Version] = true
	}

	for _, m := range migrations {
		if applied[m.Version] {
			continue
		}

		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := m.Up(tx); err != nil {
				return err
			}
			return tx.Create(&SchemaMigration{
				Version:   m.Version,
				AppliedAt: time.Now(),
			}).Error
		}); err != nil {
			return err
		}
	}

	return nil
}
