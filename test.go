package main

import "fmt"

func main(){

	a := []int{2,5,7,9,46}
	b:= []int{1,3,4,8,23,25}
	c := make([]int,11)


	i, j, k := 0,0, 0
	lena := len(a)
	lenb := len(b)

	for k = 0; k < (lena + lenb) ; k++ {
			//fmt.Println(k)
		if i < (lenb - 1) {
			if b[i] < a[j]{
				c[k] = b[i]
				i = i + 1
				fmt.Println("i",i)
			}else{
				c[k] = a[j]
				j = j + 1
				fmt.Println("j", j)
			}
		}
		if j < lena -1 {
			c[k] = a[j]
			j = j + 1
		}

		}
	

	fmt.Println(c)




}