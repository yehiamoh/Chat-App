package routes

import (
	"chat-app/pkg/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"gorm.io/gorm"
)


var store *sessions.CookieStore

func Init(){
	if err:=godotenv.Load();err!=nil{
		log.Fatal("Error in loading Enviroment variable")
	}
	store=sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.Options=&sessions.Options{
		Path: "/",
		MaxAge:   86400 * 30, // 30 days
		HttpOnly: true,
	}
	gothic.Store=store

	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_AUTH_CLIENT_ID"),
			os.Getenv("GOOGLE_AUTH_CLIENT_SECRET"),
			"http://localhost:8000/auth/callback", // Adjust this URL based on your setup
			"email", "profile",
		),
	)
}

func(r *Router)AuthMiddleWare()gin.HandlerFunc{
	return func(c *gin.Context) {
		session,err:=store.Get(c.Request,"user-session")
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error": "Failed to get session"})
			c.Abort()
			return
		}
		if auth,ok:=session.Values["authenticated"].(bool);!ok ||!auth{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}



func (r *Router) BeginAuthRoute(c *gin.Context){
	gothic.BeginAuthHandler(c.Writer,c.Request)
}




func (r *Router) CallBackAuthRoute(c *gin.Context){
	user,err:=gothic.CompleteUserAuth(c.Writer,c.Request)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":fmt.Sprintf("Error Completing auth %v",err.Error()),
		})
		return
	}

	var dbUser models.User

	result:=r.db.Where("google_id = ?",user.UserID).First(&dbUser)

	if result.Error!=nil{
		if result.Error== gorm.ErrRecordNotFound{
			dbUser=models.User{
				GoogleID: user.UserID,
				Email:    user.Email,
				Name:     user.Name,
				Picture:  user.AvatarURL,
			}
			if err:=r.db.Create(&dbUser).Error;err!=nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}
		}
	}else{
		dbUser.Email = user.Email
		dbUser.Name = user.Name
		dbUser.Picture = user.AvatarURL
		dbUser.LastLoginAt = time.Now()
		if err := r.db.Save(&dbUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
	}
	session,err:=store.Get(c.Request,"user-session")
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get session"})
		return
	}

	session.Values["authenticated"]=true
	session.Values["user_id"]=dbUser.GoogleID
	session.Values["email"]=dbUser.Email
	session.Values["name"]=dbUser.Name

	if err := session.Save(c.Request, c.Writer);err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	// Return user information
	c.JSON(http.StatusOK, gin.H{
		"user": dbUser,
	})
}


func (r *Router) Logout(c *gin.Context){
	session,err:=store.Get(c.Request,"user-session")
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error": "Failed to get session",
		})
	}
	session.Values["authenticated"]=false
	session.Options.MaxAge=-1

	if err:=session.Save(c.Request,c.Writer);err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error": "Failed to save session",
		})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}




func (r *Router)GetUserFromSession(c *gin.Context){
	session,err:=store.Get(c.Request,"user-session")
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get session"})
		return
	}

	if auth,ok:=session.Values["authenticated"].(bool);!ok||!auth{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	userData := gin.H{
		"user_id": session.Values["user_id"],
		"email":   session.Values["email"],
		"name":    session.Values["name"],
	}

	c.JSON(http.StatusOK,userData)
}