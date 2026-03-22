package infrastructure

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/usecases"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore/repositories"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/filestore/hdstore"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router/modules"
	"github.com/gorilla/mux"
)

func SetupModules(r *mux.Router, config entities.Settings) error {
	// Repository settings
	settings, err := repositories.NewSettingsRepository(config)
	if err != nil {
		return errors.Join(errors.New("failed to create settings repository"), err)
	}

	fileStorage := hdstore.NewHDFileStorage(config)

	// Repositories
	authRepository := repositories.NewAuthenticationRepository(settings)
	productRepository := repositories.NewProductRepository(settings)

	// Use Cases
	authUseCases := usecases.NewAuthenticationUseCase(authRepository, config.Paseto.SecurityKey)
	_ = usecases.NewUserUseCases(fileStorage, config.FileStorage.StorageFolder)
	productUseCases := usecases.NewProductUseCases(productRepository)

	// Modules
	authModule := modules.NewAuthModule(authUseCases)
	productModule := modules.NewProductModule(productUseCases)

	apiSubRouter := r.PathPrefix("/api").Subrouter()

	_, _ = authModule.Setup(apiSubRouter)
	_, _ = productModule.Setup(apiSubRouter)

	r.Use(router.LoggingMiddleware)

	// Home URL handler returns the current server time
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		serverTime, err := settings.ServerTime(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = fmt.Fprintf(w, "%v", serverTime.UTC().Unix())

		if err != nil {
			log.Println(err)
		}
	})

	return nil
}
