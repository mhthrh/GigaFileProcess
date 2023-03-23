package Server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mhthrh/GigaFileProcess/Entity"
	"github.com/mhthrh/GigaFileProcess/Rabbit"
	"github.com/mhthrh/GigaFileProcess/Validation"
	"net/http"
	"time"
)

func Run(ctx *gin.Context) {
	var request Entity.RunRequest
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, "cannot pars request object")
		ctx.Abort()
		return
	}
	if err := Validation.Priority(request.Priority); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		ctx.Abort()
		return
	}
	if err := Validation.Count(request.Count); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		ctx.Abort()
		return
	}
	if request.FileName == "" {
		ctx.JSON(http.StatusBadRequest, "file name be empty")
		ctx.Abort()
		return
	}
	count := request.Count / request.Priority
	rabbit, err := Rabbit.New("")
	if err != nil {
		if request.FileName == "" {
			ctx.JSON(http.StatusInternalServerError, "cannot connect to rabbit MQ")
			ctx.Abort()
			return
		}
	}
	for i := 0; i < count; i++ {
		if err := rabbit.DeclareQueue(fmt.Sprintf("%s-%d", request.ID.String(), i)); err != nil {
			ctx.JSON(http.StatusInternalServerError, "cannot create queue rabbit MQ")
			ctx.Abort()
			return
		}
	}

}
func Version(context *gin.Context) {
	context.JSON(http.StatusOK, "Ver:1.0.0")
}
func NotFound(context *gin.Context) {
	context.JSON(http.StatusOK, struct {
		Time        time.Time `json:"time"`
		Description string    `json:"description"`
	}{
		Time:        time.Now(),
		Description: "Workers are working, coming soon!!!",
	})
}
