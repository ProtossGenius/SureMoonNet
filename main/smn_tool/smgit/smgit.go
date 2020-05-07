package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
)

/*
 *   command:
 *   pull push
 */
var (
	Comment string
)

const (
	PULL = "pull"
	PUSH = "push"
	SP   = "sp"
	RSP  = "rsp"

	CfgDir = "~/.smtools/smgit-cfg"
)

type FlagFunc func()

type FlagMap map[string]FlagFunc

func (this FlagMap) Register(key string, f FlagFunc) {
	this[key] = f
}

func (this FlagMap) Parse(key string) {
	if f, ok := this[key]; ok {
		f()
	}
}
func ec(name string, arg ...string) {
	c := exec.Command(name, arg...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	check(err)
}
func ffPull() {
	//git stash
	fmt.Println("doing stash")
	ec("git", "stash", "save", fmt.Sprintf(`'save when %s'`, time.Now().Format("2016-01-02 15:04:05")))
	//git pull
	fmt.Println("doing pull")
	ec("git", "pull")
	fmt.Println("doing pop --index")
	//git pop
	ec("git", "stash", "pop", "--index")
}

func ffPush() {
	if Comment == "" {
		panic(fmt.Errorf("no comment message"))
	}
	//make test
	fmt.Println("make test")
	ec("make", "test")
	//make clean
	fmt.Println("make clean")
	ec("make", "clean")
	//git add .
	fmt.Println("git add -A")
	ec("git", "add", "-A")
	//git commit -m
	fmt.Println("git commit -m ", fmt.Sprintf(`"%s"`, Comment))
	ec("git", "commit", "-m", fmt.Sprintf(`""%s"`, Comment))
	fmt.Println("git push")
	ec("git", "push")
}

func ffSp() {
	//create directory
	if !smn_file.IsFileExist(CfgDir) {
		err := os.MkdirAll(CfgDir, os.ModePerm)
		check(err)
	}
}

func ffRsp() {
}

var flagMap = FlagMap{
	PULL: ffPull,
	PUSH: ffPush,
	SP:   ffSp,
	RSP:  ffRsp,
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	flag.StringVar(&Comment, "m", "", "comment message for push")
	flag.Parse()
	doFlag := false
	args := flag.Args()
	fmt.Println(args)
	argMap := make(map[string]bool, len(args))
	if Comment != "" {
		argMap[PUSH] = true
	}
	for _, arg := range args {
		argMap[arg] = true
	}
	for arg := range argMap {
		flagMap.Parse(arg)
		doFlag = true
	}
	if !doFlag {
		fmt.Println(`smgit pull : pull from remote
------------equals--------------
		git stash save ''
		git pull 
		git stash pop
################################
		smgit -m [push] : push to remote, push not must
------------equals--------------
		make test
		make clean
		git add .
		git comment -m "..."
		git push
################################ wait for add
		sp  : startup pull
		rsp : remove statup pull		
$ git status 
`)
		ec("git", "status")
	}
	fmt.Println("################### smgit FINISH")
}
