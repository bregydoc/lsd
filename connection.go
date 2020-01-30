package lsd

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)



func launchServer() {
	r := gin.New()
	m := melody.New()
	lock := new(sync.Mutex)

	sessionsMap := map[*melody.Session]string{}

	r.GET("/ws", func(c *gin.Context) {
		if err := m.HandleRequest(c.Writer, c.Request); err != nil {
			panic(err)
		}
	})

	m.HandleConnect(func(s *melody.Session) {
		lock.Lock()
		s.Request.

		lock.Unlock()
	})
}