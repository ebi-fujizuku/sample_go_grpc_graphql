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
func PrintStart(msg string){
	if msg == ""{
		fmt.Println("Start",":", getFunctionName())
	}else{
		fmt.Println("Start",msg,":", getFunctionName())
	}
}
func PrintEnd(msg string){
	if msg == ""{
		fmt.Println("End",":",getFunctionName())
	}else{
		fmt.Println("End",msg,":",getFunctionName())
	}
}
