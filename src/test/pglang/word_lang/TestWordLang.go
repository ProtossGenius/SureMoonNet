package main

import (
	"basis/smn_analysis/lw_analysis"
	"basis/smn_data"
	"bufio"
	"fmt"
	"os"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	file, err := os.Open("./datas/code_tools/type_def/test.tdf")
	fin := bufio.NewScanner(file)
	fin.Split(bufio.ScanWords)
	check(err)
	sm := lw_analysis.GetStructStateMachine()
	for fin.Scan() {
		sm.Read(&lw_analysis.LangInput{Word: fin.Text()})
	}
	result := sm.GetResultChan()
	go func() {
		for {
			time.Sleep(1)
		}
	}()
	for {
		res := <-result
		if res == nil {
			continue
		}
		out := res.(*lw_analysis.ResultStruct)
		fmt.Println(smn_data.ValToJson(*out.Result))
	}
}
