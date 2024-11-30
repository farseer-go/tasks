package tasks

import (
	"context"
	"time"

	"github.com/farseer-go/fs/asyncLocal"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/stopwatch"
	"github.com/farseer-go/fs/trace"
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
	// 立即执行
	taskFn(&TaskContext{
		sw: stopwatch.StartNew(),
	})
	Run(taskName, interval, taskFn, ctx)
}

// 运行任务
func runTask(taskName string, interval time.Duration, taskFn func(context *TaskContext)) (nextInterval time.Duration) {
	// 这里需要提前设置默认的间隔时间。如果发生异常时，不提前设置会=0
	nextInterval = interval
	entryTask := container.Resolve[trace.IManager]().EntryTask(taskName)
	var err error
	try := exception.Try(func() {
		taskContext := &TaskContext{
			sw: stopwatch.StartNew(),
		}
		taskFn(taskContext)
		flog.ComponentInfof("task", "%s，耗时：%s", taskName, taskContext.sw.GetMillisecondsText())
		if taskContext.nextRunAt.Year() >= 2022 {
			nextInterval = taskContext.nextRunAt.Sub(time.Now())
		}
	})
	try.CatchException(func(exp any) {
		err = flog.Errorf("[%s] throw exception：%s", taskName, exp)
	})
	container.Resolve[trace.IManager]().Push(entryTask, err)
	asyncLocal.Release()
	return
}
