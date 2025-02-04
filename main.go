package main

import (
	"github.com/gin-gonic/gin"
	"PracticaClases/routes"
)

func main() {
	routerPrincipal := gin.Default()
	routes.SetupProductosRoutes(routerPrincipal)
	routerReplicacion := gin.Default()
	routes.SetupProductosRoutes(routerReplicacion)

	// servidor principal en una goroutine
	go func() {
		if err := routerPrincipal.Run(":8080"); err != nil {
			panic(err)
		}
	}()
	// servidor de replicación en otra goroutine
	go func() {
		if err := routerReplicacion.Run(":8081"); err != nil {
			panic(err)
		}
	}()
	select {}
}