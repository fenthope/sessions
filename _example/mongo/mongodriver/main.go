package main

import (
	"context"

	"github.com/fenthope/sessions"
	"github.com/fenthope/sessions/mongo/mongodriver"
	"github.com/infinite-iroha/touka"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := touka.Default()
	mongoOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(mongoOptions)
	if err != nil {
		// handle err
	}

	if err := client.Connect(context.Background()); err != nil {
		// handle err
	}

	c := client.Database("test").Collection("sessions")
	store := mongodriver.NewStore(c, 3600, true, []byte("secret"))
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
