package cookie

import (
	"testing"

	"github.com/fenthope/sessions"
	"github.com/fenthope/sessions/tester"
)

var newStore = func(_ *testing.T) sessions.Store {
	store := NewStore([]byte("secret"))
	return store
}

func TestCookie_SessionGetSet(t *testing.T) {
	tester.GetSet(t, newStore)
}

func TestCookie_SessionDeleteKey(t *testing.T) {
	tester.DeleteKey(t, newStore)
}

func TestCookie_SessionFlashes(t *testing.T) {
	tester.Flashes(t, newStore)
}

func TestCookie_SessionClear(t *testing.T) {
	tester.Clear(t, newStore)
}

func TestCookie_SessionOptions(t *testing.T) {
	tester.Options(t, newStore)
}

func TestCookie_SessionMany(t *testing.T) {
	tester.Many(t, newStore)
}

func TestCookie_SessionManyStores(t *testing.T) {
	tester.ManyStores(t, newStore)
}
