package main

import (
	"log"

	"github.com/idalmasso/storytellers/backend/api"
	"github.com/idalmasso/storytellers/backend/database/mongodb"
	"github.com/idalmasso/storytellers/backend/service"
	"github.com/joho/godotenv"
)


func main(){
	err := godotenv.Load(".env")
	 if err != nil {
    log.Fatalf("Error loading .env file")
  }
	db:=mongodb.MongoDatabase{}
	app,err:=service.CreateApp(&db)
	if err!=nil{
		log.Fatal("Error on service create: ",err)
	}
	srv :=  api.NewStoryServer(&app)
	srv.Run()
}
