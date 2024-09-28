package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kTowkA/GophKeeper/client"
	"github.com/kTowkA/GophKeeper/client/config"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
)

func main() {
	// При указании флага линковщика -ldflags и указании флага -X можно собрать приложение с определенными значениями buildVersion,buildDate,buildCommit
	// -X main.buildVersion=версия
	fmt.Println(buildVersion)
	// -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' для получения даты сборки
	fmt.Println(buildDate)

	cfg := config.LoadConfig()

	err := client.Run(context.Background(), *cfg)
	if err != nil {
		log.Fatal(err)
	}
}
