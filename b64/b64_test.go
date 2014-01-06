package b64_test

import (
	"github.com/pebbe/b64"

	"fmt"
	"math"
)

func ExampleEncode() {
	fmt.Println(b64.Encode(0))
	fmt.Println(b64.Encode(100))
	fmt.Println(b64.Encode(256))
	fmt.Println(b64.Encode(1000))
	fmt.Println(b64.Encode(10000))
	fmt.Println(b64.Encode(math.MaxUint32))

	// Output:
	// A
	// Bk
	// EA
	// Po
	// CcQ
	// D/////
}

func ExampleDecode() {
	var v uint32
	var e error

	v, e = b64.Decode("A")
	fmt.Println(v, e)

	v, e = b64.Decode("Bk")
	fmt.Println(v, e)

	v, e = b64.Decode("EA")
	fmt.Println(v, e)

	v, e = b64.Decode("Po")
	fmt.Println(v, e)

	v, e = b64.Decode("CcQ")
	fmt.Println(v, e)

	v, e = b64.Decode("D/////")
	fmt.Println(v, e)

	v, e = b64.Decode("D//@//")
	fmt.Println(v, e)

	// Output:
	// 0 <nil>
	// 100 <nil>
	// 256 <nil>
	// 1000 <nil>
	// 10000 <nil>
	// 4294967295 <nil>
	// 0 Illegal character in base64 value: @
}
