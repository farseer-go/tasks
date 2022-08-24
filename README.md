# tasks
Process-level tasks

## What are the functions?
* `tasks`
    * `struct`
        * `TaskContext（运行任务的上下文）`
            * `.SetNextTime（设置休眠时间）`
            * `.SetNextDuration（设置休眠时间）`
    * `func`
        * `Run （运行一个任务）`

## Getting Started
```go
// custom your job
// context can set Execute NextTime Or Use Default Interval
func testRun(context *tasks.TaskContext) {
	fmt.Println(time.Now())
}

// "testRun" = taskName
// true = Execute Now
// 1*time.Second = Interval time
// testRun = job func
tasks.Run("testRun", true, 1*time.Second, testRun)
```