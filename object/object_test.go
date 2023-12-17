package object

import (
	"fmt"
	"testing"
)

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "Hello Now"}
	diff2 := &String{Value: "Hello Now"}
	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with the same Value have different hashes")
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with the same Value have different hashes")
	}
	if hello1.HashKey() == diff1.HashKey() {
		fmt.Println(hello1.HashKey())
		fmt.Println(diff2.HashKey())
		t.Errorf("strings with the different Value have same hashes")
	}
}
