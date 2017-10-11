package oas2

import (
	"fmt"
	"os"
)

func ExampleID_String() {
	opID := OperationID("addPet")

	fmt.Fprint(os.Stdout, opID.String())

	// Output:
	// addPet
}
