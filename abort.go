package gen

import "sync"

// abortGroup propagates one generation failure to every participant of a
// concurrent generation round. Abort records only the first error and closes
// a broadcast channel, so spawners and workers can observe cancellation
// through a non-blocking receive without ever blocking on error delivery.
type abortGroup struct {
	once sync.Once
	done chan struct{}
	errs chan error
}

func newAbortGroup() *abortGroup {
	return &abortGroup{done: make(chan struct{}), errs: make(chan error, 1)}
}

// Aborted returns a channel that is closed once Abort has been called.
func (a *abortGroup) Aborted() <-chan struct{} {
	return a.done
}

// Abort records err as the round's failure and broadcasts cancellation.
// The send never blocks: errs has capacity one and Abort is the sole sender
// (guarded by once). The send happens before done is closed, so any caller
// observing the closed channel can read the recorded error afterwards.
func (a *abortGroup) Abort(err error) {
	a.once.Do(func() {
		a.errs <- err
		close(a.done)
	})
}

// Err drains the single recorded error, if any. Call it after all workers of
// the round have finished so no failure can race with the drain.
func (a *abortGroup) Err() error {
	select {
	case err := <-a.errs:
		return err
	default:
		return nil
	}
}
