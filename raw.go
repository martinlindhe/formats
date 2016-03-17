package formats

import (
	"errors"
	"fmt"
	"os"
)

type probeFunc func(f *os.File) bool

func (pf probeFunc) Probe(f *os.File) bool {
	return pf(f)
}

func probeX(f *os.File) {

	probes := []func(*os.File) bool{
		psxExeProbe, psxBiosProbe,
	}

	for _, p := range probes {
		if p(f) {
			fmt.Printf("yay, %v", p)
		}
	}
	/*
		if ProbeFunc(psxExeProbe) {
			fmt.Printf("yay psx-exe\n")
		} else if ProbeFunc(psxBiosProbe) {
			fmt.Printf("yay psx-bios\n")
		} else {
			fmt.Printf("not recognized\n")
		}
	*/
}

func readBytes(f *os.File, want int) (data []byte, err error) {
	data = make([]byte, want)
	got, err := f.Read(data)
	if err != nil {
		return nil, err
	}
	if got != want {
		return nil, errors.New("not enough data")
	}
	return data, nil
}
