package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lvyonghuan/Ubik-Util/ujson"
	"github.com/lvyonghuan/Ubik-Util/ulog"
)

func response(status int, info any) gin.H {
	return gin.H{
		"status": status,
		"info":   info,
	}
}

func successResponse(c *gin.Context, info any) {
	marshalInfo, _ := ujson.Marshal(info)

	c.JSON(http.StatusOK, response(200, marshalInfo))
}

func errorResponse(c *gin.Context, status int, info any) {
	marshalInfo, _ := ujson.Marshal(info)

	c.JSON(http.StatusOK, response(status, marshalInfo))
}

// fatalErrHandel handles fatal errors by logging them and returning a 500 response.
func fatalErrHandel(c *gin.Context, err error) {
	l := ulog.NewLogWithoutPost(ulog.Debug, true, "./logs")
	l.Fatal(err)
	marshalInfo, _ := ujson.Marshal("Internal Server Error: " + err.Error())
	c.JSON(http.StatusOK, response(500, marshalInfo))
	os.Exit(1)
}
