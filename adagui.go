package adagui

import (
    "log"
)

//----------------------------------------------------------------------------

/*
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
*/

func check(err error) {
    if err != nil {
        log.Fatal(err)
    }
}


