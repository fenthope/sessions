// Package tester is a package to test each packages of session stores, such as
// cookie, redis, memcached, mongo, memstore.  You can use this to test your own session
// stores.
package tester

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fenthope/sessions"
	"github.com/infinite-iroha/touka"
)

type storeFactory func(*testing.T) sessions.Store

const (
	sessionName = "mysession"
	ok          = "ok"
)

func GetSet(t *testing.T, newStore storeFactory) {
	r := touka.Default()
	r.Use(sessions.Sessions(sessionName, newStore(t)))

	r.GET("/set", func(c *touka.Context) {
		session := sessions.Default(c)
		session.Set("key", ok)
		_ = session.Save()
		c.String(http.StatusOK, ok)
	})

	r.GET("/get", func(c *touka.Context) {
		session := sessions.Default(c)
		if session.Get("key") != ok {
			t.Error("Session writing failed")
		}
		_ = session.Save()
		c.String(http.StatusOK, ok)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/set", nil)
	r.ServeHTTP(res1, req1)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/get", nil)
	copyCookies(req2, res1)
	r.ServeHTTP(res2, req2)
}

func DeleteKey(t *testing.T, newStore storeFactory) {
	r := touka.Default()
	r.Use(sessions.Sessions(sessionName, newStore(t)))

	r.GET("/set", func(c *touka.Context) {
		session := sessions.Default(c)
		session.Set("key", ok)
		_ = session.Save()
		c.String(http.StatusOK, ok)
	})

	r.GET("/delete", func(c *touka.Context) {
		session := sessions.Default(c)
		session.Delete("key")
		_ = session.Save()
		c.String(http.StatusOK, ok)
	})

	r.GET("/get", func(c *touka.Context) {
		session := sessions.Default(c)
		if session.Get("key") != nil {
			t.Error("Session deleting failed")
		}
		_ = session.Save()
		c.String(http.StatusOK, ok)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/set", nil)
	r.ServeHTTP(res1, req1)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/delete", nil)
	copyCookies(req2, res1)
	r.ServeHTTP(res2, req2)

	res3 := httptest.NewRecorder()
	req3, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/get", nil)
	copyCookies(req3, res2)
	r.ServeHTTP(res3, req3)
}

func Flashes(t *testing.T, newStore storeFactory) {
	r := touka.Default()
	store := newStore(t)
	r.Use(sessions.Sessions(sessionName, store))

	r.GET("/set", func(c *touka.Context) {
		session := sessions.Default(c)
		session.AddFlash(ok)
		_ = session.Save()
		c.String(http.StatusOK, ok)
	})

	r.GET("/flash", func(c *touka.Context) {
		session := sessions.Default(c)
		l := len(session.Flashes())
		if l != 1 {
			t.Error("Flashes count does not equal 1. Equals ", l)
		}
		_ = session.Save()
		c.String(http.StatusOK, ok)
	})

	r.GET("/check", func(c *touka.Context) {
		session := sessions.Default(c)
		l := len(session.Flashes())
		if l != 0 {
			t.Error("flashes count is not 0 after reading. Equals ", l)
		}
		_ = session.Save()
		c.String(http.StatusOK, ok)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/set", nil)
	r.ServeHTTP(res1, req1)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/flash", nil)
	copyCookies(req2, res1)
	r.ServeHTTP(res2, req2)

	res3 := httptest.NewRecorder()
	req3, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/check", nil)
	copyCookies(req3, res2)
	r.ServeHTTP(res3, req3)
}

func Clear(t *testing.T, newStore storeFactory) {
	data := map[string]string{
		"key": "val",
		"foo": "bar",
	}
	r := touka.Default()
	store := newStore(t)
	r.Use(sessions.Sessions(sessionName, store))

	r.GET("/set", func(c *touka.Context) {
		session := sessions.Default(c)
		for k, v := range data {
			session.Set(k, v)
		}
		session.Clear()
		_ = session.Save()
		c.String(http.StatusOK, ok)
	})

	r.GET("/check", func(c *touka.Context) {
		session := sessions.Default(c)
		for k, v := range data {
			if session.Get(k) == v {
				t.Fatal("Session clear failed")
			}
		}
		_ = session.Save()
		c.String(http.StatusOK, ok)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/set", nil)
	r.ServeHTTP(res1, req1)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/check", nil)
	copyCookies(req2, res1)
	r.ServeHTTP(res2, req2)
}

func Options(t *testing.T, newStore storeFactory) {
	r := touka.Default()
	store := newStore(t)
	store.Options(sessions.Options{
		Domain: "localhost",
	})
	r.Use(sessions.Sessions(sessionName, store))

	r.GET("/domain", func(c *touka.Context) {
		session := sessions.Default(c)
		session.Set("key", ok)
		session.Options(sessions.Options{
			Path: "/foo/bar/bat",
		})
		_ = session.Save()
		c.String(http.StatusOK, ok)
	})
	r.GET("/path", func(c *touka.Context) {
		session := sessions.Default(c)
		session.Set("key", ok)
		_ = session.Save()
		c.String(http.StatusOK, ok)
	})
	r.GET("/set", func(c *touka.Context) {
		session := sessions.Default(c)
		session.Set("key", ok)
		_ = session.Save()
		c.String(http.StatusOK, ok)
	})
	r.GET("/expire", func(c *touka.Context) {
		session := sessions.Default(c)
		session.Options(sessions.Options{
			MaxAge: -1,
		})
		_ = session.Save()
		c.String(http.StatusOK, ok)
	})
	r.GET("/check", func(c *touka.Context) {
		session := sessions.Default(c)
		val := session.Get("key")
		if val != nil {
			t.Fatal("Session expiration failed")
		}
		c.String(http.StatusOK, ok)
	})

	testOptionSameSitego(t, r)

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/domain", nil)
	r.ServeHTTP(res1, req1)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/path", nil)
	r.ServeHTTP(res2, req2)

	res3 := httptest.NewRecorder()
	req3, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/set", nil)
	r.ServeHTTP(res3, req3)

	res4 := httptest.NewRecorder()
	req4, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/expire", nil)
	r.ServeHTTP(res4, req4)

	res5 := httptest.NewRecorder()
	req5, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/check", nil)
	r.ServeHTTP(res5, req5)

	for _, c := range res1.Header().Values("Set-Cookie") {
		s := strings.Split(c, "; ")
		if len(s) < 2 {
			t.Fatal("No Path=/foo/bar/batt found in options")
		}
		if s[1] != "Path=/foo/bar/bat" {
			t.Fatal("Error writing path with options:", s[1])
		}
	}

	for _, c := range res2.Header().Values("Set-Cookie") {
		s := strings.Split(c, "; ")
		if len(s) < 2 {
			t.Fatal("No Domain=localhost found in options")
		}

		if s[1] != "Domain=localhost" {
			t.Fatal("Error writing domain with options:", s[1])
		}
	}
}

func Many(t *testing.T, newStore storeFactory) {
	r := touka.Default()
	sessionNames := []string{"a", "b"}

	r.Use(sessions.SessionsMany(sessionNames, newStore(t)))

	r.GET("/set", func(c *touka.Context) {
		sessionA := sessions.DefaultMany(c, "a")
		sessionA.Set("hello", "world")
		_ = sessionA.Save()

		sessionB := sessions.DefaultMany(c, "b")
		sessionB.Set("foo", "bar")
		_ = sessionB.Save()
		c.String(http.StatusOK, ok)
	})

	r.GET("/get", func(c *touka.Context) {
		sessionA := sessions.DefaultMany(c, "a")
		if sessionA.Get("hello") != "world" {
			t.Error("Session writing failed")
		}
		_ = sessionA.Save()

		sessionB := sessions.DefaultMany(c, "b")
		if sessionB.Get("foo") != "bar" {
			t.Error("Session writing failed")
		}
		_ = sessionB.Save()
		c.String(http.StatusOK, ok)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/set", nil)
	r.ServeHTTP(res1, req1)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/get", nil)
	header := ""
	for _, x := range res1.Header()["Set-Cookie"] {
		header += strings.Split(x, ";")[0] + "; \n"
	}
	req2.Header.Set("Cookie", header)
	r.ServeHTTP(res2, req2)
}

func ManyStores(t *testing.T, newStore storeFactory) {
	r := touka.Default()

	store := newStore(t)
	sessionStores := []sessions.SessionStore{
		{Name: "a", Store: store},
		{Name: "b", Store: store},
	}

	r.Use(sessions.SessionsManyStores(sessionStores))

	r.GET("/set", func(c *touka.Context) {
		sessionA := sessions.DefaultMany(c, "a")
		sessionA.Set("hello", "world")
		_ = sessionA.Save()

		sessionB := sessions.DefaultMany(c, "b")
		sessionB.Set("foo", "bar")
		_ = sessionB.Save()
		c.String(http.StatusOK, ok)
	})

	r.GET("/get", func(c *touka.Context) {
		sessionA := sessions.DefaultMany(c, "a")
		if sessionA.Get("hello") != "world" {
			t.Error("Session writing failed")
		}
		_ = sessionA.Save()

		sessionB := sessions.DefaultMany(c, "b")
		if sessionB.Get("foo") != "bar" {
			t.Error("Session writing failed")
		}
		_ = sessionB.Save()
		c.String(http.StatusOK, ok)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequestWithContext(context.Background(), "GET", "/set", nil)
	r.ServeHTTP(res1, req1)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequestWithContext(context.Background(), "GET", "/get", nil)
	header := ""
	for _, x := range res1.Header()["Set-Cookie"] {
		header += strings.Split(x, ";")[0] + "; \n"
	}
	req2.Header.Set("Cookie", header)
	r.ServeHTTP(res2, req2)
}

func copyCookies(req *http.Request, res *httptest.ResponseRecorder) {
	req.Header.Set("Cookie", strings.Join(res.Header().Values("Set-Cookie"), "; "))
}
