package main

import "github.com/gin-gonic/gin"
import "net/http"
import "io/ioutil"
import "strconv"

func fetch(conf Configuration, name string) ([]byte, int, error) {
	response, err := http.Get(conf.RemoteUrl + "/" + name)
	if err != nil {
		return []byte(""), 0, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte(""), 0, err
	}
	return contents, response.StatusCode, nil
}
func treeHandler(conf Configuration, c *gin.Context) {
	name := c.Params.ByName("name")
	indicator_ids, ok := c.GetQueryArray("indicator_ids[]")
	if ok != true {
		c.JSON(418, gin.H{"error": "You are a teapot and did not send indicator_ids in the query string", "name": name})
		return
	}
	retries := conf.Retries
	var errors []error
	for retries > 0 {
		contents, status, err := fetch(conf, name)
		if err != nil {
			errors = append(errors, err)
		}
		if status == 200 {
			collection, err := ParseTree(contents)
			if err != nil {
				c.JSON(500, gin.H{"error": err})
				return
			}
			filter := func(ind Indicator) bool {
				for i := range indicator_ids {
					iid, _ := strconv.Atoi(indicator_ids[i])
					if iid == ind.Id {
						return true
					}
				}
				return false
			}
			new_collection := collection.FilterByIndicators(filter)
			new_themes := new_collection.Themes
			c.IndentedJSON(200, new_themes)
			return
		}
		retries--
	}
	c.JSON(500, errors)
}
func main() {
	conf, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	r := createRouter()
	r.GET("/tree/:name", func(c *gin.Context) {
		treeHandler(conf, c)
	})
	r.Run(conf.Listen)
}
