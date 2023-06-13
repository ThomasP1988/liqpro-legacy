package socketredis

import (
	"strings"
)

func RemoveDot(dst *string, src *string) {
	*dst = strings.Replace(*src, ".", "", 1)
}

type ConsumerFormat struct {
	b        strings.Builder
	counter  int
	afterDot bool
	buff     []byte
	i        int
	dot      []byte
	zero     []byte
}

func NewConsumerFormat() *ConsumerFormat {
	csmfmt := ConsumerFormat{
		counter:  8,
		afterDot: false,
		buff:     make([]byte, 16),
		i:        0,
		dot:      []byte("."),
		zero:     []byte("0"),
	}
	csmfmt.b.Grow(16)
	return &csmfmt
}

func Format8Decimal(dst *string, src *string, csm *ConsumerFormat) {

	csm.buff = []byte(*src)
	for csm.i = 0; csm.i < len(csm.buff); csm.i++ {
		if csm.buff[csm.i] == csm.dot[0] {
			csm.afterDot = true
		} else {
			csm.b.WriteByte(csm.buff[csm.i])
			if csm.afterDot {
				csm.counter--
			}
		}
	}

	for csm.i = 0; csm.i < csm.counter; csm.i++ {
		csm.b.WriteByte(csm.zero[0])
	}
	*dst = csm.b.String()
	csm.b.Reset()
	csm.counter = 8
	csm.afterDot = false
}
