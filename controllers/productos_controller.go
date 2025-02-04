package controllers

import (
	"github.com/gin-gonic/gin"
	"sync"
	"time"
	"io"
)

type Producto struct {
	ID           string `json:"id"`
	Nombre       string `json:"nombre"`
	Cantidad     int    `json:"cantidad"`
	CodigoBarras string `json:"codigo_barras"`
}

var (
	productos        []Producto
	productosMu      sync.Mutex
	ultimaActualizacion time.Time
)

// Crear un nuevo producto
func CrearProducto(c *gin.Context) {
	var nuevoProducto Producto
	if err := c.ShouldBindJSON(&nuevoProducto); err != nil {
		c.JSON(400, gin.H{"error": "Error al decodificar el cuerpo de la solicitud"})
		return
	}

	productosMu.Lock()
	productos = append(productos, nuevoProducto)
	ultimaActualizacion = time.Now()
	productosMu.Unlock()

	c.JSON(201, gin.H{"message": "Producto creado", "producto": nuevoProducto})
}

// Obtener todos los productos
func ObtenerProductos(c *gin.Context) {
	productosMu.Lock()
	defer productosMu.Unlock()

	c.JSON(200, gin.H{"productos": productos})
}

// Obtener un producto por su ID
func ObtenerProductoPorID(c *gin.Context) {
	id := c.Param("id")

	productosMu.Lock()
	defer productosMu.Unlock()

	for _, producto := range productos {
		if producto.ID == id {
			c.JSON(200, gin.H{"producto": producto})
			return
		}
	}

	c.JSON(404, gin.H{"error": "Producto no encontrado"})
}

// Actualizar un producto por su ID
func ActualizarProducto(c *gin.Context) {
	id := c.Param("id")

	var productoActualizado Producto
	if err := c.ShouldBindJSON(&productoActualizado); err != nil {
		c.JSON(400, gin.H{"error": "Error al decodificar el cuerpo de la solicitud"})
		return
	}

	productosMu.Lock()
	defer productosMu.Unlock()

	for i, producto := range productos {
		if producto.ID == id {
			productos[i] = productoActualizado
			ultimaActualizacion = time.Now()
			c.JSON(200, gin.H{"message": "Producto actualizado", "producto": productoActualizado})
			return
		}
	}

	c.JSON(404, gin.H{"error": "Producto no encontrado"})
}

// Eliminar un producto por su ID
func EliminarProducto(c *gin.Context) {
	id := c.Param("id")

	productosMu.Lock()
	defer productosMu.Unlock()

	for i, producto := range productos {
		if producto.ID == id {
			productos = append(productos[:i], productos[i+1:]...)
			ultimaActualizacion = time.Now()
			c.JSON(200, gin.H{"message": "Producto eliminado"})
			return
		}
	}

	c.JSON(404, gin.H{"error": "Producto no encontrado"})
}

// Short Pulling: El cliente hace solicitudes periódicas
func ShortPulling(c *gin.Context) {
	productosMu.Lock()
	defer productosMu.Unlock()

	c.JSON(200, gin.H{
		"productos":          productos,
		"ultima_actualizacion": ultimaActualizacion,
	})
}

// Long Pulling: El servidor mantiene la conexión abierta hasta que haya cambios
func LongPulling(c *gin.Context) {
	productosMu.Lock()
	ultimaModificacion := ultimaActualizacion
	productosMu.Unlock()

	// Esperar hasta que haya cambios o se alcance un tiempo máximo de espera
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			productosMu.Lock()
			if !ultimaActualizacion.Equal(ultimaModificacion) {
				productosMu.Unlock()
				c.JSON(200, gin.H{
					"productos":          productos,
					"ultima_actualizacion": ultimaActualizacion,
				})
				return
			}
			productosMu.Unlock()
		case <-timeout:
			c.JSON(200, gin.H{
				"message": "No hay cambios",
			})
			return
		}
	}
}

//Obtener productos uno por uno con un intervalo de tiempo
func ObtenerProductosIncremental(c *gin.Context) {
	productosMu.Lock()
	defer productosMu.Unlock()

	intervalo := 2 * time.Second
	// Crear un canal para enviar los productos uno por uno
	productosChan := make(chan Producto)
	go func() {
		for _, producto := range productos {
			productosChan <- producto
			time.Sleep(intervalo)
		}
		close(productosChan)
	}()
	// Enviar los productos al cliente uno por uno
	c.Stream(func(w io.Writer) bool {
		if producto, ok := <-productosChan; ok {
			c.SSEvent("producto", producto)
			return true
		}
		return false
	})
}