package main

import (
	"github.com/joho/godotenv"
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Fatal("Error loading .env.test file")
	}
}
