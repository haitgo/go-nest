package httplib

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSession(t *testing.T) {
	g := gin.Default()
	g.Use(func(c *gin.Context) {
		d := Session{
			MaxAge:  3600,
			Writer:  c.Writer,
			Request: c.Request,
			RdsHost: "localhost:6379",
			RdsDB:   1,
		}.Use()
		d.Set("name", "李老汉")
	})
	g.GET("/", func(c *gin.Context) {
		html := `<html>
				<body>
				<form method="POST" action="/" >
				<input type="text" name="Name" value="" />
				<input type="text" name="Age" value="" /> 
				<input type="submit" name="Name" value="提交" />
				</form>
				</body>
				</html>`
		c.Header("content-type", "text/html;charset=utf-8")
		c.String(200, html)
	})
	g.POST("/", func(c *gin.Context) {
		d := Session{}.Get(c.Request)
		c.String(200, d.GetString("name"))
	})
	g.Run(":9923")
}
