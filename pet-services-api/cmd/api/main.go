// @title Pet Services API
// @version 1.0
// @description API para gerenciamento de serviços pet.
// @termsOfService http://swagger.io/terms/
// @contact.name Guilherme Amorim
// @contact.email guilherme@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT Authorization header usando o esquema Bearer. Exemplo: 'Authorization: Bearer {token}'
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pet-services-api/internal/config"
	"pet-services-api/internal/database"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/routes"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	logger := &logging.DefaultLogger{}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, sqlDB := database.SetupDatabaseConnection(ctx)
	if db == nil {
		slog.Error("[Start] Banco de dados não disponível, encerrando aplicação")
		os.Exit(1)
	}
	defer database.Shutdown(ctx, db)
	defer sqlDB.Close()

	if err := database.RunMigrations(db); err != nil {
		slog.Error("[Start] Falha ao rodar migrações do banco", "error", err)
		os.Exit(1)
	}
	slog.Info("[Start] Migrações do banco concluídas com sucesso")

	storageInput := database.StorageInput{
		DB:         db,
		BucketName: "",
	}

	router := routes.SetupRouter(storageInput, ctx, logger)

	server := &http.Server{
		Addr:    config.GetServerPort(),
		Handler: router,
	}

	go func() {
		slog.Info("[Start] Servidor HTTP iniciado", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("[Start] Falha ao iniciar servidor HTTP", "error", err)
			os.Exit(1)
		}
	}()

	slog.Info("[Start] Serviço rodando. Pressione Ctrl+C para encerrar.")

	<-ctx.Done()
	slog.Info("[Start] Encerrando aplicação com shutdown gracioso...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("[Shutdown] Erro ao encerrar servidor HTTP", "error", err)
	} else {
		slog.Info("[Shutdown] Servidor HTTP encerrado com sucesso")
	}

	slog.Info("[Shutdown] Aplicação finalizada")
}
