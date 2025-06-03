package main

import (
	"github.com/fenthope/sessions"
	"github.com/fenthope/sessions/mongo/mongomgo"
	"github.com/globalsign/mgo"
	"github.com/infinite-iroha/touka"
)

func main() {
	r := touka.Default()
	session, err := mgo.Dial("localhost:27017/test")
	if err != nil {
		// handle err
	}

	c := session.DB("").C("sessions")
	store := mongomgo.NewStore(c, 3600, true, []byte("secret"))
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
