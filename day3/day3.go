package main

import ( 
    "fmt" 
    "bufio" 
    "os" 
    "strconv"
    "regexp"
    s "strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func getMemory(f string) string {
    file, err := os.Open(f)
    check(err)
    defer file.Close()

    var sb s.Builder

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        sb.WriteString(scanner.Text())
    }

    return sb.String()
}

// Input: Reports as a list of lists.
// Output: Sum of all multiplication instructions.
func part1(memory string) int64 {
    sum := int64(0)
    re := regexp.MustCompile(`mul\(([0-9]+),([0-9]+)\)`)
    instructions := re.FindAllStringSubmatch(memory, -1)

    for _, instr := range(instructions) {
        val1, err := strconv.ParseInt(instr[1], 10, 64)
        check(err)
        val2, err := strconv.ParseInt(instr[2], 10, 64)
        check(err)
        product := val1 * val2
        sum += product
    }

    return sum
}

// Input: Reports as a list of lists.
// Output: Sum of all multiplication instructions after a do(), without instructions after a don't()
func part2(memory string) int64 {
    sum := int64(0)
    re := regexp.MustCompile(`mul\(([0-9]+),([0-9]+)\)`)

    do_strings := s.Split(memory, "do()")

    for _, mem_block := range(do_strings) {
        dont_sec := s.Index(mem_block, "don't()")
        mem := ""
        if dont_sec == -1 {
            mem = mem_block
        } else {
            mem = mem_block[:dont_sec]
        }
        instructions := re.FindAllStringSubmatch(mem, -1)

        for _, instr := range(instructions) {
            val1, err := strconv.ParseInt(instr[1], 10, 64)
            check(err)
            val2, err := strconv.ParseInt(instr[2], 10, 64)
            check(err)
            product := val1 * val2
            sum += product
        }
    }

    return sum
}

func main() {
    example := getMemory("example.txt")
    example2 := getMemory("example2.txt")
    memory := getMemory("input.txt")

    fmt.Printf("example1: %v, should be 161\n", part1(example))
    fmt.Printf("example2: %v, should be 48\n", part2(example2))

    fmt.Printf("part1: %v\n", part1(memory))
    fmt.Printf("part2: %v\n", part2(memory))
}

