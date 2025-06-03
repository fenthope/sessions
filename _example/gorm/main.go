package main

import (
	"github.com/fenthope/sessions"
	gormsessions "github.com/fenthope/sessions/gorm"
	"github.com/infinite-iroha/touka"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	store := gormsessions.NewStore(db, true, []byte("secret"))

	r := touka.Default()
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
