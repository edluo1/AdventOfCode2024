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

func getInput(f string) []int { // Value to count.
    file, err := os.Open(f)
    check(err)
    defer file.Close()

    scanner := bufio.NewScanner(file)

    // Make sure to copy each line so that it doesn't break
    var stones []int
    for scanner.Scan() {
        split := s.Split(scanner.Text(), " ")
        for _, val := range split {
            n, err := strconv.Atoi(val)
            check(err)
            stones = append(stones, n)
        }
    }

    return stones
}

func getDigitCount(n int) {

}

func blink(st map[int]int) map[int]int { // Value to count.
    nextMap := make(map[int]int)

    for k, v := range(st) {
        if k == 0 {
            nextMap[1] += v
        } else {
            str_form := strconv.Itoa(k)
            l_str := len(str_form)
            if l_str % 2 == 1 {
                nextMap[k*2024] += v
            } else {
                left, err := strconv.Atoi(str_form[:l_str/2])
                check(err)
                right, err := strconv.Atoi(str_form[l_str/2:])
                check(err)
                nextMap[left] += v
                nextMap[right] += v
            }
        }
    }

    return nextMap
}

// Input: Stones with engravings. On blink:
// 0s become 1s.
// Even digit counts are split in two.
// Odd digit counts multiply the value by 2024.
// Output: Total stone count after 25 blinks.
func part1(stones []int) int {
    // Add values to map.
    stoneMap := make(map[int]int)

    for _, val := range(stones) {
        stoneMap[val] += 1
    }
    
    for range 25 { // Each blink
        stoneMap = blink(stoneMap)
    }
    
    stoneSum := 0
    for _, val := range(stoneMap) {
        stoneSum += val
    }
    return stoneSum
}

// Input: Same stones, same rules.
// Output: Total stone count after 75 blinks.
func part2(stones []int) int {
    // Add values to map.
    stoneMap := make(map[int]int)

    for _, val := range(stones) {
        stoneMap[val] += 1
    }
    
    for range 75 { // Each blink
        stoneMap = blink(stoneMap)
    }
    
    stoneSum := 0
    for _, val := range(stoneMap) {
        stoneSum += val
    }
    return stoneSum
}

func main() {
    example := getInput("example.txt")
    stones := getInput("input.txt")

    fmt.Printf("example1: %v, should be 55312\n", part1(example))
    fmt.Printf("example2: %v, no example solution\n", part2(example))

    fmt.Printf("part1: %v\n", part1(stones))
    fmt.Printf("part2: %v\n", part2(stones))
}
