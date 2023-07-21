package main

import (
	"fmt"
	"openai-proxy/jwt"
)

func main() {
	println("Please enter your username:")
	var line string
	_, err := fmt.Scanln(&line)
	if err != nil {
		return
	}
	fmt.Printf("Your input username is: %s \n", line)
	jwt.CreateJwt(line)
}
