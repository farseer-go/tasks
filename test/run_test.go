package test

import (
	"context"
	"github.com/farseer-go/tasks"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	assert.Panics(t, func() {
		tasks.Run("testRun", 0, func(context *tasks.TaskContext) {

		}, context.Background())
	})

	lock := &sync.Mutex{}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	count := 0
	tasks.Run("testRun", 10*time.Millisecond, func(context *tasks.TaskContext) {
		lock.Lock()
		defer lock.Unlock()
		count++
	}, ctx)

	time.Sleep(50 * time.Millisecond)
	lock.Lock()
	defer lock.Unlock()

	assert.Equal(t, 4, count)
}

func TestRunNow(t *testing.T) {
	assert.Panics(t, func() {
		tasks.RunNow("testRunNow", 0, func(context *tasks.TaskContext) {

		}, context.Background())
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	var lock sync.Mutex
	count := 0
	tasks.RunNow("testRunNow", 10*time.Millisecond, func(context *tasks.TaskContext) {
		lock.Lock()
		defer lock.Unlock()

		count++
	}, ctx)

	time.Sleep(50 * time.Millisecond)

	lock.Lock()
	defer lock.Unlock()

	assert.Equal(t, 5, count)
}

func TestRunPanic(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	tasks.RunNow("testRun", 500*time.Millisecond, func(context *tasks.TaskContext) {
		context.SetNextDuration(0)
	}, ctx)
	tasks.RunNow("testRun", 500*time.Millisecond, func(context *tasks.TaskContext) {
		context.SetNextTime(time.Now().Add(-1 * time.Second))
	}, ctx)

	time.Sleep(50 * time.Millisecond)
}

func TestSetNextDuration(t *testing.T) {
	count := 0
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	var lock sync.Mutex
	tasks.Run("testRun", 20*time.Millisecond, func(context *tasks.TaskContext) {
		lock.Lock()
		defer lock.Unlock()

		context.SetNextDuration(10 * time.Millisecond)
		count++
	}, ctx)
	time.Sleep(50 * time.Millisecond)
	lock.Lock()
	defer lock.Unlock()

	assert.Equal(t, 3, count)
}

func TestSetNextTime(t *testing.T) {
	var lock sync.Mutex
	count := 0
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	tasks.Run("testRun", 20*time.Millisecond, func(context *tasks.TaskContext) {
		lock.Lock()
		defer lock.Unlock()

		context.SetNextTime(time.Now().Add(10 * time.Millisecond))
		count++
	}, ctx)
	time.Sleep(50 * time.Millisecond)
	cancel()
	lock.Lock()
	defer lock.Unlock()
	
	assert.Equal(t, 3, count)
}
