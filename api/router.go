package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type API struct {
	ready    bool
	engine   *gin.Engine
	accounts map[string]string
	s        ServiceHTTP
}

func (api *API) registerRoutes() {
	lsd := api.engine.Group("/v1/lsd", gin.BasicAuth(api.accounts))

	lsd.POST("/send-notification", func(c *gin.Context) {
		payload := new(NotificationPayload)
		if err := c.BindJSON(payload); err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		result, err := api.s.SendNotificationHTTP(payload)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	lsd.POST("/generate-new-key-pair", func(c *gin.Context) {
		payload := new(NewKeyPairPayload)
		if err := c.BindJSON(payload); err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		result, err := api.s.GenerateNewKeyPairHTTP(payload)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	lsd.POST("/get-key-pair", func(c *gin.Context) {
		payload := new(KeyPairPayload)
		if err := c.BindJSON(payload); err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		result, err := api.s.GetKeyPairHTTP(payload)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	api.ready = true
}
