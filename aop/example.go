package aop

import (
	"fmt"
	"reflect"
)

func init()  {
	RegisterPoint(reflect.TypeOf((*Account)(nil)))
	RegisterAspect(&Aspect{})
}

type Aspect struct {}

func (a *Aspect) Before(point *JoinPoint) bool {
	fmt.Println("before",point.Method.Name)
	return true
}

func (a *Aspect) After(point *JoinPoint) {
	fmt.Println("after",point.Method.Name)
}

func (a *Aspect) Finally(point *JoinPoint) {
	fmt.Println("finally",point.Method.Name)
}


//匹配到的方法名
func (a *Aspect) GetAspectExpress() string {
	return ".*\\."
}

type Account struct {
	Id int
}

func (h *Account) HelloAccount() {
	fmt.Println("HelloAccount")
}

func example()  {
	h := &Account{}
	h.HelloAccount()
}