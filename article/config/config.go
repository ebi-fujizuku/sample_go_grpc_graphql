package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)


type configStruct struct{
	Port         int
	Sqlite3_path string
}
var Conf configStruct
func NewConfig(envFilepath string){
	err := godotenv.Load(envFilepath)
	if err != nil {
		panic(errors.New(".envが読み込み出来ませんでした"))
	}

	Conf = configStruct{
		Port:	8080,
		Sqlite3_path: os.Getenv("SQLITE3_PATH"),
	}
	fmt.Println("Conf",Conf)
}
