package main

import (
	"fmt"

	"github.com/object88/cprofile/test/ipsum/bacon"
)

func main() {
	message := bacon.GenerateBaconIpsum()
	fmt.Printf(message)
}
