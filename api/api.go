package api

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

func New(service ServiceHTTP, accounts map[string]string) *API {
	if accounts == nil || len(accounts) == 0 {
		accounts = map[string]string{
			"admin": "admin",
		}
	}
	gin.SetMode(gin.ReleaseMode)
	p := &API{accounts: accounts, s: service, engine: gin.New()}
	p.registerRoutes()
	return p
}

func (api *API) Run(port int) error {
	if !api.ready {
		return errors.New("api not initialized")
	}
	return api.engine.Run(fmt.Sprintf(":%d", port))
}
