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

func getReports(f string) [][]int {
    file, err := os.Open(f)
    check(err)
    defer file.Close()

    scanner:= bufio.NewScanner(file)

    var reports [][]int
    for scanner.Scan() {
        split := s.Split(scanner.Text(), " ")
        vals_int := make([]int, len(split))
        for idx, elem := range(split) {
            val, err := strconv.Atoi(elem)
            check(err)
            vals_int[idx] = val
        }
        reports = append(reports, vals_int)
    }

    return reports
}

func absDiffInt(x, y int) int {
    if x < y {
        return y - x
    }
    return x - y
}

// Input: Reports as a list of lists.
// Identify the differences and see if they are constant and the diff isn't too high.
func part1(reports [][]int) int {
    safe := 0
    for _, report := range(reports) {
        increasing := false
        isSafe := true

        for j, val := range(report) {
            if j == 1 {
                if val > report[j-1] {
                    increasing = true
                } else if val < report[j-1] {
                    increasing = false
                } else {
                    isSafe = false
                    break
                }
            } 
            if j >= 1 {
                isSafe = checkSafeWithLastValue(report, j, 1, increasing)
                if !isSafe {
                    break
                }
            }
        }

        if isSafe {
            safe += 1
        }
    }
    return safe
}

func checkSafeWithLastValue(report []int, idx int, gap int, increasing bool) bool {
    if idx < gap { // Interpreting comparisons with non-existing values as true.
        return true
    }
    if idx >= len(report) {
        return true
    }
    // // fmt.Printf("compared %v with %v while increasing: %v\n", report[idx-gap], report[idx], increasing)
    diff := absDiffInt(report[idx], report[idx-gap])
    if diff < 1 || diff > 3 {
        // fmt.Printf("diff too big or equal, ")
        return false
    }
    if report[idx] > report[idx-gap] && !increasing {
        // fmt.Printf("increasing when shouldn't, ")
        return false
    } else if report[idx] < report[idx-gap] && increasing {
        // fmt.Printf("decreasing when shouldn't, ")
        return false
    }

    return true
}

// Input: Reports as a list of lists.
// Check to see if the bad level can be mitigated.
// It's possible that it could be 1 2 5 3 4 5 6 where the issue happens in the last safe value
// Or that it could be 1 2 3 20 4 5 6 where the issue happens on the bad number
// So we have to test skipping both index i and i-1 where i is the one that caused the issue
// And we need to make sure we can handle both directions like with 10 20 9 8 7 6 5 4
func part2(reports [][]int) int {
    safe := 0

    for _, report := range(reports) {
        // Case 1: skip the first value
        increasing_skip1 := report[2] - report[1] > 0
        isSafe := true
        for i := 2; i < len(report) && isSafe; i++ {
            safe := checkSafeWithLastValue(report, i, 1, increasing_skip1)
            if !safe {
                isSafe = false
            }
        }
        if isSafe {
            safe += 1
            // fmt.Printf("%v: skipped 0, safe\n", report)
            continue
        } else {
            // fmt.Printf("%v: didn't skip 0, unsafe\n", report)
        }

        // Case 2: skip the second value
        increasing_skip2 := report[2] - report[0] > 0
        // Make sure we're safe with the jump
        isSafe = checkSafeWithLastValue(report, 2, 2, increasing_skip2)
        for i := 3; i < len(report) && isSafe; i++ {
            safe := checkSafeWithLastValue(report, i, 1, increasing_skip2)
            if !safe {
                isSafe = false
            }
        }
        if isSafe {
            safe += 1
            // fmt.Printf("%v: skipped 1, safe\n", report)
            continue
        } else {
            // fmt.Printf("%v: didn't skip 1, unsafe\n", report)
        }

        // Case 3: skip later value
        // Need to check what happens if we skip the offender or the one before the offender.
        increasing := report[1] - report[0] > 0
        skipIdx := -1
        isSafe = checkSafeWithLastValue(report, 1, 1, increasing) // Doing this check now 
        for i := 2; i < len(report) && isSafe; i++ {
            if !skipUsed(skipIdx) && i == len(report)-1 {
                // fmt.Printf("%v: skipped %v, ", report, i)
                break
            }
            safe := checkSafeWithLastValue(report, i, 1, increasing)
            if !safe {
                if skipUsed(skipIdx) { 
                    // Already skipped a value, return false.
                    isSafe = false
                    break
                }

                // Otherwise, check if we need to do 1 X 3 4 or 1 2 X 4. Consider four values: the one ahead, this one, and the two behind.
                if i != 2 && checkSafeWithLastValue(report, i, 2, increasing) && checkSafeWithLastValue(report, i+1, 1, increasing) {
                    // Skipping idx 2 since we already handled skipping 1 last case
                    skipIdx = i-1
                } else if checkSafeWithLastValue(report, i+1, 2, increasing) && checkSafeWithLastValue(report, i-1, 1, increasing) {
                    skipIdx = i
                } else {
                    // Neither skip works, return false
                    isSafe = false
                    break
                }
                i += 1
            }
        }
        if isSafe {
            // fmt.Printf("safe\n")
            safe += 1
        }
    }
    return safe
}

func skipUsed(skipIdx int) bool {
    return skipIdx != -1
}

func main() {
    example := getReports("example.txt")
    reports := getReports("input.txt")

    fmt.Printf("example1: %v, should be 2\n", part1(example))
    fmt.Printf("example2: %v, should be 6 (added extra examples)\n", part2(example))

    fmt.Printf("part1: %v\n", part1(reports))
    fmt.Printf("part2: %v\n", part2(reports))
}

