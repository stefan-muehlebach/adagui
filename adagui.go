//go:generate go generate ./props

package adagui

import (
    "log"
)

func check(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

