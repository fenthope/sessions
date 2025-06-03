package main

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/fenthope/sessions"
	"github.com/fenthope/sessions/memcached"
	"github.com/infinite-iroha/touka"
)

func main() {
	r := touka.Default()
	store := memcached.NewStore(memcache.New("localhost:11211"), "", []byte("secret"))
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
