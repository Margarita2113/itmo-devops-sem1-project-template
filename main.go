package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5"

	"project_sem/internal/postgres"
	"project_sem/internal/server"
)

func main() {
	fmt.Println("Start application")

	sigs := make(chan os.Signal, 1)
	// Указываем, что нас интересуют SIGINT, SIGSTOP и SIGTERM
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP)

	db, err := postgres.NewDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	err = server.NewServer(db)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Server started")
	<-sigs
	fmt.Println("application stopped")
}
