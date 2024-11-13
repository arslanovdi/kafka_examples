package model

import "fmt"

type Human struct {
	Id      int32   `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	Age     int32   `json:"age,omitempty"`
	Student bool    `json:"student,omitempty"`
	City    *string `json:"city,omitempty"`
}

func (h Human) String() string {
	str := fmt.Sprintf("id:%d name:%s age:%d student:%v", h.Id, h.Name, h.Age, h.Student)
	if h.City != nil {
		str += fmt.Sprintf(" city: %s", *h.City)
	}
	return str
}
