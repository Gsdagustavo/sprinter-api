package infrastructure

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	usecases "github.com/Gsdagustavo/sprinter-api/domain/usecases/impl"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore/repositories/impl"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/filestore/hdstore"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/mail"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router/logger"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router/modules"
	"github.com/gorilla/mux"
)

func SetupModules(r *mux.Router, settings entities.Settings) error {
	// Repository settings
	settingsRepository, err := repositories.NewSettingsRepository(settings)
	if err != nil {
		return derr.JoinError("failed to create settings repository", err)
	}

	fileStorage := hdstore.NewHDFileStorage(settings)

	mailSender, err := mail.NewMailSender(settings)
	if err != nil {
		return derr.JoinError("failed to create mail sender", err)
	}

	// Repositories
	authRepository := repositories.NewAuthenticationRepository(settingsRepository)
	productRepository := repositories.NewProductRepository(settingsRepository)
	userRepository := repositories.NewUserRepository(settingsRepository)

	// Use Cases
	authUseCases := usecases.NewAuthenticationUseCase(authRepository, settings.PasetoSettings.SecurityKey, fileStorage, mailSender)
	userUseCases := usecases.NewUserUseCases(fileStorage, settings.FileStorageSettings, userRepository)
	productUseCases := usecases.NewProductUseCases(productRepository)

	// Modules
	authModule := modules.NewAuthModule(authUseCases)
	productModule := modules.NewProductModule(productUseCases)
	userModule := modules.NewUserModule(userUseCases)

	modules := []router.Module{
		authModule,
		productModule,
		userModule,
	}

	apiSubRouter := r.PathPrefix("/api").Subrouter()

	for _, module := range modules {
		moduleRouter := apiSubRouter.PathPrefix(module.Path()).Subrouter()

		protectedRouter := moduleRouter.NewRoute().Subrouter()
		for _, mw := range module.Middlewares() {
			protectedRouter.Use(mw)
		}

		for _, d := range module.Routes() {
			target := protectedRouter
			if d.Public {
				target = moduleRouter
			}
			target.HandleFunc(d.Path, d.Handler).Methods(d.HttpMethods...)
		}
	}

	r.Use(router.LoggingMiddleware)

	// Home URL handler returns the current server time
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		serverTime, err := settingsRepository.ServerTime(ctx)
		if err != nil {
			router.HandleError(w, err)
			return
		}

		_, err = fmt.Fprintf(w, "%v", serverTime.UTC().Unix())
		if err != nil {
			slog.ErrorContext(ctx, "failed to respond time", logger.Err(err))
		}
	})

	return nil
}
