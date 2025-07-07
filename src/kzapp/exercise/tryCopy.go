package exercise

import "fmt"

// enum like 
type CorgiType string
const (
	Wanpachi CorgiType = "Wanpachi"
	Yamper   CorgiType = "Yamper"
)

// data class like
type Character struct {
	Kind string
	Gender  string
}

// data class like
type Corgi struct {
	ID       int
	Type     CorgiType
	Name     string
	Trainner []string
	Character Character
}

func CreateCopies() (Corgi, Corgi, Corgi, Corgi) {
	makeOne := Corgi{
		ID:       1,
		Type:     Wanpachi,
		Name:     "WanJan",
		Trainner: make([]string, 0),
	}

	newOne := Corgi{
		ID:       2,
		Type:     Yamper,
		Name:     "Lulu",
		Trainner: []string{"Lisa", "Jorge"},
	}

	shallowCopy := newOne

	newSlice := make([]string, len(newOne.Trainner))
	copy(newSlice, newOne.Trainner)
	deepCopy := Corgi{
		ID:       newOne.ID,
		Name:     newOne.Name,
		Type:     newOne.Type,
		Trainner: newSlice,
	}
	logContent(&makeOne, &newOne, &shallowCopy, &deepCopy)

	return makeOne, newOne, shallowCopy, deepCopy
}

func ChangeUnknownType(newOne *Corgi, makeOne *Corgi, shallowCopy *Corgi, deepCopy *Corgi) {
	fmt.Println("----Change the type of newOne--------------------")
	newOne.Type = "Unknown"
	logContent(makeOne, newOne, shallowCopy, deepCopy)
}

func MakeNewName(newOne *Corgi, makeOne *Corgi, shallowCopy *Corgi, deepCopy *Corgi) {
	fmt.Println("----Change the name of newOne--------------------")
	newOne.Name = "Coco"
	logContent(makeOne, newOne, shallowCopy, deepCopy)
}

func LostTrainner(newOne *Corgi, makeOne *Corgi, shallowCopy *Corgi, deepCopy *Corgi) {
	fmt.Println("----Change the trainner of newOne--------------------")
	newOne.Trainner[0] = "Missing"
	newOne.Trainner = append(newOne.Trainner, "NewTrainner")
	logContent(makeOne, newOne, shallowCopy, deepCopy)
}

func ChangeFunctionType(newOneParameter Corgi, newOnePointer *Corgi) {
	fmt.Println("----Change the type of FunctionType newOne--------------------")
	newOneParameter.Type = "FunctionType"
	newOneParameter.ID = 3
	newOneParameter.Name = "FunctionType"
	newOneParameter.Trainner[0] = "FunctionType"

	fmt.Println("originOne", newOnePointer)
	fmt.Println("parameterOne", newOneParameter)
}

func logContent(makeOne *Corgi, newOne *Corgi, shallowCopy *Corgi, deepCopy *Corgi) {
	fmt.Println("makeOne", makeOne)
	fmt.Println("newOne", newOne)
	fmt.Println("shallowCopy", shallowCopy)
	fmt.Println("deepCopy", deepCopy)
}