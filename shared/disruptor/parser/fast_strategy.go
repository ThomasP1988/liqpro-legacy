package socketredis

import "time"

type FastWaitStrategy struct{}

// TODO: time should be decreased when in production
func NewFastWaitStrategy() FastWaitStrategy    { return FastWaitStrategy{} }
func (this FastWaitStrategy) Gate(count int64) { time.Sleep(time.Nanosecond) }
func (this FastWaitStrategy) Idle(count int64) { time.Sleep(time.Microsecond * 50) }
