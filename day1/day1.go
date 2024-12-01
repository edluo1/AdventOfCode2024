package main

import ( 
    "fmt" 
    "bufio" 
    "os" 
    "slices" 
    "strconv"
    s "strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func getLists(f string) ([]int, []int) {
    file, err := os.Open(f)
    check(err)
    defer file.Close()

    scanner:= bufio.NewScanner(file)

    var f1 []int
    var f2 []int
    for scanner.Scan() {
        split := s.Split(scanner.Text(), "   ")
        n1, err := strconv.Atoi(split[0])
        check(err)
        n2, err := strconv.Atoi(split[1])
        check(err)

        f1 = append(f1, n1)
        f2 = append(f2, n2)
    }

    return f1, f2
}

func absDiffInt(x, y int) int {
    if x < y {
        return y - x
    }
    return x - y
}

// Input: Two lists of numbers side by side. 
// Need to sort both of them and then sum the differences (absolute values).
func part1(n1 []int, n2 []int) int {
    slices.Sort(n1)
    slices.Sort(n2)
    sum := 0
    for idx, _ := range n1 {
        v := absDiffInt(n1[idx], n2[idx])
        sum += v
    }

    return sum
}

// Input: Two lists of numbers side by side. 
// Need to multiply value on left by number of time value shows up on right.
func part2(n1 []int, n2 []int) int {
    freq := make(map[int]int)

    for _, val := range n2 {
        freq[val] += 1
    }

    score := 0
    for _, val := range n1 {
        score += val * freq[val]
    }

    return score
}

func main() {
    l1, l2 := getLists("input.txt")

    fmt.Printf("part1: %v\n", part1(l1, l2))
    fmt.Printf("part2: %v\n", part2(l1, l2))
}

