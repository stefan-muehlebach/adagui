//go:build ignore
// +build ignore

package main

import (
    "log"
    "os"
    "path"
    "runtime"
    "text/template"
)

const colorPropTempl = `
func (pe *PropertyEmbed) {{ .Name }}() (color.Color) {
    return pe.prop.Color({{ .Name }})
}
func (pe *PropertyEmbed) Set{{ .Name }}(c color.Color) {
    pe.prop.SetColor({{ .Name }}, c)
}
`
const fontPropTempl = `
func (pe *PropertyEmbed) {{ .Name }}() (*opentype.Font) {
    return pe.prop.Font({{ .Name }})
}
func (pe *PropertyEmbed) Set{{ .Name }}(f *opentype.Font) {
    pe.prop.SetFont({{ .Name }}, f)
}
`
const sizePropTempl = `
func (pe *PropertyEmbed) {{ .Name }}() (float64) {
    return pe.prop.Size({{ .Name }})
}
func (pe *PropertyEmbed) Set{{ .Name }}(s float64) {
    pe.prop.SetSize({{ .Name }}, s)
}
`

type propInfo struct {
    Templ  *template.Template
    NameList []string
}

type propName struct {
    Name string
}

func main() {
    _, dirName, _, _ := runtime.Caller(0)
    filePath := path.Join(path.Dir(dirName), "propsEmbed.go")
    os.Remove(filePath)
    propFile, err := os.Create(filePath)
    if err != nil {
        log.Fatalf("Unable to open file: %v", err)
        return
    }
    defer propFile.Close()
    propFile.WriteString(`//
// THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDTI
//

package adagui

import (
    "golang.org/x/image/font/opentype"
    "github.com/stefan-muehlebach/gg/color"
)

type PropertyEmbed struct {
    prop *Properties
}

func (pe *PropertyEmbed) Init(parent *Properties) {
    pe.prop = NewProperties(parent)
}

`)

    propList := []propInfo{
        propInfo{
            Templ: template.Must(template.New("color").Parse(colorPropTempl)),
            NameList: []string{
                "Color",
                "PressedColor",
                "SelectedColor",
                "BorderColor",
                "PressedBorderColor",
                "SelectedBorderColor",
                "TextColor",
                "PressedTextColor",
                "SelectedTextColor",
                "LineColor",
                "PressedLineColor",
                "SelectedLineColor",
                "BarColor",
                "PressedBarColor",
                "BackgroundColor",
                "MenuBackgroundColor",
            },
        },
        propInfo{
            Templ: template.Must(template.New("font").Parse(fontPropTempl)),
            NameList: []string{"Font"},
        },
        propInfo{
            Templ: template.Must(template.New("size").Parse(sizePropTempl)),
            NameList: []string{
                "Width",
                "Height",
                // "Size",
                "BorderWidth",
                "PressedBorderWidth",
                "SelectedBorderWidth",
                "LineWidth",
                "InnerPadding",
                "Padding",
                "CornerRadius",
                "FontSize",
                "BarWidth",
                "CtrlWidth",
            },
        },
    }

    for _, pi := range propList {
        for _, pn := range pi.NameList {
            pr := propName{Name: pn}
            if err := pi.Templ.Execute(propFile, pr); err != nil {
                log.Fatalf("Unable to write file: %v", err)
            }
        }
    }
}

