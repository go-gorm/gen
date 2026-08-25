package pools

import (
	"testing"
	"time"
)

const poolTestTimeout = time.Second

func TestPoolTracksCapacityAndBlocksAtLimit(t *testing.T) {
	pool := NewPool(1)
	if got, want := pool.Size(), 1; got != want {
		t.Fatalf("Size() = %d, want %d", got, want)
	}
	if got := pool.Num(); got != 0 {
		t.Fatalf("initial Num() = %d, want 0", got)
	}

	pool.Wait()
	if got := pool.Num(); got != 1 {
		t.Fatalf("Num() after Wait = %d, want 1", got)
	}

	started := make(chan struct{})
	acquired := make(chan struct{})
	go func() {
		close(started)
		pool.Wait()
		close(acquired)
	}()
	<-started

	select {
	case <-acquired:
		t.Fatal("second Wait should block while the pool is full")
	case <-time.After(20 * time.Millisecond):
	}

	pool.Done()
	select {
	case <-acquired:
	case <-time.After(poolTestTimeout):
		t.Fatal("second Wait did not continue after a token was returned")
	}
	pool.Done()
	pool.WaitAll()

	if got := pool.Num(); got != 0 {
		t.Fatalf("final Num() = %d, want 0", got)
	}
}

func TestPoolWaitAllContracts(t *testing.T) {
	pool := NewPool(2)
	pool.Wait()
	pool.Wait()

	done := pool.AsyncWaitAll()
	select {
	case <-done:
		t.Fatal("AsyncWaitAll completed before tokens were returned")
	case <-time.After(20 * time.Millisecond):
	}

	pool.Done()
	pool.Done()
	select {
	case <-done:
	case <-time.After(poolTestTimeout):
		t.Fatal("AsyncWaitAll did not complete")
	}
}

func TestPoolNegativeSizeIsNoop(t *testing.T) {
	pool := NewPool(-1)
	if pool.Size() != 0 || pool.Num() != 0 {
		t.Fatalf("negative-size pool should be disabled: size=%d num=%d", pool.Size(), pool.Num())
	}
	pool.Wait()
	pool.Done()
	pool.WaitAll()

	select {
	case <-pool.AsyncWaitAll():
	case <-time.After(poolTestTimeout):
		t.Fatal("disabled pool should already be idle")
	}
}

func TestPoolZeroSizeReportsNoCapacity(t *testing.T) {
	pool := NewPool(0)
	if pool.Size() != 0 || pool.Num() != 0 {
		t.Fatalf("zero-size pool: size=%d num=%d", pool.Size(), pool.Num())
	}
}
