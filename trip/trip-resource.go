package main
import "github.com/gin-gonic/gin"




func init() {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use()


}