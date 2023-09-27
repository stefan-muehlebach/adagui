package binding

//-----------------------------------------------------------------------------

/*
type BasicBinder struct {
    callback atomic.Value
    dataListenerPairLock sync.RWMutex
    dataListenerPair annotatedListener
}

func (binder *BasicBinder) Bind(data DataItem) {
    listener := NewDataListener(func() {
        f := binder.callback.Load()
        if f != nil {
            f.(func(DataItem))(data)
        }
    })
    data.AddListener(listener)
    listenerInfo := annotatedListener {
        data:     data,
        listener: listener,
    }

    binder.dataListenerPairLock.Lock()
    binder.unbindLocked()
    binder.dataListenerPair = listenerInfo
    binder.dataListenerPairLock.Unlock()
}

func (binder *BasicBinder) CallWithData(f func(data DataItem)) {
    binder.dataListenerPairLock.RLock()
    data := binder.dataListenerPair.data
    binder.dataListenerPairLock.RUnlock()
    f(data)
}

func (binder *BasicBinder) SetCallback(f func(data DataItem)) {
    binder.callback.Store(f)
}

func (binder *BasicBinder) Unbind() {
    binder.dataListenerPairLock.Lock()
    binder.unbindLocked()
    binder.dataListenerPairLock.Unlock()
}

func (binder *BasicBinder) unbindLocked() {
    previousListener := binder.dataListenerPair
    binder.dataListenerPair = annotatedListener{nil, nil}

    if previousListener.listener == nil || previousListener.data == nil {
        return
    }
    previousListener.data.RemoveListener(previousListener.listener)
}

type annotatedListener struct {
    data     DataItem
    listener DataListener
}

//-----------------------------------------------------------------------------

var (
    once sync.Once
    queue *UnboundedFuncChan
)

func queueItem(f func()) {
    once.Do(func() {
        queue = NewUnboundedFuncChan()
        go func() {
            for f := range queue.Out() {
                f()
            }
        }()
    })
    queue.In() <- f
}

func waitForItems() {
    done := make(chan struct{})
    queue.In() <- func() { close(done) }
    <- done
}

//-----------------------------------------------------------------------------

type UnboundedFuncChan struct {
    in, out chan func()
    close   chan struct{}
    q       []func()
}

func NewUnboundedFuncChan() *UnboundedFuncChan {
    ch := &UnboundedFuncChan{
        // The size of Func is less than 16 bytes, we use 16 to fit
        // a CPU cache line (L2, 256 Bytes), which may reduce cache misses.
        in:    make(chan func(), 16),
        out:   make(chan func(), 16),
        close: make(chan struct{}),
    }
    go ch.processing()
    return ch
}

// In returns the send channel of the given channel, which can be used to
// send values to the channel.
func (ch *UnboundedFuncChan) In() chan<- func() { return ch.in }

// Out returns the receive channel of the given channel, which can be used
// to receive values from the channel.
func (ch *UnboundedFuncChan) Out() <-chan func() { return ch.out }

// Close closes the channel.
func (ch *UnboundedFuncChan) Close() { ch.close <- struct{}{} }

func (ch *UnboundedFuncChan) processing() {
    // This is a preallocation of the internal unbounded buffer.
    // The size is randomly picked. But if one changes the size, the
    // reallocation size at the subsequent for loop should also be
    // changed too. Furthermore, there is no memory leak since the
    // queue is garbage collected.
    ch.q = make([]func(), 0, 1<<10)
    for {
        select {
        case e, ok := <-ch.in:
            if !ok {
                // We don't want the input channel be accidentally closed
                // via close() instead of Close(). If that happens, it is
                // a misuse, do a panic as warning.
                panic("async: misuse of unbounded channel, In() was closed")
            }
            ch.q = append(ch.q, e)
        case <-ch.close:
            ch.closed()
            return
        }
        for len(ch.q) > 0 {
            select {
            case ch.out <- ch.q[0]:
                ch.q[0] = nil // de-reference earlier to help GC
                ch.q = ch.q[1:]
            case e, ok := <-ch.in:
                if !ok {
                    // We don't want the input channel be accidentally closed
                    // via close() instead of Close(). If that happens, it is
                    // a misuse, do a panic as warning.
                    panic("async: misuse of unbounded channel, In() was closed")
                }
                ch.q = append(ch.q, e)
            case <-ch.close:
                ch.closed()
                return
            }
        }
        // If the remaining capacity is too small, we prefer to
        // reallocate the entire buffer.
        if cap(ch.q) < 1<<5 {
            ch.q = make([]func(), 0, 1<<10)
        }
    }
}

func (ch *UnboundedFuncChan) closed() {
    close(ch.in)
    for e := range ch.in {
        ch.q = append(ch.q, e)
    }
    for len(ch.q) > 0 {
        select {
        case ch.out <- ch.q[0]:
            ch.q[0] = nil // de-reference earlier to help GC
            ch.q = ch.q[1:]
        default:
        }
    }
    close(ch.out)
    close(ch.close)
}
*/
