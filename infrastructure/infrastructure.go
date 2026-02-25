package infrastructure

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/VitorFranciscoDev/sprinter-api/domain/entities"
	"github.com/VitorFranciscoDev/sprinter-api/domain/usecases"
	"github.com/VitorFranciscoDev/sprinter-api/infrastructure/datastore/repositories"
	"github.com/VitorFranciscoDev/sprinter-api/infrastructure/router"
	"github.com/VitorFranciscoDev/sprinter-api/infrastructure/router/modules"
	"github.com/gorilla/mux"
)

func SetupModules(r *mux.Router, config entities.Config) error {
	// Repository settings
	settings, err := repositories.NewSettingsRepository(config)
	if err != nil {
		return errors.Join(errors.New("failed to create settings repository"), err)
	}

	// Repositories
	authRepository := repositories.NewAuthenticationRepository(settings)

	// Use Cases
	authUseCases := usecases.NewAuthenticationUseCase(authRepository, config.Paseto.SecurityKey)

	// Modules
	authModule := modules.NewAuthModule(*authUseCases)

	apiSubRouter := r.PathPrefix("/api").Subrouter()

	_, _ = authModule.Setup(apiSubRouter)

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
