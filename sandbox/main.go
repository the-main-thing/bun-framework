package main

import "fmt"


func main() {
	str := "Hello"
	b := []byte(str)
	modifyString(&b)
	for index, char := range b {
		if char == 'o' {
			prev := prevRune(&b, index)
			prevAsString := string(prev)
			fmt.Println(prevAsString)
		}
	}

	fmt.Println(str)
	fmt.Println("pointer", string(b))
}

func modifyString(input *[]byte) {
	(*input)[0] = 'G'
}

func prevRune(input *[]byte, index int) rune {
	if index == 0 {
		return 0
	}

	return rune((*input)[index - 1])
}
