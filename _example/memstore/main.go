package main

import (
	"github.com/fenthope/sessions"
	"github.com/fenthope/sessions/memstore"
	"github.com/infinite-iroha/touka"
)

func main() {
	r := touka.Default()
	store := memstore.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/incr", func(c *touka.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, touka.H{"count": count})
	})
	r.Run(":8000")
}
