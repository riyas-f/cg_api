package database

import (
	"testing"
)

type AnotherStruct struct {
	Phone             string               `json:"some_data" db:"hello_world"`
	PhoneProvider     string               `json:"hello" db:"provider"`
	DontIndexMe       string               `json:"he" db:"-"`
	HaveAnotherStruct AnotherAnotherStruct `json:"-" db:"-"`
}

type AnotherAnotherStruct struct {
	Geh   string `db:"-"`
	Value string `db:"index_me_pls"`
}

type Data struct {
	Name       string        `db:"name"`
	SomeStruct AnotherStruct //won't get index on getField
}

func TestResolveRecursive(t *testing.T) {
	Data := &Data{
		Name: "hello",
		SomeStruct: AnotherStruct{
			Phone:         "Now i am exist",
			PhoneProvider: "Hehe",
			DontIndexMe:   "This is a secret, yes",
			HaveAnotherStruct: AnotherAnotherStruct{
				Geh:   "WHy are you geh",
				Value: "This is not a secret yeah",
			},
		},
	}

	names, value, _ := getField(Data, true)
	t.Log(names)
	t.Log(value)
}
