package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"

	"go.uber.org/zap"

	"github.com/senago/technopark-dbms/internal/api"
)

const (
	configPathEnvVar = "CONFIG_PATH"
	defaultAddress   = "0.0.0.0"
	defaultPort      = "8080"
)

func main() {
	// -------------------- Set up viper -------------------- //

	viper.AutomaticEnv()

	viper.SetConfigFile("/cmd/configs/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file: %s\n", err)
	}

	viper.SetDefault("service.bind.address", defaultAddress)
	viper.SetDefault("service.bind.port", defaultPort)

	// -------------------- Set up logging -------------------- //

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to set up logger: %s\n", err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	// -------------------- Set up database -------------------- //

	dbPool, err := pgxpool.Connect(context.Background(), viper.GetString("db.connection_string"))
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err)
	}
	defer dbPool.Close()

	// -------------------- Set up service -------------------- //

	svc, err := api.NewAPIService(sugar, dbPool)
	if err != nil {
		log.Fatalf("error creating service instance: %s", err)
	}

	go svc.Serve(viper.GetString("service.bind.address") + ":" + viper.GetString("service.bind.port"))

	// -------------------- Listen for INT signal -------------------- //

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(viper.GetInt("service.shutdown_timeout")),
	)
	defer cancel()

	if err := svc.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
