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

type order struct {
    left int
    right int
}

func getPageSetup(f string) ([]order, [][]int) {
    file, err := os.Open(f)
    check(err)
    defer file.Close()

    scanner := bufio.NewScanner(file)

    // Create the input.
    var setup []order
    for scanner.Scan() {
        line := scanner.Text()
        if (line == "") {
            break
        }

        parts := s.Split(line, "|")
        p1, err := strconv.Atoi(parts[0])
        check(err)
        p2, err := strconv.Atoi(parts[1])
        check(err)
        ord := order{ p1, p2 }
        setup = append(setup, ord)
    }

    // Create the page setup
    var updates [][]int
    for scanner.Scan() {
        line := scanner.Text()
        if (scanner.Text() == "") {
            break
        }
        pages_str := s.Split(line, ",")
        pages_int := make([]int, len(pages_str))
        for i := range pages_str {
            val, err := strconv.Atoi(pages_str[i])
            check(err)
            pages_int[i] = val
        }
        updates = append(updates, pages_int)
    }

    return setup, updates
}

// Input: Page ordering, pages produced
// The updates should be reduced to those where the setup is always followed.
// Return the middle page number and add them together.
// The input should have 49 choose 2 = 1176 unique ones and 49 numbers.
func part1(setup []order, updates [][]int) int {
    middle_sum := 0

    // Setup: for each page list all pages that must be to the left and not to the right.
    pages_to_left := make(map[int]map[int]bool)
    for _, ord := range(setup) {
        if len(pages_to_left[ord.right]) == 0 {
            pages_to_left[ord.right] = make(map[int]bool)
        }
        pages_to_left[ord.right][ord.left] = true
    }

    for _, update := range(updates) {
        // Setup a map of illegal pages from the current page.
        middle_page := update[len(update)/2]
        illegal_pages := make(map[int]bool)
        legal := true

        for _, page := range(update) {
            if illegal_pages[page] == true {
                // Illegal page, break
                legal = false
                break
            }

            // Add new pages in.
            left_pages := pages_to_left[page]
            for k, v := range(left_pages) {
                illegal_pages[k] = v
            }
        }

        if legal {
            middle_sum += middle_page
        }
    }

    return middle_sum
}

// Make a sequence of the values in update.
// NOTE: 12|23, 23|34, 34|12 is possible in the input, meaning the graph is acyclic.
// But that being said we can just use Go's built-in SortFunc to do this. Not fast but it works.
func sort_using_rules(gt_map map[int]map[int]bool, update []int) []int {
    n := make([]int, len(update))
    copy(n, update)
    slices.SortFunc(n, func(a, b int) int {
        if a == b {
            return 0
        } else if gt_map[b][a] == true {
            return 1
        } else {
            return -1
        }
    })
    return n
}

// Input: Page ordering, pages produced
// The page updates that are in the wrong order should be put back.
// Return the middle page number of ONLY those in the wrong order put back into the correct order.
func part2(setup []order, updates [][]int) int {
    middle_sum := 0

    // Setup: for each page list all pages that must be to the left and not to the right.
    pages_to_left := make(map[int]map[int]bool)
    for _, ord := range(setup) {
        if len(pages_to_left[ord.right]) == 0 {
            pages_to_left[ord.right] = make(map[int]bool)
        }
        pages_to_left[ord.right][ord.left] = true
    }

    for _, update := range(updates) {
        // Setup a map of illegal pages from the current page.
        illegal_pages := make(map[int]bool)
        legal := true
        correct_order := make([]int, len(update))

        for _, page := range(update) {
            if illegal_pages[page] == true {
                // Illegal page, sort them.
                legal = false

                correct_order = sort_using_rules(pages_to_left, update)
                break
            }

            // Add new pages in.
            left_pages := pages_to_left[page]
            for k, v := range(left_pages) {
                illegal_pages[k] = v
            }
        }
        // Only add if the original order was WRONG.
        if !legal {
            middle_sum += correct_order[len(correct_order)/2]
        }
    }

    return middle_sum
}

func main() {
    ex_rules, ex_updates := getPageSetup("example.txt")
    rules, updates := getPageSetup("input.txt")

    fmt.Printf("example1: %v, should be 143\n", part1(ex_rules, ex_updates))
    fmt.Printf("example2: %v, should be 123\n", part2(ex_rules, ex_updates))

    fmt.Printf("part1: %v\n", part1(rules, updates))
    fmt.Printf("part2: %v\n", part2(rules, updates))
}
