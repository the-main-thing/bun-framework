package main

import (
	"fmt"
	// "os"
	// "strings"
)

func main() {
	// f, err := os.ReadFile("./tsconfig.json")
	// if err != nil {
	// 	fmt.Println("Error opening file")
	// 	os.Exit(1)
	// }
	//
	// removeComments(&f)
	//
	// fmt.Println(string(f))
	//
	// err = os.WriteFile("./stripped.json", f, 0644)
	// if err != nil {
	// 	fmt.Println("Error writing file")
	// 	os.Exit(1)
	// }
	// word := []rune("as_")
	// word[0] = 't'
	// sliced := word[0:0+3]
	// fmt.Println(string(sliced))

	t := "type_"
	r := []rune(t)

	fmt.Println(string(r[0:len("type")]))
	fmt.Println(string(r[len("type")]))

	for i := 1; i < 1; i++ {
		fmt.Println("Nooooo")
	}
	for i := 1; i < 2; i++ {
		fmt.Println("Yay!")
	}

}

//
// type ContentType int
//
// const (
// 	CODE ContentType = iota
// 	SINGLE_LINE_COMMENT
// 	MULTILINE_COMMENT
// 	STRING
// )
//
// func removeComments(tsconfig *[]byte) {
// 	var state ContentType
// 	state = CODE
// 	var stringOpeningSymbol byte
// 	for i, b := range *tsconfig {
// 		switch state {
// 		case STRING:
// 			if b == stringOpeningSymbol {
// 				// i is 100% bigger than 0, because we can't enter STRING state before 0 index
// 				if (*tsconfig)[i-1] == '\\' {
// 					break
// 				}
// 				state = CODE
// 			}
// 			continue
// 		case MULTILINE_COMMENT:
// 			(*tsconfig)[i] = ' '
// 			if b == '*' {
// 				nextIndex := i + 1
// 				if nextIndex < len(*tsconfig) {
// 					nextChar := (*tsconfig)[nextIndex]
// 					if nextChar == '/' {
// 						(*tsconfig)[nextIndex] = ' '
// 						state = CODE
// 					}
// 				}
// 			}
// 			break
//
// 		case SINGLE_LINE_COMMENT:
// 			(*tsconfig)[i] = ' '
// 			if b == '\n' {
// 				state = CODE
// 			}
// 			break
//
// 		case CODE:
// 			if b == '/' {
// 				nextIndex := i + 1
// 				if nextIndex >= len(*tsconfig) {
// 					(*tsconfig)[i] = ' '
// 					state = SINGLE_LINE_COMMENT
// 					break
// 				}
// 				nextChar := (*tsconfig)[nextIndex]
// 				if nextChar == '/' {
// 					(*tsconfig)[i] = ' '
// 					state = SINGLE_LINE_COMMENT
// 					break
// 				}
// 				if nextChar == '*' {
// 					(*tsconfig)[i] = ' '
// 					state = MULTILINE_COMMENT
// 					break
// 				}
// 				break
// 			}
// 			if b == '\'' || b == '"' || b == '`' {
// 				stringOpeningSymbol = b
// 				state = STRING
// 			}
// 			break
// 		default:
// 			panic("Unknown state")
// 		}
//
// 	}
// }
