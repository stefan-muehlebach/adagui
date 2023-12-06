package adagui

import (
    "image/color"
    "time"
    "github.com/stefan-muehlebach/gg/geom"
)

// Als nächstes folgen die Typen, mit welchen man Animation in die Bude
// bringen kann.
type AnimationCurve func(float64) float64

const (
    AnimationRepeatForever = -1
    DurationStandard       = time.Millisecond * 300
    DurationShort          = time.Millisecond * 150
)

var (
    AnimationEaseInOut = animationEaseInOut
    AnimationEaseIn    = animationEaseIn   
    AnimationEaseOut   = animationEaseOut
    AnimationLinear    = animationLinear
)

type Animation struct {
    AutoReverse bool
    Curve       AnimationCurve
    Duration    time.Duration
    RepeatCount int
    Tick        func(float64)
}

func NewAnimation(d time.Duration, fn func(float64)) (*Animation) {
    return &Animation{Duration: d, Tick: fn}
}

func (a *Animation) Start() {
    CurrentScreen().StartAnimation(a)
}

func (a *Animation) Stop() {
    CurrentScreen().StopAnimation(a)
}

func animationEaseIn(val float64) (float64) {
    return val*val
}

func animationEaseInOut(val float64) (float64) {
    if val <= 0.5 {
        return val*val*2
    }
    return -1 + (4-val*2)*val
}

func animationEaseOut(val float64) (float64) {
    return val*(2-val)
}

func animationLinear(val float64) (float64) {
    return val
}

func NewColorAnimation(start, stop color.Color, d time.Duration,
        fn func(color.Color)) (*Animation) {
    r1, g1, b1, a1 := start.RGBA()
    r2, g2, b2, a2 := stop.RGBA()

    rStart := int(r1 >> 8)
    gStart := int(g1 >> 8)
    bStart := int(b1 >> 8)
    aStart := int(a1 >> 8)
    rDelta := float64(int(r2 >> 8) - rStart)
    gDelta := float64(int(g2 >> 8) - gStart)
    bDelta := float64(int(b2 >> 8) - bStart)
    aDelta := float64(int(a2 >> 8) - aStart)

    return &Animation{
        Duration: d,
        Tick: func(done float64) {
                  fn(color.NRGBA{
                      R: scaleChannel(rStart, rDelta, done),
                      G: scaleChannel(gStart, gDelta, done),
                      B: scaleChannel(bStart, bDelta, done),
                      A: scaleChannel(aStart, aDelta, done),
                  })}}
}

func NewPositionAnimation(start, stop geom.Point, d time.Duration,
        fn func(geom.Point)) (*Animation) {
    xDelta := float64(stop.X - start.X)
    yDelta := float64(stop.Y - start.Y)

    return &Animation{
        Duration: d,
        Tick: func(done float64) {
                  fn(geom.Point{scaleVal(start.X, xDelta, done),
                                scaleVal(start.Y, yDelta, done)})
        }}
}

func NewSizeAnimation(start, stop geom.Point, d time.Duration,
        fn func(geom.Point)) (*Animation) {
    widthDelta  := float64(stop.X - start.X)
    heightDelta := float64(stop.Y - start.Y)

    return &Animation{
        Duration: d,
        Tick: func(done float64) {
                  fn(geom.Point{scaleVal(start.X, widthDelta, done),
                                scaleVal(start.Y, heightDelta, done)})
        }}
}

func scaleChannel(start int, diff, done float64) (uint8) {
    return uint8(start + int(diff*done))
}

func scaleVal(start float64, delta, done float64) (float64) {
    return start + delta*done
}

