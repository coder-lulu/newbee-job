package base

import (
	"context"
	"os"
	"testing"

	"github.com/coder-lulu/newbee-job/ent"
	_ "github.com/coder-lulu/newbee-job/ent/runtime"
	"github.com/coder-lulu/newbee-job/internal/svc"
	"github.com/coder-lulu/newbee-job/types/job"
	_ "github.com/go-sql-driver/mysql"
)

// NEWBEE_JOB_TEST_DSN must target a dedicated empty test database.
// The test does not start queue consumers or schedule any task.
func TestInitDatabaseEmptyAndExistingTasks(t *testing.T) {
	dsn := os.Getenv("NEWBEE_JOB_TEST_DSN")
	if dsn == "" {
		t.Skip("set NEWBEE_JOB_TEST_DSN to a dedicated empty MySQL test database")
	}
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		t.Fatal("could not open test database")
	}
	t.Cleanup(func() { client.Close() })
	ctx := context.Background()
	logic := NewInitDatabaseLogic(ctx, &svc.ServiceContext{DB: client})
	for i := 0; i < 2; i++ {
		resp, err := logic.InitDatabase(&job.Empty{})
		if err != nil || resp == nil {
			t.Fatalf("initialization %d failed", i)
		}
		count, err := client.Task.Query().Count(ctx)
		if err != nil || count != 0 {
			t.Fatalf("expected empty task table after initialization; count=%d", count)
		}
	}
	existing, err := client.Task.Create().SetName("bootstrap-regression").SetTaskGroup("test").SetCronExpression("@every 60s").SetPattern("bootstrap:regression").SetPayload(`{"preserve":true}`).SetStatus(2).Save(ctx)
	if err != nil {
		t.Fatal("could not create test task")
	}
	t.Cleanup(func() { client.Task.DeleteOneID(existing.ID).Exec(context.Background()) })
	existing, err = client.Task.Get(ctx, existing.ID)
	if err != nil {
		t.Fatal("could not read baseline test task")
	}
	for i := 0; i < 2; i++ {
		resp, err := logic.InitDatabase(&job.Empty{})
		if err != nil || resp == nil {
			t.Fatalf("repeat initialization %d failed", i)
		}
		rows, err := client.Task.Query().All(ctx)
		if err != nil || len(rows) != 1 {
			t.Fatal("initialization changed task count")
		}
		got := rows[0]
		if got.ID != existing.ID || got.Name != existing.Name || got.Status != existing.Status || got.TaskGroup != existing.TaskGroup || got.CronExpression != existing.CronExpression || got.Pattern != existing.Pattern || got.Payload != existing.Payload || !got.CreatedAt.Equal(existing.CreatedAt) || !got.UpdatedAt.Equal(existing.UpdatedAt) {
			t.Fatal("initialization changed existing task")
		}
	}
}
