package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Note: This is just an example for a tutorial

func main() {

	router := gin.Default()

	store := sessions.NewCookieStore([]byte("sessionSuperSecret"))
	router.Use(sessions.Sessions("sessionName", store))

	api := router.Group("/api/v1")

	// no authentication endpoints
	api.POST("/login", loginHandler)
	api.GET("/message/:msg", noAuthMessageHandler)

	//user authentication endpoints
	userAuth := api.Group("/user")
	userAuth.Use(AuthenticationRequired("user"))
	userAuth.GET("/message/:msg", userMessageHandler)

	// admin authentication endpoints
	adminAuth := api.Group("/admin")
	adminAuth.Use(AuthenticationRequired("admin"))
	adminAuth.GET("/message/:msg", adminMessageHandler)

	// subscriber authentication endpoints
	subscriberAuth := api.Group("/subscriber")
	subscriberAuth.Use(AuthenticationRequired("subscriber"))
	subscriberAuth.GET("/message/:msg", subscriberMessageHandler)

	logout := api.Group("/")
	logout.Use(AuthenticationRequired())
	logout.GET("/logout", logoutHandler)

	router.Use(AuthenticationRequired("admin", "subscriber"))
	{
		api.POST("/post/message/:msg", postMessageHandler)
	}

	router.Run(":8080")
}
