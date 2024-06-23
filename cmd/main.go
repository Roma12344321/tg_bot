package main

import (
	"context"
	"github.com/spf13/viper"
	"log"
	"myBot/pkg/handler"
	"myBot/pkg/service"
	"myBot/pkg/storage"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	initConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	chanOs := make(chan os.Signal)
	signal.Notify(chanOs, syscall.SIGINT, syscall.SIGTERM)
	m := &sync.Map{}
	storages := storage.NewStorage(m)
	services := service.NewService(storages, ctx)
	handlers := handler.NewHandler(services)
	go handlers.InitBot(ctx)
	<-chanOs
	cancel()
}

func initConfig() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}
