package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const USAGE = `befunge is an esoteric programming language.
Usage:
	befunge [flags] file.bf
Flags:`

func readFile(filename string) ([25][80]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return [25][80]byte{}, errors.New(fmt.Sprintf("Could not open file: '%s'", filename))
	}
	defer file.Close()

	var result [25][80]byte

	scanner := bufio.NewScanner(file)
	line := 0
	for scanner.Scan() {
		bytes := []byte(scanner.Text())
		for i, c := range bytes {
			result[line][i] = c
		}
		for i := len(bytes); i < 80; i++ {
			result[line][i] = ' '
		}
		line += 1
		if line >= 25 {
			break
		}
	}
	for j := line; j < 25; j++ {
		for i := 0; i < 80; i++ {
			result[j][i] = ' '
		}
	}
	return result, nil
}

func programString(program [25][80]byte) string {
	result := ""
	for _, line := range program {
		lineString := ""
		for _, c := range line {
			lineString += string(c)
		}
		lineString = strings.TrimRight(lineString, " ")
		if lineString != "" {
			result += lineString
			result += "\n"
		}
	}
	return result
}

func main() {
	var debug bool
	flag.BoolVar(
		&debug, "debug", false,
		"pause the program and show the PC and stack after every step.",
	)
	flag.Usage = func() {
		fmt.Println(USAGE)
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println(USAGE)
		flag.PrintDefaults()
		return
	}

	filename := flag.Args()[0]
	program, err := readFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	programs := programString(program)
	stack := make([]int, 0)

	pcx := 0
	pcy := 0
	dirx := 1
	diry := 0

	var n, c int

	skipping := false
	stringing := false

	if debug {
		fmt.Printf("program:\n%s\n", programs)
	}

loop:
	for {
		if debug {
			bufio.NewReader(os.Stdin).ReadBytes('\n') 
			fmt.Printf("PC: (%d, %d)\n", pcx, pcy)
			fmt.Printf("character: %q\n", program[pcy][pcx])
		}

		if skipping {
			skipping = false
			pcx, pcy = step(pcx, pcy, dirx, diry)
			continue
		}

		if stringing {
			if program[pcy][pcx] == '"' {
				stringing = false
			} else {
				stack = push(stack, int(program[pcy][pcx]))
			}
		} else {
			switch command := program[pcy][pcx]; command {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				stack = push(stack, int(command-'0'))
			case '+':
				stack = add(stack)
			case '-':
				stack = sub(stack)
			case '*':
				stack = mul(stack)
			case '/':
				stack = div(stack)
			case '%':
				stack = mod(stack)
			case '!':
				stack = not(stack)
			case '>':
				dirx = 1
				diry = 0
			case '<':
				dirx = -1
				diry = 0
			case 'v':
				dirx = 0
				diry = 1
			case '^':
				dirx = 0
				diry = -1
			case '?':
				dir := rand.Intn(4)
				switch dir {
				case 0:
					dirx = 1
					diry = 0
				case 1:
					dirx = -1
					diry = 0
				case 2:
					dirx = 0
					diry = 1
				default:
					dirx = 0
					diry = -1
				}
			case '_':
				diry = 0
				stack, dirx = condition(stack)
			case '|':
				dirx = 0
				stack, diry = condition(stack)
			case '&':
				stack, err = inputNum(stack)
				if err != nil {
					fmt.Println(err)
					return
				}
			case '~':
				stack, err = inputChar(stack)
				if err != nil {
					fmt.Println(err)
					return
				}
			case '.':
				stack, n = pop(stack)
				fmt.Printf("%d ",  n)
			case ',':
				stack, c = pop(stack)
				fmt.Printf("%c", c)
			case '#':
				skipping = true
			case ':':
				stack = duplicate(stack)
			case '$':
				stack, _ = pop(stack)
			case '\\':
				stack = swap(stack)
			case '`':
				stack = compare(stack)
			case 'g':
				stack = get(program, stack)
			case 'p':
				program, stack = put(program, stack)
			case '"':
				stringing = true
			case '@':
				break loop
			case ' ':
			default:
				fmt.Printf("Invalid character: %q\n", command)
				return
			}
		}

		pcx, pcy = step(pcx, pcy, dirx, diry)

		if debug {
			fmt.Printf("stack: %v\n", stack)
		}
	}
}

func step(pcx, pcy, dirx, diry int) (int, int) {
	pcx += dirx
	if pcx == 80 {
		pcx = 0
	} else if pcx == -1 {
		pcx = 79
	}
	pcy += diry
	if pcy == 25 {
		pcy = 0
	} else if pcy == -1 {
		pcy = 24
	}
	return pcx, pcy
}

func push(stack []int, n int) []int {
	return append(stack, n)
}

func pop(stack []int) ([]int, int) {
	result := peek(stack)
	if len(stack) == 0 {
		return make([]int, 0), result
	}
	return stack[:len(stack)-1], result
}

func peek(stack []int) int {
	if len(stack) == 0 {
		return 0
	}
	return stack[len(stack)-1]
}

func add(stack []int) []int {
	stack, b := pop(stack)
	stack, a := pop(stack)
	return push(stack, a+b)
}

func sub(stack []int) []int {
	stack, b := pop(stack)
	stack, a := pop(stack)
	return push(stack, a-b)
}

func mul(stack []int) []int {
	stack, b := pop(stack)
	stack, a := pop(stack)
	return push(stack, a*b)
}

func div(stack []int) []int {
	stack, b := pop(stack)
	stack, a := pop(stack)
	return push(stack, a/b)
}

func mod(stack []int) []int {
	stack, b := pop(stack)
	stack, a := pop(stack)
	return push(stack, a%b)
}

func not(stack []int) []int {
	stack, a := pop(stack)
	var res int
	if a == 0 {
		res = 1
	} else {
		res = 0
	}
	return push(stack, res)
}

func condition(stack []int) ([]int, int) {
	stack, a := pop(stack)
	if a == 0 {
		return stack, 1
	}
	return stack, -1
}

func inputNum(stack []int) ([]int, error) {
	fmt.Print("Enter a number: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return []int{}, errors.New("Error reading input")
	}
	line = line[:len(line)-1]
	n, err := strconv.ParseInt(line, 10, 0)
	if err != nil {
		return []int{}, errors.New("Invalid number")
	}
	return push(stack, int(n)), nil
}

func inputChar(stack []int) ([]int, error) {
	fmt.Print("Enter a character: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return []int{}, errors.New("Error reading input")
	}
	line = line[:len(line)-1]
	if len(line) == 0 {
		return []int{}, errors.New("No input")
	}
	if len(line) > 1 {
		fmt.Printf("Ignoring extra characters, input is %q\n", line[0])
	}
	return push(stack, int(line[0])), nil
}

func duplicate(stack []int) []int {
	stack, n := pop(stack)
	stack = push(stack, n)
	stack = push(stack, n)
	return stack
}

func swap(stack []int) []int {
	stack, a := pop(stack)
	stack, b := pop(stack)
	stack = push(stack, a)
	stack = push(stack, b)
	return stack
}

func compare(stack []int) []int {
	stack, b := pop(stack)
	stack, a := pop(stack)
	var n int
	if a > b {
		n = 1
	} else {
		n = 0
	}
	return push(stack, n)
}

func get(program [25][80]byte, stack []int) []int {
	stack, y := pop(stack)
	stack, x := pop(stack)
	v := program[y][x]
	return push(stack, int(v))
}

func put(program [25][80]byte, stack []int) ([25][80]byte, []int) {
	stack, y := pop(stack)
	stack, x := pop(stack)
	stack, v := pop(stack)
	program[y][x] = byte(v)
	return program, stack
}
