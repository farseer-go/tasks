package tasks

import (
	"context"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/stopwatch"
	"time"
)

// Run 运行一个任务，运行前先休眠
// interval:任务运行的间隔时间
// taskFn:要运行的任务
func Run(taskName string, interval time.Duration, taskFn func(context *TaskContext), ctx context.Context) {
	// 不立即运行，则先休眠interval时间
	if interval <= 0 {
		panic("interval参数，必须大于0")
	}

	go func() {
		taskInterval := interval
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(taskInterval):
				taskInterval = runTask(taskName, interval, taskFn)
			}
		}
	}()
}

// RunNow 运行一个任务
// interval:任务运行的间隔时间
// taskFn:要运行的任务
func RunNow(taskName string, interval time.Duration, taskFn func(context *TaskContext), ctx context.Context) {
	// 不立即运行，则先休眠interval时间
	if interval <= 0 {
		panic("interval参数，必须大于0")
	}

	go func() {
		taskInterval := time.Duration(0)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(taskInterval):
				taskInterval = runTask(taskName, interval, taskFn)
			}
		}
	}()
}

// 运行任务
func runTask(taskName string, interval time.Duration, taskFn func(context *TaskContext)) time.Duration {
	defer func() {
		if r := recover(); r != nil {
			flog.Errorf("taskFn [%s] throw exception：%s", taskName, r)
		}
	}()
	taskContext := &TaskContext{
		sw: stopwatch.StartNew(),
	}
	taskFn(taskContext)
	if taskContext.nextRunAt.Year() >= 2022 {
		return taskContext.nextRunAt.Sub(time.Now())
	}
	return interval
}
