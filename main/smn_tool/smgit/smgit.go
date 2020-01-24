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

	CfgPath = "~/.smtools/smgit-cfg"
)

type FlagFunc func() error

type FlagMap map[string]FlagFunc

func (this FlagMap) Register(key string, f FlagFunc) {
	this[key] = f
}

func (this FlagMap) Parse(key string) error {
	if f, ok := this[key]; ok {
		return f()
	}
	return nil
}

func ffPull() error {
	//git stash
	fmt.Println("doing stash")
	c := exec.Command("git", "stash", "save", fmt.Sprintf(`'save when %s'`, time.Now().Format("2016-01-02 15:04:05")))
	if err := c.Run(); err != nil {
		return err
	}
	//git pull
	fmt.Println("doing pull")
	c = exec.Command("git", "pull")
	if err := c.Run(); err != nil {
		return err
	}
	fmt.Println("doing pop --index")
	//git pop
	c = exec.Command("git", "stash", "pop", "--index")
	return c.Run()
}

func ffPush() error {
	if Comment == "" {
		panic(fmt.Errorf("no comment message"))
	}
	//git add .
	fmt.Println("git add .")
	c := exec.Command("git", "add", ".")
	c.Stdout = os.Stdout
	if err := c.Run(); err != nil {
		return err
	}
	//git commit -m
	c = exec.Command("git", "commit", "-m", fmt.Sprintf(`""%s"`, Comment))
	if err := c.Run(); err != nil {
		return err
	}
	c = exec.Command("git", "push")
	return c.Run()
}

func ffSp() error {
	//create directory
	if !smn_file.IsFileExist(CfgPath) {
		if err := os.MkdirAll(CfgPath, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func ffRsp() error {
	return nil
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
	fmt.Println(Comment)
	doFlag := false
	args := flag.Args()
	for _, arg := range args {
		check(flagMap.Parse(arg))
		doFlag = true
	}
	if !doFlag {
		fmt.Println(`pull : pull from remote
------------equals--------------
		git stash save ''
		git pull 
		git stash pop
################################
		push : push to remote
------------equals--------------
		git add .
		git comment -m "..."
		git push
################################
		sp  : startup pull
		rsp : remove statup pull		
`)
	}
}
