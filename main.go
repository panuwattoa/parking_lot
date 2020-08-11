package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"parkinglot/command"
	"regexp"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		inputCmd()
		return
	}
	filename := os.Args[1]
	fmt.Println(filename)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	line := strings.Split(string(file), "\n")
	for _, message := range line {
		cmd(message)
	}
}

func inputCmd() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	cmd(input)
	inputCmd()
}

func splitCmd(msg string) (cmd string, param string) {
	param = ""
	tempMsg := regexp.MustCompile(`\s+`).Split(msg, 2)
	switch len(tempMsg) {
	case 1:
		cmd = tempMsg[0]
	case 2:
		cmd = tempMsg[0]
		param = tempMsg[1]
	}
	return cmd, param
}

func cmd(msg string) {
	cmd, param := splitCmd(msg)
	commandFunc, ok := command.CommandMapper[cmd]
	if !ok {
		// Invalid command
		fmt.Println("Bad command '" + cmd + "'")
		inputCmd()
		return
	}
	commandFunc(cmd, param)
	fmt.Println()
}
