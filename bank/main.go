package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/figormartins/ms-it/infrastructure/repository"
	"github.com/figormartins/ms-it/usecase"
	_ "github.com/lib/pq"
)

func main() {
	db := setupDb()
	defer db.Close()
}

func setupTransactionUseCase(db *sql.DB) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepository(db)
	useCase := usecase.NewUseCaseTransaction(transactionRepository)
	return useCase
}

func setupDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"db", "5432", "postgres", "root", "codebank")

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("error connecting database")
	}

	return db
}
