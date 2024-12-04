package main

import ( 
    "fmt" 
    "bufio" 
    "os" 
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func getwordsearch(f string) [][]byte {
    file, err := os.Open(f)
    check(err)
    defer file.Close()

    scanner:= bufio.NewScanner(file)

    // Need to increase the capacity since default somehow breaks it
    const maxCapacity = 512*1024  
    buf := make([]byte, maxCapacity)
    scanner.Buffer(buf, maxCapacity)

    var wordsearch [][]byte
    for scanner.Scan() {
        wordsearch = append(wordsearch, scanner.Bytes())
    }

    return wordsearch
}

const XMAS = "XMAS"
const SAMX = "SAMX"

func wordExists(search [][]byte, x int, y int, word string, dx int, dy int) bool {
    for i, _ := range(word) {
        if search[y+dy*i][x+dx*i] != word[i] {
            return false
        }
    }
    return true
}

// Returns how many times the word is found in the search.
// Only search right, down-right, down, down-left; we can find the reverse word the other way
func find_in_search(search [][]byte, x int, y int, word string) int {
    found := 0
    len_y := len(search)
    len_x := len(search[y])

    reach := len(word)-1

    // Horizontal search
    if x + reach < len_x {
        if wordExists(search, x, y, word, 1, 0) {
            found += 1
        }
    }

    // Diagonal right search
    if x + reach < len_x && y + reach < len_y {
        if wordExists(search, x, y, word, 1, 1) {
            found += 1
        }
    }

    // Down search
    if y + reach < len_y {
        if wordExists(search, x, y, word, 0, 1) {
            found += 1
        }
    }

    // Down-left search
    if y + reach < len_y && x >= reach {
        if wordExists(search, x, y, word, -1, 1) {
            found += 1
        }
    }

    return found
}

// Input: Word search
// Find "XMAS" in each direction (left, right, up, down, diagonally)
func part1(wordsearch [][]byte) int {
    xmases := 0
    for y, line := range(wordsearch) {
        for x, c := range(line) {
            if c == 'X' {
                xmases += find_in_search(wordsearch, x, y, XMAS)
            } else if c == 'S' { // Reversed text
                xmases += find_in_search(wordsearch, x, y, SAMX)
            }
        }
    }
    return xmases
}

const MAS = "MAS"
const SAM = "SAM"

func crossExists_right(search [][]byte, x int, y int, word string) bool {
    reach := len(word)-1
    for i, _ := range(word) {
        if search[y+i][x+i] != word[i] {
            return false
        }
        if search[y+i][x+reach-i] != word[i] {
            return false
        }
    }
    return true
}

func crossExists_down(search [][]byte, x int, y int, word string) bool {
    reach := len(word)-1
    for i, _ := range(word) {
        if search[y+i][x+i] != word[i] {
            return false
        }
        if search[y+reach-i][x+i] != word[i] {
            return false
        }
    }
    return true
}

// Returns how many times the word is found in the search.
// Only search right, down-right, down, down-left; we can find the reverse word the other way
func find_cross_search(search [][]byte, x int, y int, word string) int {
    found := 0
    len_y := len(search)
    len_x := len(search[y])

    reach := len(word)-1

    // Finding squares
    if x + reach < len_x && y + reach < len_y{
        if crossExists_right(search, x, y, word) {
            found += 1
        }
        if crossExists_down(search, x, y, word) {
            found += 1
        }
    }

    return found
}

// Input: Word search
// Find "MAS" crossing each other (left, right, up, down, diagonally)
func part2(wordsearch [][]byte) int {
    xmases := 0
    for y, line := range(wordsearch) {
        for x, c := range(line) {
            if c == 'M' {
                xmases += find_cross_search(wordsearch, x, y, MAS)
            } else if c == 'S' { // Reversed text
                xmases += find_cross_search(wordsearch, x, y, SAM)
            }
        }
    }
    return xmases
}

func main() {
    example := getwordsearch("example.txt")
    wordsearch := getwordsearch("input.txt")

    fmt.Printf("example1: %v, should be 18\n", part1(example))
    fmt.Printf("example2: %v, should be 9\n", part2(example))

    fmt.Printf("part1: %v\n", part1(wordsearch))
    fmt.Printf("part2: %v\n", part2(wordsearch))
}
