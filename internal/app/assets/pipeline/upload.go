package pipeline

import (
	"time"

	"github.com/kartFr/Asset-Reuploader/internal/retry"
	"github.com/kartFr/Asset-Reuploader/internal/taskqueue"
)

// TaskQueuer abstracts *taskqueue.Queue[R] and *taskqueue.SmoothQueue[R]
// (both expose QueueTask with identical signatures).
type TaskQueuer[R any] interface {
	QueueTask(func() (R, error)) chan taskqueue.TaskResult[R]
}

// QueuedDo queues fn via q and runs retry.Do inside the worker.
// It collapses the `res := <-q.QueueTask(func() { return retry.Do(...) })`
// boilerplate triplicated across animation/sound/mesh.
func QueuedDo[R any](q TaskQueuer[R], tries int, delay time.Duration, fn func(try int) (R, error)) (R, error) {
	res := <-q.QueueTask(func() (R, error) {
		return retry.Do(retry.NewOptions(retry.Tries(tries), retry.Delay(delay)), fn)
	})
	return res.Result, res.Error
}
