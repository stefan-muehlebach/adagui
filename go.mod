module github.com/stefan-muehlebach/adagui

go 1.23.2

replace github.com/stefan-muehlebach/adatft => ../adatft

replace github.com/stefan-muehlebach/gg => ../gg

replace github.com/stefan-muehlebach/ledgrid => ../ledgrid

replace github.com/stefan-muehlebach/mandel => ../mandel

require (
	github.com/cpmech/gosl v1.2.12
	github.com/stefan-muehlebach/adatft v1.2.1
	github.com/stefan-muehlebach/gg v1.3.4
	github.com/stefan-muehlebach/mandel v0.0.0-00010101000000-000000000000
	golang.org/x/image v0.21.0
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/text v0.19.0 // indirect
	periph.io/x/conn/v3 v3.7.1 // indirect
	periph.io/x/host/v3 v3.8.2 // indirect
)
