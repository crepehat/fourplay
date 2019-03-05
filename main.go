package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/crepehat/fourplay/fourplay"
)

func RunTests(path string) {
	maxTurns := fourplay.BoardWidth * fourplay.BoardHeight / 2

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		file, err := os.Open(filepath.Join(path, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// fmt.Println(scanner.Text())
			lineSplit := strings.Split(scanner.Text(), " ")
			board, err := fourplay.CreateFromSequence(lineSplit[0])
			// board.Print()
			if err != nil {
				fmt.Println(err)
			}

			result := board.NegaMax(-maxTurns, maxTurns)

			fmt.Println(lineSplit[0], lineSplit[1], result)
			givenAnswer, err := strconv.Atoi(lineSplit[1])
			if int8(givenAnswer) != result {
				panic(lineSplit[0])
			}
		}
	}
}

// 14363756335665245414
func main() {
	var runTests = flag.Bool("tests", false, "should i run the tests")
	var path = flag.String("path", "./tests", "where should i look for test files")
	flag.Parse()

	if *runTests {
		fmt.Println("swag")
		RunTests(*path)
	}

	board, err := fourplay.CreateFromSequence("441432421426")
	board.Print()
	if err != nil {
		fmt.Println(err)
	}

	maxTurns := fourplay.BoardWidth * fourplay.BoardHeight / 2
	result := board.NegaMax(-maxTurns, maxTurns)
	fmt.Println(result)

}
