package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/aldotp/ecommerce-go-api/pkg/consts"
	"github.com/gin-gonic/gin"
)

func Index(g *gin.Engine, version, name string) {
	g.GET("/", func(context *gin.Context) {
		context.JSON(200, struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}{
			Name:    name,
			Version: version,
		})
	})
}

func CatchInternalServerError(errMessage error) error {
	if errors.Is(errMessage, consts.ErrInternal) {
		return consts.ErrInternal
	}

	fileName := fmt.Sprintf("./storage/error/error-%s.log", time.Now().Format("2006-01-02"))

	// open log file
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}

	defer logFile.Close()

	// set log out put
	log.SetOutput(logFile)

	log.SetFlags(log.LstdFlags)

	_, fileName, line, _ := runtime.Caller(1)
	log.Printf("[Error] in [%s:%d] %v", fileName, line, errMessage.Error())

	return consts.ErrInternal
}
