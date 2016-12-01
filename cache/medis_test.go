package cache

import (
	"fmt"
	"io/ioutil"
	"testing"
)

type Best struct {
	Name string
}

func TestMedis(t *testing.T) {
	ioutil.WriteFile("test.txt", []byte(""), 0666)
	b := NewMedis("test.txt")
	//b.Clean()
	//*/
	fmt.Println("test.o")
	str := `func(unc(unc(unc(unc(unc(unc(unc(unc(unc(unc(unc(unc(unc(unc(unc(unc(unc(unc(unc(unc(unc(unc(unc(){
	
	}aa`
	b.Set("wang", "wang")
	b.Set("hait", &Best{str}, 90)
	b.Set("tao", "tao")
	b.Set("haha", "haha")
	b.Del("haha")

	//
	///*
	fmt.Println("读取")
	c := NewMedis("test.txt")
	d := new(Best)
	c.GetBy("hait", &d)
	fmt.Println(d)
}
