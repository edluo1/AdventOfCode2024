package main

import ( 
    "fmt" 
    "bufio" 
    "os"
    "strconv"
    s "strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func check_overflow(e error) bool {
    if e != nil {
        return true
    }
    return false
}

type equation struct {
    total int
    values []int
}

func toEquation(line string) equation {
    sides := s.Split(line, ": ")
    testVal, err := strconv.Atoi(sides[0])
    check(err)
    split := s.Split(sides[1], " ")
    vals_int := make([]int, len(split))
    for idx, elem := range(split) {
        val, err := strconv.Atoi(elem)
        check(err)
        vals_int[idx] = val
    }
    return equation{ testVal, vals_int }
}

func getEquations(f string) []equation {
    file, err := os.Open(f)
    check(err)
    defer file.Close()

    scanner := bufio.NewScanner(file)

    // Make sure to copy each line so that it doesn't break
    var eq []equation
    for scanner.Scan() {
        line := scanner.Text()
        eq = append(eq, toEquation(line))
    }

    return eq
}

func concatenate(left int, right int) int {
    combine := fmt.Sprintf("%d", left) + fmt.Sprintf("%d", right)
    val, err := strconv.Atoi(combine)
    if check_overflow(err) {
        return -1
    }
    return val
}

func solve_eq(eq equation, val int, idx int, ver2 bool) bool {
    if val > eq.total {
        return false
    } else if idx == len(eq.values) {
        return val == eq.total
    }
    add := solve_eq(eq, val + eq.values[idx], idx+1, ver2)
    multiply := solve_eq(eq, val * eq.values[idx], idx+1, ver2)
    var conc bool
    if ver2 {
        concatenate_value := concatenate(val, eq.values[idx])
        if concatenate_value == -1 {
            conc = false
        } else {
            conc = solve_eq(eq, concatenate_value, idx+1, ver2)
        }
    } else {
        conc = false
    }

    return add || multiply || conc
}

func solvable(eq equation, part2 bool) bool {
    // If you try to bruteforce it it'll take O(2^n) time.
    // Can reduce time by cutting any values that go over.
    return solve_eq(eq, 0, 0, part2)
}

// Input: list of equations
// Figure out which ones can be solved using + and * only
// explicitly from left to right, no PEMDAS
func part1(equations []equation) int {
    test_sum := 0

    for _, eq := range(equations) {
        if solvable(eq, false) {
            test_sum += eq.total
        }
    }

    return test_sum
}

// Input: list of equations
// Figure out which ones can be solved using +, *, and concatenate only
// explicitly from left to right, no PEMDAS
func part2(equations []equation) int {
    test_sum := 0

    for _, eq := range(equations) {
        if solvable(eq, true) {
            test_sum += eq.total
        }
    }

    return test_sum
}

func main() {
    example := getEquations("example.txt")
    equations := getEquations("input.txt")

    fmt.Printf("example1: %v, should be 3749\n", part1(example))
    fmt.Printf("example2: %v, should be 11387\n", part2(example))

    fmt.Printf("part1: %v\n", part1(equations))
    fmt.Printf("part2: %v\n", part2(equations))
}
