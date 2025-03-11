package dsn

import (
	"fmt"
	"os"
)

func FromEnv() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		fmt.Println("No host")
		return ""
	}
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		fmt.Println("Password is empty!")
	} else {
		fmt.Println("Password loaded successfully!")
	}

	dbname := os.Getenv("DB_NAME")
	// И вот мы возвращаем dsn, который необходим для подключения к БД
	// fmt.Printf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}
