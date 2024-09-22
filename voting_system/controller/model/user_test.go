package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUser(t *testing.T) {
	u := &User{ID: 1}
	u.ID = 1
	code, _ := json.Marshal(u)
	fmt.Println("User", string(code), u.ID)
}

func TestCharacter(t *testing.T) {
	c := &Character{}
	c.ID = 2
	code, _ := json.Marshal(c)
	fmt.Println("User", string(code))
}

func TestVote(t *testing.T) {
	v := &Vote{}
	v.ID = 2
	code, _ := json.Marshal(v)
	fmt.Println("User", string(code))
}
