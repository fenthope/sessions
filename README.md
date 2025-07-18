# sessions

Touka(灯花)框架的sessions中间件, 移植自[Gin sessions](https://github.com/gin-contrib/sessions)

Touka middleware for session management with multi-backend support:

- [cookie-based](#cookie-based)
- [Redis](#redis)
- [memcached](#memcached)
- [MongoDB](#mongodb)
- [GORM](#gorm)
- [memstore](#memstore)
- [PostgreSQL](#postgresql)
- [Filesystem](#Filesystem)

## 许可证 License

MIT License Copyright (c) 2016 Gin-Gonic

Apache 2.0 License Copyright © 2025 Infinite-Iroha & WJQSERVER. All rights reserved.

## Usage

### Start using it

Download and install it:

```bash
go get github.com/fenthope/sessions
```

Import it in your code:

```go
import "github.com/fenthope/sessions"
```

## Basic Examples

### single session

```go
package main

import (
  "github.com/fenthope/sessions"
  "github.com/fenthope/sessions/cookie"
  "github.com/infinite-iroha/touka"
)

func main() {
  r := touka.Default()
  store := cookie.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("mysession", store))

  r.GET("/hello", func(c *touka.Context) {
    session := sessions.Default(c)

    if session.Get("hello") != "world" {
      session.Set("hello", "world")
      session.Save()
    }

    c.JSON(200, touka.H{"hello": session.Get("hello")})
  })
  r.Run(":8000")
}
```

### multiple sessions

```go
package main

import (
  "github.com/fenthope/sessions"
  "github.com/fenthope/sessions/cookie"
  "github.com/infinite-iroha/touka"
)

func main() {
  r := touka.Default()
  store := cookie.NewStore([]byte("secret"))
  sessionNames := []string{"a", "b"}
  r.Use(sessions.SessionsMany(sessionNames, store))

  r.GET("/hello", func(c *touka.Context) {
    sessionA := sessions.DefaultMany(c, "a")
    sessionB := sessions.DefaultMany(c, "b")

    if sessionA.Get("hello") != "world!" {
      sessionA.Set("hello", "world!")
      sessionA.Save()
    }

    if sessionB.Get("hello") != "world?" {
      sessionB.Set("hello", "world?")
      sessionB.Save()
    }

    c.JSON(200, touka.H{
      "a": sessionA.Get("hello"),
      "b": sessionB.Get("hello"),
    })
  })
  r.Run(":8000")
}
```

### multiple sessions with different stores

```go
package main

import (
  "github.com/fenthope/sessions"
  "github.com/fenthope/sessions/cookie"
  "github.com/infinite-iroha/touka"
)

func main() {
  r := touka.Default()
  cookieStore := cookie.NewStore([]byte("secret"))
  redisStore, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
  sessionStores := []sessions.SessionStore{
    {
      Name:  "a",
      Store: cookieStore,
    },
    {
      Name:  "b",
      Store: redisStore,
    },
  }
  r.Use(sessions.SessionsManyStores(sessionStores))

  r.GET("/hello", func(c *touka.Context) {
    sessionA := sessions.DefaultMany(c, "a")
    sessionB := sessions.DefaultMany(c, "b")

    if sessionA.Get("hello") != "world!" {
      sessionA.Set("hello", "world!")
      sessionA.Save()
    }

    if sessionB.Get("hello") != "world?" {
      sessionB.Set("hello", "world?")
      sessionB.Save()
    }

    c.JSON(200, touka.H{
      "a": sessionA.Get("hello"),
      "b": sessionB.Get("hello"),
    })
  })
  r.Run(":8000")
}
```

## Backend Examples

### cookie-based

```go
package main

import (
  "github.com/fenthope/sessions"
  "github.com/fenthope/sessions/cookie"
  "github.com/infinite-iroha/touka"
)

func main() {
  r := touka.Default()
  store := cookie.NewStore([]byte("secret"))
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
```

### Redis

```go
package main

import (
  "github.com/fenthope/sessions"
  "github.com/fenthope/sessions/redis"
  "github.com/infinite-iroha/touka"
)

func main() {
  r := touka.Default()
  store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
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
```

### Memcached

#### ASCII Protocol

```go
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
```

#### Binary protocol (with optional SASL authentication)

```go
package main

import (
  "github.com/fenthope/sessions"
  "github.com/fenthope/sessions/memcached"
  "github.com/infinite-iroha/touka"
  "github.com/memcachier/mc"
)

func main() {
  r := touka.Default()
  client := mc.NewMC("localhost:11211", "username", "password")
  store := memcached.NewMemcacheStore(client, "", []byte("secret"))
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
```

### MongoDB

#### mgo

```go
package main

import (
  "github.com/fenthope/sessions"
  "github.com/fenthope/sessions/mongo/mongomgo"
  "github.com/infinite-iroha/touka"
  "github.com/globalsign/mgo"
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
```

#### mongo-driver

```go
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
```

### memstore

```go
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
```

### GORM

```go
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
```

### PostgreSQL

```go
package main

import (
  "database/sql"
  "github.com/fenthope/sessions"
  "github.com/fenthope/sessions/postgres"
  "github.com/infinite-iroha/touka"
)

func main() {
  r := touka.Default()
  db, err := sql.Open("postgres", "postgresql://username:password@localhost:5432/database")
  if err != nil {
    // handle err
  }

  store, err := postgres.NewStore(db, []byte("secret"))
  if err != nil {
    // handle err
  }

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
```

### Filesystem

```go
package main

import (
  "github.com/fenthope/sessions"
  "github.com/fenthope/sessions/filesystem"
  "github.com/infinite-iroha/touka"
)

func main() {
  r := touka.Default()

  var sessionPath = "/tmp/" // in case of empty string, the system's default tmp folder is used

  store := filesystem.NewStore(sessionPath,[]byte("secret"))
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
```
