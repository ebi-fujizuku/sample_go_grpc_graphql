package common

import (
	"fmt"
	"runtime"
)

func getFunctionName()string{
	programCounter,_,_,_ := runtime.Caller(2)
	fn := runtime.FuncForPC(programCounter)
	return fn.Name()
}
// func PrintStart(msg string,num int){
func PrintStart(msg string){
	// PrintDelimiter(num)
	if msg == ""{
		fmt.Println("Start",":", getFunctionName())
	}else{
		fmt.Println("Start",msg,":", getFunctionName())
	}
}
// func PrintEnd(msg string,num int){
func PrintEnd(msg string){
	if msg == ""{
		fmt.Println("End",":",getFunctionName())
	}else{
		fmt.Println("End",msg,":",getFunctionName())
	}
	// PrintDelimiter(num)
}
const (
	DELIMITER_1 = "###############################################"
	DELIMITER_2 = "=============================="
	DELIMITER_3 = "----------------"
	DELIMITER_4 = "........"
)
func PrintDelimiter(num int){
	delimiter := ""
	switch num{
	case 1:
		delimiter = DELIMITER_1
	case 2:
		delimiter = DELIMITER_2
	case 3:
		delimiter = DELIMITER_3
	case 4:
		delimiter = DELIMITER_4
	}
	fmt.Println(delimiter)
}
