package adagui

import (
    "image/color"
    "github.com/stefan-muehlebach/adagui/touch"
    "github.com/stefan-muehlebach/gg/colornames"
    "github.com/stefan-muehlebach/gg/geom"
    "github.com/stefan-muehlebach/gg"
)

//-----------------------------------------------------------------------------

// GUI daten sowohl fuer PageButton als auch Drawer
var (
    flapWidth            = 22.0
    flapHeight           = 50.0
    flapInset            =  9.0
    flapRectRad          =  6.0
    flapFillColor        = colornames.Gray.Alpha(0.5)
    flapPushedFillColor  = colornames.Gray.Alpha(0.5)
    flapArrowColor       = pr.Color(Color)
    flapPushedArrowColor = pr.Color(PressedColor)
    flapArrowWidth       =  8.0
    flapSize = []geom.Point{
	    geom.Point{flapWidth, flapHeight},
	    geom.Point{flapHeight, flapWidth},
	    geom.Point{flapWidth, flapHeight},
	    geom.Point{flapHeight, flapWidth},
    }
    flapRectInsets = [][]geom.Point{
        {
            geom.Point{-flapRectRad, 0.0},
            geom.Point{},
        },
        {
            geom.Point{0.0, -flapRectRad},
            geom.Point{},
        },
        {
            geom.Point{},
            geom.Point{flapRectRad, 0.0},
        },
        {
            geom.Point{},
            geom.Point{0.0, flapRectRad},
        },
    }
    pgBtnFillColor = flapFillColor
    drwFillColor   = colornames.LightGray
)

func DrawArrow(gc *gg.Context, dst geom.Rectangle, border Border) {
    switch border {
    case Left:
        gc.MoveTo(dst.Max.X, dst.Min.Y)
        gc.LineTo(dst.Min.X, 0.5*(dst.Min.Y+dst.Max.Y))
        gc.LineTo(dst.Max.X, dst.Max.Y)
    case Top:
        gc.MoveTo(dst.Min.X, dst.Max.Y)
        gc.LineTo(0.5*(dst.Min.X+dst.Max.X), dst.Min.Y)
        gc.LineTo(dst.Max.X, dst.Max.Y)
    case Right:
        gc.MoveTo(dst.Min.X, dst.Min.Y)
        gc.LineTo(dst.Max.X, 0.5*(dst.Min.Y+dst.Max.Y))
        gc.LineTo(dst.Min.X, dst.Max.Y)
    case Bottom:
        gc.MoveTo(dst.Min.X, dst.Min.Y)
        gc.LineTo(0.5*(dst.Min.X+dst.Max.X), dst.Max.Y)
        gc.LineTo(dst.Max.X, dst.Min.Y)
    }
}

//-----------------------------------------------------------------------------

// PageButton dienen vorallem fuer den Wechsel zwischen den Windows, koennen
// aber auch fuer anderes verwendet werden.
type PageButton struct {
    LeafEmbed
    border Border
    pushed bool
}

func NewPageButton(border Border) (*PageButton) {
    b := &PageButton{}
    b.Wrapper = b
    b.Init(nil)
    b.SetMinSize(geom.Point{flapWidth, flapHeight})
    b.border = border
    b.pushed = false
    return b
}

/*
func (b *PageButton) SetPos(pt geom.Point) {
    switch b.pos {
    case Left:
        pt.X = 0.0
        pt.Y -= 0.5*flapHeight
    case Top:
        pt.X -= 0.5*flapHeight
        pt.Y = 0.0
    case Right:
        pt.X = 298.0
        pt.Y -= 0.5*flapHeight
    case Bottom:
        pt.X -= 0.5*flapHeight
        pt.Y = 218.0
    }
    b.Wrappee().SetPos(pt)

    b.ExtRect = b.Rect().Sub(b.Rect().Min)
    b.ExtRect = geom.Rectangle{
        b.ExtRect.Min.Add(flapRectInsets[b.pos][0]),
        b.ExtRect.Max.Add(flapRectInsets[b.pos][1]),
    }
}
*/

func (b *PageButton) Paint(gc *gg.Context) {
    //log.Printf("PageButton.Paint()")
    b.Marks.UnmarkNeedsPaint()
    //gc.Push()
    //gc.Translate(b.Rect().Min.AsCoord())
    gc.DrawRoundedRectangle(-flapRectRad, 0.0, b.Size().X+flapRectRad,
            b.Size().Y, flapRectRad)
    gc.SetFillColor(flapFillColor)
    gc.Fill()

    if b.pushed {
        gc.SetStrokeColor(flapPushedArrowColor)
    } else {
        gc.SetStrokeColor(flapArrowColor)
    }
    gc.SetStrokeWidth(flapArrowWidth)
    DrawArrow(gc, b.LocalBounds().Inset(flapInset, flapInset), b.border)
    gc.Stroke()
    //gc.Pop()
}

func (b *PageButton) OnInputEvent(evt touch.Event) {
    //log.Printf("%T: %v", b, evt)
    switch evt.Type {
    case touch.TypePress, touch.TypeEnter:
        b.pushed = true
        b.Mark(MarkNeedsPaint)
    case touch.TypeLeave:
        b.pushed = false
        b.Mark(MarkNeedsPaint)
    case touch.TypeRelease:
        b.pushed = false
        b.Mark(MarkNeedsPaint)
    }
    b.CallTouchFunc(evt)
}

// Der Drawer (engl. Schublade) kann eine Reihe von weiteren Widgets aufnehmen
// und laesst sich bei Nichtbedarf am Rand des Bildschirms auf ein kleines
// Icon zusammenklappen. Eine Antwort auf den beschraenkten Platz des Adafruit
// TFT-Bildschirm.
var (
    drwSizeChange = [][]geom.Point{
        {
            geom.Point{},
            geom.Point{100.0, 0.0},
        },
	    {
            geom.Point{},
            geom.Point{0.0, 100.0},
        },
        {
            geom.Point{-100.0, 0.0},
            geom.Point{},
        },
        {
            geom.Point{0.0, -100.0},
            geom.Point{},
        },
    }
)

type Drawer struct {
    ContainerEmbed
    pos Border
    FillColor, ActiveColor color.Color
    pushed bool
    handle geom.Rectangle
    isOpen bool
    ExtRect geom.Rectangle
}

func NewDrawer(pos Border) (*Drawer) {
    d := &Drawer{}
    d.Wrapper = d
    d.Init(nil)
    d.pos = pos
    d.FillColor = pr.Color(Color)
    d.ActiveColor = pr.Color(ActiveColor)
    d.pushed = false
    d.isOpen = false
    d.SetSize(flapSize[d.pos])
    return d
}

func (d *Drawer) SetPos(pt geom.Point) {
    switch d.pos {
    case Left:
        pt.X = 0.0
    case Top:
        pt.Y = 0.0
    case Right:
        pt.X = 298.0
    case Bottom:
        pt.Y = 218.0
    }
    d.Wrappee().SetPos(pt)

    d.ExtRect = d.Rect().Sub(d.Rect().Min)
    d.ExtRect = geom.Rectangle{
        d.ExtRect.Min.Add(flapRectInsets[d.pos][0]),
        d.ExtRect.Max.Add(flapRectInsets[d.pos][1]),
    }
}

func (d *Drawer) Paint(gc *gg.Context) {
    //log.Printf("Drawer.Paint()")
    d.Marks.UnmarkNeedsPaint()
    gc.Push()
    gc.Translate(d.Rect().Min.AsCoord())
    gc.DrawRoundedRectangle(d.ExtRect.Min.X, d.ExtRect.Min.Y,
            d.ExtRect.Dx(), d.ExtRect.Dy(), flapRectRad)
    //log.Printf("Drawer.Paint():")
    if d.pushed {
        gc.SetFillColor(d.ActiveColor)
    } else {
        gc.SetFillColor(d.FillColor)
    }
    gc.Fill()

    gc.SetStrokeColor(flapArrowColor)
    gc.SetStrokeWidth(flapArrowWidth)
    DrawArrow(gc, d.ExtRect.Inset(flapInset, flapInset), (d.pos+2)%4)
    gc.Stroke()

    gc.DrawRectangle(d.ExtRect.AsCoord())
    gc.Clip()
    //d.ContainerEmbed.Paint(gc)
    gc.ResetClip()
    gc.Pop()
}

func (d *Drawer) OnInputEvent(evt touch.Event) {
    //log.Printf("Drawer.OnInputEvent(): %T, %v", d, evt)
    switch evt.Type {
    case touch.TypePress, touch.TypeEnter:
        d.pushed = true
        d.Mark(MarkNeedsPaint)
    case touch.TypeLeave, touch.TypeRelease:
        d.pushed = false
        d.Mark(MarkNeedsPaint)
    case touch.TypeTap:
        if d.isOpen {
            d.Close()
        } else {
            d.Open()
        }
    }
}

func (d *Drawer) IsOpen() (bool) {
    return d.isOpen
}

func (d *Drawer) Open() {
    if d.isOpen {
        return
    }
    d.isOpen = true
    d.SetPos(d.Pos().Add(drwSizeChange[d.pos][0]))
    d.SetSize(d.Size().Add(drwSizeChange[d.pos][1]))
    //d.Rect.Min = d.Rect.Min.Add(drwSizeChange[d.pos][0])
    //d.Rect.Max = d.Rect.Max.Add(drwSizeChange[d.pos][1])
    d.Mark(MarkNeedsPaint)
}


func (d *Drawer) Close() {
    if !d.isOpen {
        return
    }
    d.isOpen = false
    d.SetPos(d.Pos().Sub(drwSizeChange[d.pos][0]))
    d.SetSize(d.Size().Sub(drwSizeChange[d.pos][1]))
    //d.Rect.Min = d.Rect.Min.Sub(drwSizeChange[d.pos][0])
    //d.Rect.Max = d.Rect.Max.Sub(drwSizeChange[d.pos][1])
    d.Mark(MarkNeedsPaint)
}

