package backend

import (
	"github.com/gin-gonic/gin"
	srv "github.com/hhandhuan/ku-bbs/internal/service"
)

var Home = cHome{}

type cHome struct{}

func (c *cHome) IndexPage(ctx *gin.Context) {
	srv.Context(ctx).View("backend.home.index", gin.H{})
}
