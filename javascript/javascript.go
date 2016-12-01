package main // javascript

import (
	"fmt"

	"github.com/robertkrimen/otto"
)

func main() {
	vm := otto.New()
	vm.Run(`
	    abc = 2 + 2;
	    console.log("The value of abc is " + abc); // 4
	`)
	if value, err := vm.Get("abc"); err == nil {
		value.ToInteger()
	}
	fmt.Println()
	vm.Set("def", map[string]interface{}{"name": "王海涛", "age": 22})
	vm.Run(`  
		 console.log(abc)
		 
	`)
}
