package adagui

import (
    "container/list"
    "mju.net/geom"
)

//----------------------------------------------------------------------------

func max(a, b float64) (float64) {
    if a > b {
        return a
    } else {
        return b
    }
}

func min(a, b float64) (float64) {
    if a < b {
        return a
    } else {
        return b
    }
}

//----------------------------------------------------------------------------

type NullLayout struct { }

func (l *NullLayout) Layout(childList *list.List, size geom.Point) { }

func (l *NullLayout) MinSize(childList *list.List) (geom.Point) {
    //log.Printf("NullLayout.MinSize()")
    minSize := geom.Point{1.0, 1.0}
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        minSize = minSize.Max(child.MinSize())
    }
    return minSize
}

//----------------------------------------------------------------------------

type StackLayout struct {
}

func NewStackLayout() LayoutManager {
    l := &StackLayout{}
    return l
}

func NewMaxLayout() LayoutManager {
    return NewStackLayout()
}

func (l *StackLayout) Layout(childList *list.List, size geom.Point) {
    pos := geom.Point{0.0, 0.0}
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        child.SetSize(size)
        child.SetPos(pos)
    }
}

func (l *StackLayout) MinSize(childList *list.List) (geom.Point) {
    //log.Printf("StackLayout.MinSize()")
    minSize := geom.Point{0.0, 0.0}
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        if !child.Visible() {
            continue
        }
        minSize = minSize.Max(child.MinSize())
    }
    return minSize
}

//----------------------------------------------------------------------------

type CenterLayout struct {
}

func NewCenterLayout() LayoutManager {
    l := &CenterLayout{}
    return l
}

func (l *CenterLayout) Layout(childList *list.List, size geom.Point) {
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        childMin := child.MinSize()
        dp := size.Sub(childMin).Mul(0.5)
        child.SetSize(childMin)
        child.SetPos(child.Pos().Add(dp))
    }
}

func (l *CenterLayout) MinSize(childList *list.List) (geom.Point) {
    //log.Printf("CenterLayout.MinSize()")
    minSize := geom.Point{0.0, 0.0}
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        if !child.Visible() {
            continue
        }
        minSize = minSize.Max(child.MinSize())
    }
    return minSize
}

//----------------------------------------------------------------------------

var (
    boxLayPadding = 5.0
)

type BoxLayout struct {
    orient Orientation
    padding float64
}

func NewHBoxLayout(padding float64) LayoutManager {
    l := &BoxLayout{Horizontal, padding}
    return l
}

func NewVBoxLayout(padding float64) LayoutManager {
    l := &BoxLayout{Vertical, padding}
    return l
}

func (l *BoxLayout) Layout(childList *list.List, size geom.Point) {
    total := 0.0
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        if !child.Visible() {
            continue
        }
        childSize := child.MinSize()
        switch l.orient {
        case Horizontal:
            total += childSize.X
        case Vertical:
            total += childSize.Y
        }
    }
    total += l.padding * float64(childList.Len() - 1)
    pos := geom.Point{}
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        if !child.Visible() {
            continue
        }
        child.SetPos(pos)
        switch l.orient {
        case Horizontal:
            width := child.MinSize().X
            pos.X += width + l.padding
            child.SetSize(geom.Point{width, size.Y})
        case Vertical:
            height := child.MinSize().Y
            pos.Y += height + l.padding
            child.SetSize(geom.Point{size.X, height})
        }
    }
}

func (l *BoxLayout) MinSize(childList *list.List) (geom.Point) {
    //log.Printf("BoxLayout.MinSize()")
    minSize := geom.Point{}
    addPadding := false
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        if !child.Visible() {
            continue
        }
        childSize := child.MinSize()
        switch l.orient {
        case Horizontal:
            minSize.Y  = max(minSize.Y, childSize.Y)
            minSize.X += childSize.X
            if addPadding {
                minSize.X += l.padding
            }
        case Vertical:
            minSize.X  = max(minSize.X, childSize.X)
            minSize.Y += childSize.Y
            if addPadding {
                minSize.Y += l.padding
            }
        }
        addPadding = true
    }
    return minSize
}

//----------------------------------------------------------------------------

type PaddedLayout struct {
    padding float64
}

func NewPaddedLayout(padding float64) LayoutManager {
    return &PaddedLayout{padding}
}

func (l *PaddedLayout) Layout(childList *list.List, size geom.Point) {
    pos := geom.Point{l.padding, l.padding}
    siz := geom.Point{size.X-2*l.padding, size.Y-2*l.padding}
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        child.SetSize(siz)
        child.SetPos(pos)
    }
}

func (l *PaddedLayout) MinSize(childList *list.List) (geom.Point) {
    minSize := geom.Point{}
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        if !child.Visible() {
            continue
        }
        minSize = minSize.Max(child.MinSize())
    }
    return minSize.Add(geom.Point{2*l.padding, 2*l.padding})
}

