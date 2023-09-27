module mju.net/adagui

go 1.21.0

replace mju.net/gg => ../gg

replace mju.net/geom => ../geom

replace mju.net/adatft => ../adatft

replace mju.net/utils => ../utils

require (
	golang.org/x/image v0.11.0
	mju.net/adatft v0.0.0-00010101000000-000000000000
	mju.net/geom v0.0.0-00010101000000-000000000000
	mju.net/gg v0.0.0-00010101000000-000000000000
	mju.net/utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/text v0.12.0 // indirect
	periph.io/x/conn/v3 v3.7.0 // indirect
	periph.io/x/host/v3 v3.8.2 // indirect
)
