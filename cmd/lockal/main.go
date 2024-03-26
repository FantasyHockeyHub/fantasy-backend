package main

import "fmt"

func main() {
	var arr1 = make([]int, 0, 8)
	var arr2 = make([]int, 0, 8)
	fmt.Scan(&arr1)
	fmt.Scan(&arr2)
	fmt.Println(arr1)
	fmt.Println(arr2)

}
