package timer

import (
    "time"
    "sync"
)

type Ticker struct {
    duration time.Duration
    tickChan chan struct{}
    lock sync.Mutex
    stopped bool
    curNum int64
}

func New(duration time.Duration) *Ticker {
    return &Ticker{
        duration:duration,
        tickChan:make(chan struct{}),
    }
}

func (t *Ticker) Tick() <-chan struct{} {
    t.restart(0)
    return t.tickChan
}

func (t *Ticker) restart(timestamp int64) {
    go func() {
        for {
            time.Sleep(t.duration)
            t.lock.Lock()
            if t.curNum != timestamp {
                t.lock.Unlock()
                break
            } else {
                t.lock.Unlock()
            }
        }
    }()
}

func (t *Ticker) Reset() {
    t.lock.Lock()
    defer t.lock.Unlock()
    t.curNum++
    t.restart(t.curNum)
}

func (t *Ticker) Stop() {
    close(t.tickChan)
}