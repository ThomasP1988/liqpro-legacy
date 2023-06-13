package socketredis

import (
	"fmt"
	"testing"
)

// function to Benchmark ReturnGeeks()

func BenchmarkFormat8Decimal(b *testing.B) {

	csm := NewConsumerFormat()

	dst := ""
	src := "0.345430"
	for i := 0; i < b.N; i++ {
		Format8Decimal(&dst, &src, csm)
		csm.b.Grow(16)
	}
}

func ExampleFormat8Decimal() {

	csm := NewConsumerFormat()

	dst := ""
	src := "0.345430"
	Format8Decimal(&dst, &src, csm)
	fmt.Println(dst)
	// Output: 034543000
}
