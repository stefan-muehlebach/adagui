package main

import (
    "log"
    "os"
    "encoding/json"
)

type CoeffList struct {
    data []complex128
    maxFreq int
}

type FourierCoeff struct {
    freq   int
    factor complex128
}

type Complex struct {
    Re, Im float64
}

func NewComplex(c complex128) Complex {
    return Complex{real(c), imag(c)}
}

func (c *Complex) AsComplex() complex128 {
    return complex(c.Re, c.Im)
}

func NewCoeffList(d []complex128) *CoeffList {
    c := &CoeffList{}
    c.data = make([]complex128, len(d))
    copy(c.data, d)
    c.maxFreq = (len(d)-1)/2
    return c
}

func ReadCoeffList(fileName string) *CoeffList {
    c := &CoeffList{}
    in := make([]Complex, 0)
    b, err := os.ReadFile(fileName)
    if err != nil {
        log.Fatal(err)
    }
    err = json.Unmarshal(b, &in)
    if err != nil {
        log.Fatal(err)
    }
    c.data = make([]complex128, len(in))
    for i, d := range in {
        c.data[i] = d.AsComplex()
    }
    c.maxFreq = (len(in)-1)/2
    return c
}

func (c *CoeffList) Get(f int) FourierCoeff {
    if f > c.maxFreq || f < -c.maxFreq {
        log.Fatalf("frequency %d is out of bound", f)
    }
    if f >= 0 {
        return FourierCoeff{f, c.data[f]}
    } else {
        return FourierCoeff{f, c.data[len(c.data)+f]}
    }
}

