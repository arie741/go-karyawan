package main

import (
	"context"
	"html/template"
	"log"

	"github.com/arie741/go-karyawan/internal/auth"
	"github.com/arie741/go-karyawan/internal/database"
	"github.com/arie741/go-karyawan/internal/formatter"
	"github.com/arie741/go-karyawan/internal/karyawan"
	"github.com/arie741/go-karyawan/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	mongoClient := database.Connect()
	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			log.Println(err)
		}
	}()

	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"formatRupiah": formatter.FormatRupiah,
	})
	router.LoadHTMLGlob("templates/**/*")

	// LOGIN
	authHandler := auth.NewHandler(auth.NewRepository(database.GetCollection(mongoClient, "accounts")))
	router.GET("/login", authHandler.LoginPage)
	router.POST("/login", authHandler.Login)
	router.POST("/logout", authHandler.Logout)

	// KARYAWAN
	karyawanHandler := karyawan.NewHandler(karyawan.NewRepository(database.GetCollection(mongoClient, "karyawan")))
	authorized := router.Group("/", middleware.Authorization())

	authorized.GET("/", karyawanHandler.List)
	authorized.GET("/karyawan/:id", karyawanHandler.FindById)
	authorized.POST("/add-karyawan", karyawanHandler.Add)
	authorized.POST("/edit/:id", karyawanHandler.Edit)
	authorized.DELETE("/karyawan/:id", karyawanHandler.Delete)
	authorized.GET("/modal", karyawanHandler.Modal)
	router.Run(":3000")
}
