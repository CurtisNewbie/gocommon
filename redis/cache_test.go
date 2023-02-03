package redis

import (
	"testing"
	"time"

	"github.com/curtisnewbie/gocommon/common"
)

func TestRCache(t *testing.T) {
	UnmuteLockLog()
	common.LoadConfigFromFile("../app-conf-dev.yml")
	if _, e := InitRedisFromProp(); e != nil {
		t.Fatal(e)
	}
	c := common.EmptyExecContext()

	keypre := "test:rcache:key:"
	exp := 60 * time.Second

	rcache := NewRCache(exp)

	val, e := rcache.Get(c, "absent key")
	if e != nil {
		t.Fatal(e)
	}
	if val != "" {
		t.Fatal(val)
	}

	e = rcache.Put(c, keypre + "1", "2")
	if e != nil {
		t.Fatal(e)
	}

	val, e = rcache.GetElse(c, keypre + "1", nil)
	if e != nil {
		t.Fatal(e)
	}
	if val != "2" {
		t.Fatalf("val '%v' != \"2\"", val)
	}
}

func TestLKazyRCache(t *testing.T) {
	UnmuteLockLog()
	common.LoadConfigFromFile("../app-conf-dev.yml")
	InitRedisFromProp()
	c := common.EmptyExecContext()

	keypre := "test:rcache:key:"
	exp := 60 * time.Second

	rcache := NewLazyRCache(exp)
	e := rcache.Put(c, keypre + "1", "2")
	if e != nil {
		t.Fatal(e)
	}

	val, e := rcache.GetElse(c, keypre + "1", nil)
	if e != nil {
		t.Fatal(e)
	}
	if val != "2" {
		t.Fatalf("val '%v' != \"2\"", val)
	}
}

