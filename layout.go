package adagui

import (
//    "log"
    "math"
    "container/list"
    "github.com/stefan-muehlebach/gg/geom"
)

// Mit dem NullLayout werden die verwalteten Nodes per SetPos platziert und
// werden durch den Container nicht mehr weiter verwaltet. MinSize liefert
// die maximale Grösse aller verwalteten Nodes.
type NullLayout struct { }

func (l *NullLayout) Layout(childList *list.List, size geom.Point) {
    //stackLevel.Inc()
    //defer stackLevel.Dec()
    //log.Printf("NullLayout.Layout(cl, %v)", size)
}

func (l *NullLayout) MinSize(childList *list.List) (geom.Point) {
    //stackLevel.Inc()
    //defer stackLevel.Dec()
    //log.Printf("NullLayout.MinSize()")

    minSize := geom.Point{}
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        if !child.Visible() {
            continue
        }
        minSize = minSize.Max(child.MinSize())
    }
    return minSize
}

// Mit PaddedLayout kann sinnvollerweise nur ein Node verwaltet werden, der
// mit einem konfigurierbaren Abstand auf die ganze Grösse des Containers
// expandiert wird.
type PaddedLayout struct {
    pad [4]float64
}

// Mit den variablen Parametern pads können die Ränder definiert werden.
// Dabei gilt:
//   -       : verwende das Property 'InnerPadding'
//   a       : verwende a für alle Ränder
//   a,b     : verwende a für die horizontalen und b für die vertikalen
//             Ränder
//   a,b,c   : verwende a für links, b für oben und unten c für rechts
//   a,b,c,d : (dito) und d für den unteren Rand.
func NewPaddedLayout(pads... float64) *PaddedLayout {
    var l, t, r, b float64
    switch len(pads) {
    case 0:
        p := DefProps.Size(InnerPadding)
        l, t, r, b = p, p, p, p
    case 1:
        l, t, r, b = pads[0], pads[0], pads[0], pads[0]
    case 2:
        l, t, r, b = pads[0], pads[1], pads[0], pads[1]
    case 3:
        l, t, r, b = pads[0], pads[1], pads[2], pads[1]
    case 4:
        l, t, r, b = pads[0], pads[1], pads[2], pads[3]
    }
    return &PaddedLayout{[4]float64{l, t, r, b}}
}

func (l *PaddedLayout) Layout(childList *list.List, size geom.Point) {
    //stackLevel.Inc()
    //defer stackLevel.Dec()
    //log.Printf("PaddedLayout.Layout(cl, %v)", size)

    pos := geom.Point{l.pad[Left], l.pad[Top]}
    siz := geom.Point{size.X-(l.pad[Left]+l.pad[Right]),
            size.Y-(l.pad[Top]+l.pad[Bottom])}
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        child.SetSize(siz)
        child.SetPos(pos)
    }
}

func (l *PaddedLayout) MinSize(childList *list.List) (geom.Point) {
    //stackLevel.Inc()
    //defer stackLevel.Dec()
    //log.Printf("PaddedLayout.MinSize()")

    minSize := geom.Point{}
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*Embed).Wrapper
        if !child.Visible() {
            continue
        }
        minSize = minSize.Max(child.MinSize())
    }
    return minSize.Add(geom.Point{l.pad[Left]+l.pad[Right],
        l.pad[Top]+l.pad[Bottom]})
}

// Wie beim PaddedLayout wird man mit Stack- oder Max-Layout sinnvollerweise
// auch nur ein einziges Child verwalten können. Dieses wird auf die ganze
// Grösse des Containers ausgedehnt - einfach ohne Ränder.
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

// CenterLayout zentiert alle Kinder, verändert ihre Grössen jedoch nicht.
// Auch hier gilt das gleiche wie oben: nur für ein Kind zu verwenden.
type CenterLayout struct { }

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

// GridLayout ordnet die Kinder in Spalten oder Zeilen an - je nachdem ob
// NewColumn... oder NewRow... aufgerufen wird.
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
    return (size + float64(DefProps.Size(InnerPadding))) * float64(offset)
}

func getTrailing(size float64, offset int) (float64) {
    return getLeading(size, offset+1) - DefProps.Size(InnerPadding)
}

func (l *GridLayout) Layout(childList *list.List, size geom.Point) {
    rows := l.countRows(childList)
    padding    := DefProps.Size(InnerPadding)
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

    pad := DefProps.Size(InnerPadding)
    if l.horizontal() {
        minContentSize := geom.Point{minSize.X*float64(l.Cols),
                minSize.Y*float64(rows)}
        return minContentSize.Add(geom.Point{pad*math.Max(float64(l.Cols-1), 0),
                pad*math.Max(float64(rows-1), 0.0)})
    }

    minContentSize := geom.Point{minSize.X*float64(rows),
            minSize.Y*float64(l.Cols)}
    return minContentSize.Add(geom.Point{pad*math.Max(float64(rows-1), 0.0),
            pad*math.Max(float64(l.Cols-1), 0.0)})
}

// Das BorderLayout ist ein rechtes Monster. Damit können Fenster mit
// Titel- oder Menuzeile, rechter und linker Randspalte sowie Fusszeile
// verwaltet werden. Die einzelnen Zeilen können auch leer gelassen werden
// (verwende nil).
type BorderLayout struct {
    top, bottom, left, right Node
    padding float64
}

func NewBorderLayout(top, bottom, left, right Node) LayoutManager {
    l := &BorderLayout{top, bottom, left, right, DefProps.Size(InnerPadding)}
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

// Die BoxLayouts gibt es in zwei Varianten: horizontal oder vertikal.
// Alle verwalteten Kinder werden neben oder untereinander angeordnet und
// auf die gleiche Höhe, resp. gleiche Breite getrimmt. Bei diesem Layout
// kommen auch die Spacer-Widgets zum Einsatz: sie dehnen sich auf die
// maximale Breite, resp. Höhe aus, haben aber sonst keinen Inhalt oder
// Einfluss.
type BoxLayout struct {
    orient Orientation
    padding float64
}

func NewHBoxLayout(pads... float64) *BoxLayout {
    pad := DefProps.Size(InnerPadding)
    if len(pads) > 0 {
        pad = pads[0]
    }
    return &BoxLayout{Horizontal, pad}
}

func NewVBoxLayout(pads... float64) *BoxLayout {
    pad := DefProps.Size(InnerPadding)
    if len(pads) > 0 {
        pad = pads[0]
    }
    return &BoxLayout{Vertical, pad}
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
        extra = size.X - total - l.padding * float64(childList.Len() -
                spacers - 1)
    case Vertical:
        extra = size.Y - total - l.padding * float64(childList.Len() -
                spacers - 1)
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

// Nimmt den verfügbaren Platz (vertikal oder horizontal) in Box-Layouts
// ein. Ist zwar ein Widget, passt aber irgendwie besser zum Layout-Zeugs.
type Spacer struct {
    LeafEmbed
    FixHorizontal, FixVertical bool
}

func NewSpacer() (*Spacer) {
    s := &Spacer{}
    s.Wrapper = s
    s.Init(DefProps)
    return s
}

func (s *Spacer) ExpandHorizontal() (bool) {
    return !s.FixHorizontal
}

func (s *Spacer) ExpandVertical() (bool) {
    return !s.FixVertical
}


