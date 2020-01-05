package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

type User struct {
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	UserType string `form:"userType" json:"userType" xml:"userType" binding:"required"`
}

var (
	VALID_AUTHENTICATIONS = []string{"user", "admin", "subscriber"}
)

func noAuthMessageHandler(c *gin.Context) {
	msg := c.Param("msg")
	c.JSON(http.StatusOK, gin.H{"Your Message": msg})
}

//user Login
func loginHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)

	if strings.Trim(user.Username, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username can't be empty"})
	}
	if !funk.ContainsString(VALID_AUTHENTICATIONS, user.UserType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid auth type"})
	}

	// Note: This is just an example, in real world AuthType would be set by business logic and not the user
	session.Set("user", user.Username)
	session.Set("userType", user.UserType)

	err := session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate session token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "authentication successful"})
}

//user Logout
func logoutHandler(c *gin.Context) {
	session := sessions.Default(c)

	// this would only be hit if the user was authenticated
	session.Delete("user")
	session.Delete("authType")

	err := session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate session token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})

}

//post message (for Admin and Subscriber)
func postMessageHandler(c *gin.Context) {
	msg := c.Param("msg")

	session := sessions.Default(c)
	user := session.Get("user")

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello %s {User}, your message: %s will be posted", user, msg)})
}

//user Message
func userMessageHandler(c *gin.Context) {
	msg := c.Param("msg")

	session := sessions.Default(c)
	user := session.Get("user")
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello %s {User}, your message: %s", user, msg)})
}

//subscriber Message
func subscriberMessageHandler(c *gin.Context) {
	msg := c.Param("msg")

	session := sessions.Default(c)
	user := session.Get("user")

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello %s {Subscriber}, here's your message: %s", user, msg)})
}

//admin Message
func adminMessageHandler(c *gin.Context) {
	msg := c.Param("msg")

	session := sessions.Default(c)
	user := session.Get("user")

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello %s {Admin}, here's your message: %s", user, msg)})
}
