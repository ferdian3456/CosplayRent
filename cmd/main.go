package main

import (
	"context"
	"cosplayrent/internal/config"
	"cosplayrent/internal/exception"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "True")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	router := config.NewRouter()
	zerolog := config.NewZeroLog()
	koanf := config.NewKoanf()
	db := config.NewDB(koanf, &zerolog)
	memcacheClient := config.NewMemcacheClient(koanf)
	validator := config.NewValidator()

	config.Server(&config.ServerConfig{
		Router:   router,
		DB:       db,
		Memcache: memcacheClient,
		Log:      &zerolog,
		Validate: validator,
		Config:   koanf,
	})

	router.ServeFiles("/static/*filepath", http.Dir("../static"))
	router.PanicHandler = exception.ErrorHandler

	GO_SERVER_PORT := koanf.String("GO_SERVER")

	server := http.Server{
		Addr:    GO_SERVER_PORT,
		Handler: CORS(router),
	}

	zerolog.Info().Msg(("Server is running on " + GO_SERVER_PORT))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zerolog.Fatal().Err(err).Msg("Error Starting Server")
		}
	}()

	<-stop
	zerolog.Info().Msg("Got one of stop signals")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		zerolog.Fatal().Err(err).Msg("Timeout, forced kill!")
	}

	zerolog.Info().Msg("Server has shut down gracefully")
}
