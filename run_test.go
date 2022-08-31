package tasks

import (
	"context"
	"github.com/farseer-go/fs/flog"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	Run("testRun", 1*time.Second, testRunFN, context.Background())
	time.Sleep(5 * time.Second)
}
func TestRunNow(t *testing.T) {
	RunNow("testRun", 1*time.Second, testRunFN, context.Background())
	time.Sleep(5 * time.Second)
}

func testRunFN(context *TaskContext) {
	flog.Info("doing....")
}
