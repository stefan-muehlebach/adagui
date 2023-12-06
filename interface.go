package adagui

import (
    "container/list"
    "github.com/stefan-muehlebach/adagui/touch"
    "github.com/stefan-muehlebach/gg/geom"
    "github.com/stefan-muehlebach/gg"
)

// Alle Elemente im Scenegraph des GUI haben dieses Interface zu
// implementieren.
type Node interface {
    Wrappee() *Embed

    // Bewegt den aufrufenden Node an das Ende, resp. den Anfang der
    // Node-Liste seines Parents.
    ToBack()
    ToFront()
    IsAtFront() (bool)

    // Loescht den aufrufenden Node aus der Node-Liste seines Parents.
    Remove()

    // Setzt den Node an die angegebene Stelle, resp. retourniert die Position
    // des Widgets. Da jedes Widget seine 'Position' grundsätzlich selber
    // interpretieren kann, ist die Position nur über Methoden veränderbar,
    // resp. abrufbar.
    SetPos(p geom.Point)
    Pos() (geom.Point)

    SetSize(s geom.Point)
    Size() (geom.Point)

    SetMinSize(s geom.Point)
    MinSize() (geom.Point)

    // LocalBounds, resp. Bounds liefert das umfassende Rechteck des Nodes
    // in lokalen Koordinaten.
    LocalBounds() (geom.Rectangle)
    Bounds() (geom.Rectangle)

    // ParentBounds, resp. Rect liefert das umfassende Rechteck aus Sicht
    // des Parent-Nodes.
    ParentBounds() (geom.Rectangle)
    Rect() (geom.Rectangle)

    // Ist dieser node sichtbar?
    Visible() (bool)
    SetVisible(v bool)

    // Zeichnet das Widget im Graphik-Kontext von gc. Dabei muss sich das
    // Widget nicht um irgendwelche Koordinaten kuemmern, sondern kann davon
    // ausgehen, dass die notwendigen Transformationen durch das
    // darueberliegende Container-Widget (sofern vorhanden) gemacht wurden.
    Paint(gc *gg.Context)

    // Dient der Markierung von Nodes, bspw. um anzuzeigen, dass sie neu
    // gezeichnet werden muessen.
    Mark(m Marks)
    OnChildMarked(child Node, newMarks Marks)
    OnInputEvent(evt touch.Event)

    // Retourniert true, falls sich der Punkt pt innerhalb oder auf dem
    // Node befindet und false andernfalls.
    Contains(pt geom.Point) (bool)

    // Mit SelectTarget kann ermittelt werden, welcher Node sich an der
    // Position pt befindet. Es kann sein, dass diese Methode 'nil'
    // zurueck gibt.
    SelectTarget(pt geom.Point) (Node)

    Local2Parent(pt geom.Point) (geom.Point)
    Parent2Local(pt geom.Point) (geom.Point)
    Local2Screen(pt geom.Point) (geom.Point)
    Screen2Local(pt geom.Point) (geom.Point)

    // Methoden für die Koordinatentransformationen.
    // Die Methoden Translate, Rotate und Scale setzen jeweils unabhängig
    // voneinander eine Matrix, welche die jeweilige Transformation enthält.
    // Mit Matrix erhält man dann die Transformation, welche sich aus der 
    // Aneinanderreihung von Translate, Rotate und Scale ergibt.
    Translate(dp geom.Point)
    Rotate(a float64)
    RotateAbout(rp geom.Point, a float64)
    Scale(sx, sy float64)
    ScaleAbout(sp geom.Point, sx, sy float64)
    Matrix() (geom.Matrix)
}

// Dieses Interface implementieren zusaetzlich alle Nodes, welche als
// Container agieren koennen, d.h. eine Liste von weiteren Nodes fuehren.
type Container interface {
    Add(n ...Node)
}

// LayoutManager eben...
type LayoutManager interface {
    Layout(childList *list.List, size geom.Point)
    MinSize(childList *list.List) (geom.Point)
}

type CanvasObject interface {
    Paint(gc *gg.Context)
}

/*
type Checkable interface {
    SetCheck(checked bool)
    Checked() (bool)
    SetValue(value interface {})
    Value() (interface {})
}
*/

