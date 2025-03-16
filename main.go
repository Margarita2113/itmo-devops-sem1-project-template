package main

import (
	"fmt"
	_ "github.com/jackc/pgx/v5"
	"os"
	"os/signal"
	"project_sem/internal/postgres"
	"project_sem/internal/server"
	"syscall"
)

func main() {
	fmt.Println("Start application")
	db, err := postgres.NewDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.Migrate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Database migration complete")

	err = server.NewServer(db)
	if err != nil {
		fmt.Println(err)
	}

	sigs := make(chan os.Signal, 1)
	// Указываем, что нас интересуют SIGINT и SIGTERM
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP)

	<-sigs
	fmt.Println("application stopped")
}
