package test

import (
	"context"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/tasks"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	tasks.Run("testRun", 1*time.Second, testRunFN, context.Background())
	time.Sleep(5 * time.Second)
}
func TestRunNow(t *testing.T) {
	tasks.RunNow("testRun", 1*time.Second, testRunFN, context.Background())
	time.Sleep(5 * time.Second)
}

func testRunFN(context *tasks.TaskContext) {
	flog.Info("doing....")
}
