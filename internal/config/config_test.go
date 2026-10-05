package config

import (
	"github.com/zeromicro/go-zero/core/conf"
	"testing"
)

func TestTaskSchedulersDefaultDisabled(t *testing.T) {
	var c TaskConf
	if err := conf.LoadFromJsonBytes([]byte(`{}`), &c); err != nil {
		t.Fatal(err)
	}
	if c.EnableDPTask || c.EnableScheduledTask {
		t.Fatalf("unsafe scheduler defaults: %+v", c)
	}
	if err := conf.LoadFromJsonBytes([]byte(`{"EnableDPTask":true,"EnableScheduledTask":true}`), &c); err != nil {
		t.Fatal(err)
	}
	if !c.EnableDPTask || !c.EnableScheduledTask {
		t.Fatalf("explicit settings ignored: %+v", c)
	}
}
