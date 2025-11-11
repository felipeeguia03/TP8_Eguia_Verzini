package main

import (
	"net/http"
	"os"
	"strings"

	"backend/app"
	"backend/clients"

	"github.com/gin-gonic/gin"
)

func normalize(s string) string { return strings.TrimRight(s, "/") }

func main() {
	clients.StartDB()

	// Whitelist de orígenes permitidos (QA, PROD, Render y dev local)
	allowedOrigins := map[string]struct{}{
		"http://localhost:3000": {},
		"https://webapp-tp05-front-qa-verzini-eguia-b6h9asc9g5axc6hz.chilecentral-01.azurewebsites.net":   {},
		"https://webapp-tp05-front-prod-verzini-eguia-cjdkf2dac9g3g2as.chilecentral-01.azurewebsites.net": {},
		"https://tp8-frontend.onrender.com": {},
		"https://tp8-frontend-qa.onrender.com": {},
		"https://tp8-frontend-prod.onrender.com": {},
	}

	engine := gin.New()

	engine.Use(func(c *gin.Context) {
		origin := normalize(c.Request.Header.Get("Origin"))

		// Para que caches intermedios no mezclen respuestas por origen
		c.Writer.Header().Add("Vary", "Origin")
		c.Writer.Header().Add("Vary", "Access-Control-Request-Method")
		c.Writer.Header().Add("Vary", "Access-Control-Request-Headers")

		if origin != "" {
			if _, ok := allowedOrigins[origin]; ok {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
				// Si el preflight no envía estos, damos defaults sensatos
				if m := c.Request.Header.Get("Access-Control-Request-Method"); m != "" {
					c.Writer.Header().Set("Access-Control-Allow-Methods", m)
				} else {
					c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				}
				if h := c.Request.Header.Get("Access-Control-Request-Headers"); h != "" {
					c.Writer.Header().Set("Access-Control-Allow-Headers", h)
				} else {
					c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Auth-Token")
				}
				c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
			}
		}

		// Preflight
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent) // 204
			return
		}

		c.Next()
	})

	app.MapRoutes(engine)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	_ = engine.Run(":" + port)
}
