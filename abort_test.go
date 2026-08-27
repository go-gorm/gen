package gen

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestAbortGroup(t *testing.T) {
	errs := newAbortGroup()
	if err := errs.Err(); err != nil {
		t.Fatalf("Err() before Abort: expected nil, got %v", err)
	}
	select {
	case <-errs.Aborted():
		t.Fatal("Aborted() should not be closed before Abort")
	default:
	}

	sentinel := errors.New("first failure")
	errs.Abort(sentinel)
	errs.Abort(errors.New("ignored")) // only the first error wins

	if got := errs.Err(); !errors.Is(got, sentinel) {
		t.Fatalf("Err(): expected first error, got %v", got)
	}
	select {
	case <-errs.Aborted():
	default:
		t.Fatal("Aborted() should be closed after Abort")
	}
}

func TestAbortGroupConcurrentAbort(t *testing.T) {
	errs := newAbortGroup()
	const n = 64
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			errs.Abort(fmt.Errorf("failure %d", i))
		}(i)
	}
	wg.Wait()

	if got := errs.Err(); got == nil {
		t.Fatal("expected exactly one recorded error after concurrent Aborts")
	}
	select {
	case <-errs.Aborted():
	default:
		t.Fatal("Aborted() should be closed after concurrent Aborts")
	}
}
