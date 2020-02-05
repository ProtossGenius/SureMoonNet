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
	fmt.Println("git commit -m ", Comment)
	c = exec.Command("git", "commit", "-m", fmt.Sprintf(`""%s"`, Comment))
	if err := c.Run(); err != nil {
		return err
	}
	fmt.Println("git push")
	c = exec.Command("git", "push")
	return c.Run()
}

func ffSp() error {
	//create directory
	if !smn_file.IsFileExist(CfgDir) {
		if err := os.MkdirAll(CfgDir, os.ModePerm); err != nil {
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
		check(flagMap.Parse(arg))
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
		git add .
		git comment -m "..."
		git push
################################ wait for add
		sp  : startup pull
		rsp : remove statup pull		

`)
	}
}
