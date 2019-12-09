package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {

	var buffer string
	var bufferArgs []string
	var bufferPipe []string
	for {
		fmt.Printf("User@%s:~ ", os.Getenv("USER"))
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		buffer = scanner.Text()
		flag := processString(buffer, &bufferArgs, &bufferPipe)

		if flag == 1 {
			go processArgsPipe(bufferArgs, bufferPipe)
		} else if flag == 2 {
			 go processArgs(bufferArgs)
		} else if flag == 3 {
			 go shellCommand(bufferArgs)
		}
		fmt.Printf("\n")
	}
}

func checkPipe(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] == '|' {
			return true
		}
	}
	return false
}
func parseSpace(str string) []string {

	var rows []string
	rows = strings.Split(str, " ")
	return rows
}
func parsePipe(str string) []string {
	var rows []string
	rows = strings.Split(str, "|")
	return rows
}
func processString(str string, args *[]string, argspiep *[]string) int {

	check := checkPipe(str)
	parsePipe := parsePipe(str)
	var result int
	if check {
		*args = parseSpace(parsePipe[0])
		*argspiep = parseSpace(parsePipe[1])
		result = 1
	} else {
		*args = parseSpace(str)
		result = 2
	}
	if checkShellCommand(*args) {
		result = 3
	}
	return result

}
func shellCommand(buffer []string) {

	switch buffer[0] {

	case "cd":
		err := os.Chdir(buffer[1])
		if err != nil {
			log.Fatal(err)
		}
		break
	case "mkdir":
		err := os.Mkdir(buffer[1], 007)
		if err != nil {
			log.Fatal(err)
		}
		break
	case "rename":
		err := os.Rename(buffer[1], buffer[2])
		if err != nil {
			log.Fatal(os.LinkError{})
		}
		break
	case "remove":
		err := os.Remove(buffer[1])
		if err != nil {
			log.Fatal(os.PathError{})
		}
		break
	case "getpid":
		fmt.Printf("%d", os.Getpid())
		break

	}

}

func checkShellCommand(buffer []string) bool {

	command := []string{"cd", "mkdir", "rename", "remove", "getpid"}
	for i := 0; i < 4; i++ {
		if buffer[0] == command[i] {
			return true
		}
	}
	return false
}

func processArgsPipe(buffer []string, bufferPipe []string) {

	cmd := exec.Command(buffer[0], bufferPipe[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err!=nil{
		fmt.Fprintln(os.Stderr,err)
	}
}

func processArgs(buffer []string) {
	cmd := exec.Command(buffer[0], buffer[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err!=nil{
		fmt.Fprintln(os.Stderr,err)
	}

}
