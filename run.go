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

	flog.Infof("启动本地任务：%s，间隔时间：%s", taskName, interval.String())
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
	// InitContext 初始化同一协程上下文，避免在同一协程中多次初始化
	asyncLocal.InitContext()
	traceContext := container.Resolve[trace.IManager]().EntryTask(taskName)
	exception.Try(func() {
		taskContext := &TaskContext{
			sw: stopwatch.StartNew(),
		}
		taskFn(taskContext)
		flog.ComponentInfof("task", "%s，耗时：%s", taskName, taskContext.sw.GetMillisecondsText())
		if taskContext.nextRunAt.Year() >= 2022 {
			nextInterval = time.Until(taskContext.nextRunAt)
		}
	}).CatchException(func(exp any) {
		if traceContext.IsIgnore() { // 如果忽略了链路,则要在这里打印错误日志
			flog.Errorf("task %s 异常: %v", taskName, exp)
		}
	})
	container.Resolve[trace.IManager]().Push(traceContext, nil)
	asyncLocal.Release()
	return
}
