package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"openai-proxy/jwt"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	SecretKey := os.Getenv("SECRET_KEY")
	jwt.Secret = []byte(SecretKey)

	log.Println("Please enter your username:")
	var line string
	_, err = fmt.Scanln(&line)
	if err != nil {
		return
	}
	log.Printf("Your input username is: %s \n", line)
	jwt.CreateJwt(line)
}
