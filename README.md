# tasks
Process-level tasks

## What are the functions?
* tasks
    * struct
        * TaskContext（运行任务的上下文）
            * .SetNextTime（设置休眠时间）
            * .SetNextDuration（设置休眠时间）
    * func
        * Run （先休眠再运行任务）
        * RunNow （立即运行任务）

## Getting Started
```go
// custom your job
// context can set Execute NextTime Or Use Default Interval
func testRun(context *tasks.TaskContext) {
	flog.Info("doing...")
}

// "testRun" = taskName
// 1*time.Second = Interval time
// testRun = job func
tasks.Run("testRun", 1*time.Second, testRun, context.Background())

// "testRun" = taskName
// 1*time.Second = Interval time
// testRun = job func
tasks.RunNow("testRun", 1*time.Second, testRun, context.Background())
```