package main

import (
 "testing"
 "github.com/coder-lulu/newbee-job/internal/svc"
 "github.com/hibiken/asynq"
)

func TestBackgroundServicesDisabledOrUnavailable(t *testing.T) {
 ctx := &svc.ServiceContext{}
 ctx.Config.TaskConf.EnableDPTask = true
 ctx.Config.TaskConf.EnableScheduledTask = true
 if got := backgroundServices(ctx); len(got) != 0 { t.Fatal("disabled Asynq registered services") }
 ctx.Config.AsynqConf.Enable = true
 if got := backgroundServices(ctx); len(got) != 0 { t.Fatal("nil Asynq workers registered services") }
 ctx.AsynqServer = new(asynq.Server)
 ctx.AsynqScheduler = new(asynq.Scheduler)
 ctx.AsynqPTM = new(asynq.PeriodicTaskManager)
 if got := backgroundServices(ctx); len(got) != 3 { t.Fatalf("enabled service count = %d, want 3",len(got)) }
 ctx.Config.TaskConf.EnableDPTask = false
 ctx.Config.TaskConf.EnableScheduledTask = false
 if got := backgroundServices(ctx); len(got) != 1 { t.Fatalf("consumer-only service count = %d, want 1",len(got)) }
 ctx.Config.AsynqConf.Enable = false
 if got := backgroundServices(ctx); len(got) != 0 { t.Fatal("Asynq disabled with objects still registered services") }
}
