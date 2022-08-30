package tasks

import (
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	Run("testRun", true, 1*time.Second, testRunFN)
	time.Sleep(5 * time.Second)
}

func testRunFN(context *TaskContext) {
	flog.Info(time.Now())
}
