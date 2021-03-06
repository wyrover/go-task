package main

import (
	"github.com/itsjamie/gin-cors"
	"github.com/rorikurniadi/go-task/models"
	"github.com/rorikurniadi/go-task/resources"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := models.InitDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	statusResource := resources.NewStatusStorage(db)
	authResource := resources.AuthDB(db)
	taskResource := resources.TaskDB(db)
	tagResource := resources.TagDB(db)
	noteResource := resources.NoteDB(db)

	r := gin.Default()

	// handle CORS
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		Credentials:     true,
		ValidateHeaders: false,
	}))

	v1 := r.Group("api/v1")
	{
		v1.POST("/register", authResource.Register)
		v1.POST("/login", authResource.Login().LoginHandler)
		v1.GET("/logout", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Logout",
			})
		})

	}

	auth := r.Group("api/v1")
	auth.Use(authResource.Login().MiddlewareFunc())
	{
		auth.GET("/refresh_token", authResource.Login().RefreshHandler)

		//
		auth.GET("/users/current", authResource.CurrentUser)
		auth.GET("/users", authResource.Get)

		// task
		auth.GET("/tasks", taskResource.Get)
		auth.GET("/tasks/:id", taskResource.Show)
		auth.POST("/tasks", taskResource.Store)
		auth.PUT("/tasks/:id", taskResource.Update)
		auth.DELETE("/tasks/:id", taskResource.Destroy)
		auth.GET("/tasks/:id/notes", noteResource.GetByTask)

		// tag
		auth.GET("/tags", tagResource.Get)
		auth.GET("/tags/:id", tagResource.Show)
		auth.POST("/tags", tagResource.Store)
		auth.PUT("/tags/:id", tagResource.Update)
		auth.DELETE("/tags/:id", tagResource.Destroy)

		// note
		auth.GET("/notes", noteResource.Get)
		auth.GET("/notes/:id", noteResource.Show)
		auth.POST("/notes", noteResource.Store)
		auth.PUT("/notes/:id", noteResource.Update)
		auth.DELETE("/notes/:id", noteResource.Destroy)

		// statuses
		auth.GET("/statuses", statusResource.Get)
		auth.GET("/statuses/:id", statusResource.Show)
		auth.POST("/statuses", statusResource.Store)
	}

	r.Run(":8080")
}
