package configs

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

func GetMongoURI() string {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Failed or maybe erorr on .env file")
    }
    return os.Getenv("MONGO_URI")
}
