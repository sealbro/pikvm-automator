package main

import (
	"fmt"
	"github.com/sealbro/pikvm-automator/internal/macro"
)

func main() {
	expressions := []string{
		//"@left|@right",
		"@left+@200'100+1s+@920'560+1s+@200'200",
		"@left+@200'100+1s+@920'560+1s+@200'200|@right+1s+MetaLeft+KeyD+@0'0+@920'560+@200'200",
		//"@0'0|@920'560|@left|@200'200|@right",
		//"1ms",
		//"2s",
		//"MetaLeft",
		//"MetaLeft|@12'21",
		//"MetaLeft+KeyD",
		//"MetaLeft+KeyD+1ms|KeyC+KeyV+10ms",
		//"MetaLeft+2ms+KeyD+1ms",
		//"MetaLeft+KeyS+1ms|100s",
		//"[10](MetaLeft+KeyD+1ms|10s)",
		//"[10](MetaLeft+KeyD+1ms|10s)|MetaLeft|5s",
		//"[5]([20](MetaLeft+KeyD+1ms|12s)|KeyA|15s)",
		//"[5]([20](MetaLeft+KeyD+1ms|12s)|KeyA|15s)|MetaLeft|5s",
		//"[10](MetaLeft+KeyD+1ms|10s)MetaLeft|5s",
		//"[5]([20](MetaLeft+KeyD+1ms|12s)|KeyA|15s)MetaLeft|5s",
	}

	for _, exp := range expressions {
		macros := macro.New(exp)
		//fmt.Println(macros)
		group := macros.Parse()
		compiledExp := macros.String()
		if compiledExp != exp {
			fmt.Println(exp)
			fmt.Println(compiledExp)
		}
		if group.TotalDelay() > 0 {
			fmt.Println(group.TotalDelay(), ": ", compiledExp)
		}
	}
}
