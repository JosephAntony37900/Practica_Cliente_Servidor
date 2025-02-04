package routes

import (
	"github.com/gin-gonic/gin"
	"PracticaClases/controllers"
)

func SetupProductosRoutes(router *gin.Engine) {
	// Rutas para el servidor principal
	router.POST("/productos", controllers.CrearProducto)
	router.GET("/productos", controllers.ObtenerProductos)
	router.GET("/productos/:id", controllers.ObtenerProductoPorID)
	router.PUT("/productos/:id", controllers.ActualizarProducto)
	router.DELETE("/productos/:id", controllers.EliminarProducto)
	//rutas para el servidor de replicación
	router.GET("/productos/short-pulling", controllers.ShortPulling)
	router.GET("/productos/long-pulling", controllers.LongPulling)
	router.GET("/productos/incremental", controllers.ObtenerProductosIncremental)
}