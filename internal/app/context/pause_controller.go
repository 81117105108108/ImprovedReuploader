package context

import "sync"

type pauseController struct {
	// IsPaused is guarded by mutex — do not read directly, use IsPausedNow().
	IsPaused bool
	mutex    sync.RWMutex
	signal   chan struct{}
}

func newPauseController() *pauseController {
	signal := make(chan struct{})
	close(signal)

	return &pauseController{
		signal: signal,
	}
}

func (c *pauseController) WaitIfPaused() {
	c.mutex.RLock()
	signal := c.signal
	c.mutex.RUnlock()
	<-signal
}

// IsPausedNow reports paused state under lock (use instead of reading IsPaused field).
func (c *pauseController) IsPausedNow() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.IsPaused
}

func (c *pauseController) Pause() (success bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.IsPaused {
		return false
	}

	c.signal = make(chan struct{})
	c.IsPaused = true
	return true
}

func (c *pauseController) Unpause() (success bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.IsPaused {
		return false
	}

	close(c.signal)
	c.IsPaused = false
	return true
}
