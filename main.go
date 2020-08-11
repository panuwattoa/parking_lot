package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"parkinglot/command"
	"regexp"
)

func main() {

	if len(os.Args) < 2 {
		inputCmd()
		return
	}
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		cmd(line)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
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
