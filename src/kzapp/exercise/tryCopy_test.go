package exercise

import (
	"fmt"
	"testing"
)

func TestChangeFunctionType(t *testing.T) {
	_, newOne, shallowCopy, _ := CreateCopies()

	ChangeFunctionType(newOne, &newOne)

	if newOne.Type == "FunctionType" {
		t.Errorf("newOne type should not be changed")
	}
	if newOne.Name == "FunctionType" {
		t.Errorf("newOne type should not be changed")
	}
	if newOne.Trainner[0] != "FunctionType" {
		t.Errorf("newOne trainner is not changed")
	}

	// any create new instance it will got a new address but pointer to the same value
	if compareAddress(&newOne, &shallowCopy, false, true, true) {
		fmt.Println("shallowCopy address is not the same as newOne, but slice item is the same")
	} else {
		t.Errorf("shallowCopy slice item is different from newOne")
	}
}

func TestChangeUnknownType(t *testing.T) {
	makeOne, newOne, shallowCopy, deepCopy := CreateCopies()

	ChangeUnknownType(&newOne, &makeOne, &shallowCopy, &deepCopy)

	if compareAddress(&newOne, &shallowCopy, false, true, true) &&
		newOne.Type != shallowCopy.Type {
		fmt.Println("shallowCopy address is not the same as newOne, and type is changed")
	} else {
		t.Errorf("newOne type is not changed")
	}
}

func TestMakeNewName(t *testing.T) {
	makeOne, newOne, shallowCopy, deepCopy := CreateCopies()

	MakeNewName(&newOne, &makeOne, &shallowCopy, &deepCopy)

	if compareAddress(&newOne, &shallowCopy, false, true, true) &&
		newOne.Name != shallowCopy.Name {
		fmt.Println("shallowCopy address is not the same as newOne, and name is changed")
	} else {
		t.Errorf("newOne name is not changed")
	}
}

func TestLostTrainner(t *testing.T) {
	makeOne, newOne, shallowCopy, deepCopy := CreateCopies()

	LostTrainner(&newOne, &makeOne, &shallowCopy, &deepCopy)

	if compareAddress(&newOne, &shallowCopy, false, false, true) {
		fmt.Println("shallowCopy address and sliced item is not the same as newOne, but sliced item content is same")
	} else {
		t.Errorf("shallowCopy slice item content is different from newOne")
	}
}

func compareAddress(src, dst *Corgi, isSameAdress, isSameSliceItemAddress, isSameSliceItemContent bool) bool {
	isSameTypeAdreess := &src.Type == &dst.Type 
	isSameNameAdreess := &src.Name == &dst.Name 
	isSameIDAdreess := &src.ID == &dst.ID 
	isSameSliceAdreess := &src.Trainner == &dst.Trainner
	isSameAll := src == dst

	isSameFirstAddress := &src.Trainner[0] == &dst.Trainner[0]
	isSameFirstItem := src.Trainner[0] == dst.Trainner[0]
	fmt.Println("")
	fmt.Println("isSameTypeAdreess", isSameTypeAdreess)
	fmt.Println("isSameNameAdreess", isSameNameAdreess)
	fmt.Println("isSameIDAdreess", isSameIDAdreess)
	fmt.Println("isSameSliceAdreess", isSameSliceAdreess)
	fmt.Println("isSameAll", isSameAll)
	fmt.Println("isSameFirstAddress", isSameFirstAddress)
	fmt.Println("isSameFirstItem", isSameFirstItem)
	fmt.Println("")
	return (isSameTypeAdreess && isSameNameAdreess && isSameIDAdreess && isSameSliceAdreess && isSameAll) == isSameAdress && isSameFirstAddress == isSameSliceItemAddress && isSameFirstItem == isSameSliceItemContent
}