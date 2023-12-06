package adagui

import (
    //"log"
    "math"
    "container/list"
    "github.com/stefan-muehlebach/gg/geom"
)

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

type StackLayout struct { }

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

type CenterLayout struct { }

func NewCenterLayout() LayoutManager {
    l := &CenterLayout{}
    return l
}

func (l *CenterLayout) Layout(childList *list.List, size geom.Point) {
    //log.Printf("CenterLayout.Layout(), size: %v", size)
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

type GridLayout struct {
    Cols   int
    orient Orientation
}

func NewColumnGridLayout(cols int) (LayoutManager) {
    l := &GridLayout{Cols: cols, orient: Horizontal}
    return l
}

func NewRowGridLayout(rows int) (LayoutManager) {
    l := &GridLayout{Cols: rows, orient: Vertical}
    return l
}

func (l *GridLayout) horizontal() bool {
    return l.orient == Horizontal
}

func (l *GridLayout) countRows(childList *list.List) (int) {
    if l.Cols < 1 {
        l.Cols = 1
    }
    count := 0
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        if child.Visible() {
            count++
        }
    }
    return int(math.Ceil(float64(count) / float64(l.Cols)))
}

func getLeading(size float64, offset int) (float64) {
    return (size + float64(pr.Size(PaddingSize))) * float64(offset)
}

func getTrailing(size float64, offset int) (float64) {
    return getLeading(size, offset+1) - pr.Size(PaddingSize)
}

func (l *GridLayout) Layout(childList *list.List, size geom.Point) {
    rows := l.countRows(childList)
    padding    := pr.Size(PaddingSize)
    padWidth   := float64(l.Cols-1) * padding
    padHeight  := float64(rows-1) * padding
    cellWidth  := float64(size.X-padWidth) / float64(l.Cols)
    cellHeight := float64(size.Y-padHeight) / float64(rows)

    if !l.horizontal() {
        padWidth, padHeight = padHeight, padWidth
        cellWidth  = float64(size.X-padWidth) / float64(rows)
        cellHeight = float64(size.Y-padHeight) / float64(l.Cols)
    }
    row, col := 0, 0
    i := 0
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        if !child.Visible() {
            continue
        }

        x1 := getLeading(cellWidth, col)
        y1 := getLeading(cellHeight, row)
        x2 := getTrailing(cellWidth, col)
        y2 := getTrailing(cellHeight, row)

        child.SetPos(geom.Point{x1, y1})
        child.SetSize(geom.Point{x2-x1, y2-y1})

        if l.horizontal() {
            if (i+1)%l.Cols == 0 {
                row++
                col=0
            } else {
                col++
            }
        } else {
            if (i+1)%l.Cols == 0 {
                col++
                row = 0
            } else {
                row++
            }
        }
        i++
    }
}

func (l *GridLayout) MinSize(childList *list.List) (geom.Point) {
    rows := l.countRows(childList)
    minSize := geom.Point{0, 0}
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        if !child.Visible() {
            continue
        }
        minSize = minSize.Max(child.MinSize())
    }

    if l.horizontal() {
        minContentSize := geom.Point{minSize.X*float64(l.Cols), minSize.Y*float64(rows)}
        return minContentSize.Add(geom.Point{pr.Size(PaddingSize)*math.Max(float64(l.Cols-1), 0.0),
                pr.Size(PaddingSize)*math.Max(float64(rows-1), 0.0)})
    }

    minContentSize := geom.Point{minSize.X*float64(rows), minSize.Y*float64(l.Cols)}
    return minContentSize.Add(geom.Point{pr.Size(PaddingSize)*math.Max(float64(rows-1), 0.0),
            pr.Size(PaddingSize)*math.Max(float64(l.Cols-1), 0.0)})
}

//----------------------------------------------------------------------------

type BorderLayout struct {
    top, bottom, left, right Node
    padding float64
}

func NewBorderLayout(top, bottom, left, right Node) LayoutManager {
    l := &BorderLayout{top, bottom, left, right, pr.Size(PaddingSize)}
    return l
}

func (l *BorderLayout) Layout(childList *list.List, size geom.Point) {
    var topSize, bottomSize, leftSize, rightSize geom.Point
    if l.top != nil && l.top.Visible() {
        topHeight := l.top.MinSize().Y
        l.top.SetSize(geom.Point{size.X, topHeight})
        l.top.SetPos(geom.Point{0, 0})
        topSize = geom.Point{size.X, topHeight+l.padding}
    }
    if l.bottom != nil && l.bottom.Visible() {
        bottomHeight := l.bottom.MinSize().Y
        l.bottom.SetSize(geom.Point{size.X, bottomHeight})
        l.bottom.SetPos(geom.Point{0, size.Y-bottomHeight})
        bottomSize = geom.Point{size.X, bottomHeight+l.padding}
    }
    if l.left != nil && l.left.Visible() {
        leftWidth := l.left.MinSize().X
        l.left.SetSize(geom.Point{leftWidth, size.Y-topSize.Y-bottomSize.Y})
        l.left.SetPos(geom.Point{0, topSize.Y})
        leftSize = geom.Point{leftWidth+l.padding, size.Y-topSize.Y-bottomSize.Y}
    }
    if l.right != nil && l.right.Visible() {
        rightWidth := l.right.MinSize().X
        l.right.SetSize(geom.Point{rightWidth, size.Y-topSize.Y-bottomSize.Y})
        l.right.SetPos(geom.Point{size.X-rightWidth, topSize.Y})
        rightSize = geom.Point{rightWidth+l.padding, size.Y-topSize.Y-bottomSize.Y}
    }

    middleSize := geom.Point{size.X-leftSize.X-rightSize.X, size.Y-topSize.Y-bottomSize.Y}
    middlePos  := geom.Point{leftSize.X, topSize.Y}
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        if !child.Visible() {
            continue
        }
        if child != l.top && child != l.bottom && child != l.left && child != l.right {
            child.SetSize(middleSize)
            child.SetPos(middlePos)
        }
    }
}

func (l *BorderLayout) MinSize(childList *list.List) (geom.Point) {
    minSize := geom.Point{0, 0}
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        if !child.Visible() {
            continue
        }
        if child != l.top && child != l.bottom && child != l.left && child != l.right {
            minSize = minSize.Max(child.MinSize())
        }
    }
    if l.left != nil && l.left.Visible() {
        leftMin := l.left.MinSize()
        minHeight := max(minSize.Y, leftMin.Y)
        minSize = geom.Point{minSize.X+leftMin.X+l.padding, minHeight}
    }
    if l.right != nil && l.right.Visible() {
        rightMin := l.right.MinSize()
        minHeight := max(minSize.Y, rightMin.Y)
        minSize = geom.Point{minSize.X+rightMin.X+l.padding, minHeight}
    }
    if l.top != nil && l.top.Visible() {
        topMin := l.top.MinSize()
        minWidth := max(minSize.X, topMin.X)
        minSize = geom.Point{minWidth, minSize.Y+topMin.Y+l.padding}
    }
    if l.bottom != nil && l.bottom.Visible() {
        bottomMin := l.bottom.MinSize()
        minWidth := max(minSize.X, bottomMin.X)
        minSize = geom.Point{minWidth, minSize.Y+bottomMin.Y+l.padding}
    }

    return minSize
}

//----------------------------------------------------------------------------

type BoxLayout struct {
    orient Orientation
    padding float64
}

func NewHBoxLayout() LayoutManager {
    l := &BoxLayout{Horizontal, pr.Size(PaddingSize)}
    return l
}

func NewVBoxLayout() LayoutManager {
    l := &BoxLayout{Vertical, pr.Size(PaddingSize)}
    return l
}

func (l *BoxLayout) isSpacer(obj Node) (bool) {
    spc, ok := obj.(*Spacer)
    if !ok {
        return false
    }
    if l.orient == Horizontal {
        return spc.ExpandHorizontal()
    }
    return spc.ExpandVertical()
}

func (l *BoxLayout) Layout(childList *list.List, size geom.Point) {
    spacers := 0
    total   := 0.0
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        if !child.Visible() {
            continue
        }
        if l.isSpacer(child) {
            spacers++
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
    extra, extraCell := 0.0, 0.0
    switch l.orient {
    case Horizontal:
        extra = size.X - total - l.padding * float64(childList.Len() - spacers - 1)
    case Vertical:
        extra = size.Y - total - l.padding * float64(childList.Len() - spacers - 1)
    }
    if spacers > 0 {
        extraCell = extra / float64(spacers)
    }
    pos := geom.Point{}
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        if !child.Visible() {
            continue
        }
        if l.isSpacer(child) {
            switch l.orient {
            case Horizontal:
                pos.X += extraCell
            case Vertical:
                pos.Y += extraCell
            }
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
        if !child.Visible() || l.isSpacer(child) {
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

func NewPaddedLayout() LayoutManager {
    return &PaddedLayout{pr.Size(PaddingSize)}
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

