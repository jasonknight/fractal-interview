package main
import "github.com/gin-gonic/gin"
import "net/http"
func main() {
	conf,err := LoadConfig()
	if err != nil {
		panic(err)
	}
	r := createRouter()
	r.GET("/tree/:name",func (c *gin.Context ) {
		name := c.Params.ByName("name")	
		c.JSON(http.StatusOK,gin.H{"name": name})
	})
	r.Run(conf.Listen)
}
