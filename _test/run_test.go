package test

import (
	"context"
	"fmt"
	"github.com/farseer-go/tasks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	assert.Panics(t, func() {
		tasks.Run("testRun", 0, func(context *tasks.TaskContext) {

		}, context.Background())
	})

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	now := time.Now()
	tasks.Run("testRun", 10*time.Millisecond, func(context *tasks.TaskContext) {
		s := time.Since(now) - 10*time.Millisecond
		if s >= 6*time.Millisecond || s < time.Nanosecond {
			t.Fatal(fmt.Sprintf("时间不对:%d", s.Milliseconds()))
		}
		now = time.Now()
	}, ctx)

	time.Sleep(500 * time.Millisecond)
}

func TestRunNow(t *testing.T) {
	assert.Panics(t, func() {
		tasks.RunNow("testRunNow", 0, func(context *tasks.TaskContext) {

		}, context.Background())
	})

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	now := time.Now().Add(-10 * time.Millisecond)
	tasks.RunNow("testRunNow", 10*time.Millisecond, func(context *tasks.TaskContext) {
		s := time.Since(now) - 10*time.Millisecond
		if s >= 6*time.Millisecond || s < time.Nanosecond {
			t.Fatal(fmt.Sprintf("时间不对:%d", s.Milliseconds()))
		}
		now = time.Now()
	}, ctx)

	time.Sleep(500 * time.Millisecond)
}

func TestRunPanic(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	assert.Panics(t, func() {
		tasks.RunNow("testRun", 500*time.Millisecond, func(context *tasks.TaskContext) {
			context.SetNextDuration(0)
		}, ctx)
	})
	assert.Panics(t, func() {
		tasks.RunNow("testRun", 500*time.Millisecond, func(context *tasks.TaskContext) {
			context.SetNextTime(time.Now().Add(-1 * time.Second))
		}, ctx)
	})

	time.Sleep(50 * time.Millisecond)
}

func TestSetNextDuration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	now := time.Now().Add(20 * time.Millisecond)
	tasks.Run("testRun", 20*time.Millisecond, func(context *tasks.TaskContext) {
		s := time.Since(now)
		if s >= 6*time.Millisecond || s < time.Nanosecond {
			t.Fatal(fmt.Sprintf("时间不对:%d", s.Milliseconds()))
		}
		now = time.Now().Add(10 * time.Millisecond)
		context.SetNextDuration(10 * time.Millisecond)
	}, ctx)

	time.Sleep(500 * time.Millisecond)
}

func TestSetNextTime(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	now := time.Now().Add(20 * time.Millisecond)
	tasks.Run("testRun", 20*time.Millisecond, func(context *tasks.TaskContext) {
		s := time.Since(now)
		if s >= 6*time.Millisecond || s < time.Nanosecond {
			t.Fatal(fmt.Sprintf("时间不对:%d", s.Milliseconds()))
		}
		now = time.Now().Add(10 * time.Millisecond)
		context.SetNextTime(now)
	}, ctx)

	time.Sleep(500 * time.Millisecond)
}
