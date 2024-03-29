module github.com/stefan-muehlebach/adagui

go 1.22.1

replace github.com/stefan-muehlebach/adatft => ../adatft

replace github.com/stefan-muehlebach/gg => ../gg

require (
	github.com/stefan-muehlebach/adatft v0.0.0-00010101000000-000000000000
	github.com/stefan-muehlebach/gg v1.2.0
	golang.org/x/image v0.15.0
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/text v0.14.0 // indirect
	periph.io/x/conn/v3 v3.7.0 // indirect
	periph.io/x/host/v3 v3.8.2 // indirect
)
