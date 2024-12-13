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

// Storing as int instead of byte since we need to identify areas.
func getInput(f string) [][]uint64 {
    file, err := os.Open(f)
    check(err)
    defer file.Close()

    scanner := bufio.NewScanner(file)

    // Make sure to copy each line so that it doesn't break
    var garden [][]uint64
    for scanner.Scan() {
        line := make([]byte, len(scanner.Bytes()))
        copy(line, scanner.Bytes())
        int_line := make([]uint64, len(scanner.Bytes()))

        for i, val := range(scanner.Bytes()) {
            int_line[i] = uint64(val)
        }
        garden = append(garden, int_line)
    }

    return garden
}

type dir struct {
    x int
    y int
}

func outsideMap(x int, y int, x_size int, y_size int) bool {
    return x < 0 || y < 0 || x >= x_size || y >= y_size
}

func localPerimeter(garden [][]uint64, x, y, x_size, y_size int) int {
    directions := make([]dir, 4)
    directions[0] = dir{1, 0}
    directions[1] = dir{0, 1}
    directions[2] = dir{-1, 0}
    directions[3] = dir{0, -1}
    perimeter := 0
    current_val := garden[y][x]
    for _, d := range(directions) {
        cx := x + d.x
        cy := y + d.y
        if outsideMap(cx, cy, x_size, y_size) {
            perimeter += 1
        } else if garden[cy][cx] != current_val {
            perimeter += 1
        }
    }
    return perimeter
}

// Fill area with new ints.
func fill(garden [][]uint64, x, y, x_size, y_size int, old, new uint64) {
    if outsideMap(x, y, x_size, y_size) {
        return
    } else {
        if garden[y][x] == old {
            directions := make([]dir, 4)
            directions[0] = dir{1, 0}
            directions[1] = dir{0, 1}
            directions[2] = dir{-1, 0}
            directions[3] = dir{0, -1}
            garden[y][x] = new
            for _, d := range(directions) {
                cx := x + d.x
                cy := y + d.y
                fill(garden, cx, cy, x_size, y_size, old, new)
            }
        }
    }
}

// Input: map of garden plots.
// Output: Total price of fencing (sums of area*perimeter for each region)
// (note that two regions with the same value but not connected are separate)
// Strategy: for each 0 find the number of 9s that can be reached.
func part1(garden [][]uint64) int {
    areas := make(map[uint64]int)
    perimeters := make(map[uint64]int)
    y_size := len(garden)
    x_size := len(garden[0])

    // Need to fill the garden values first so we can distinguish each of the areas.
    new_id := uint64(100)
    for y, row := range(garden) {
        for x, val := range(row) {
            if val < 100 {
                fill(garden, x, y, x_size, y_size, val, new_id)
                new_id += 1
            }
        }
    }
    
    for y, row := range(garden) {
        for x, val := range(row) {
            areas[val] += 1
            perimeters[val] += localPerimeter(garden, x, y, x_size, y_size)
        }
    }
    
    totalCost := 0
    for k, _ := range(areas) {
        // fmt.Println(k, areas[k], perimeters[k])
        totalCost += areas[k] * perimeters[k]
    }
    return totalCost
}

// Turns out: the number of sides in a polygon always equals the number of corners.
func sideCount(garden [][]uint64, x_size, y_size int, shape uint64) int {
    // Look at windows of 2x2 areas and count the times there was just one or three of the value.
    sides := 0
    for x := -1; x < x_size; x++ {
        for y := -1; y < y_size; y++ {
            coords := make([]dir, 4)
            working_vals := make([]bool, 4)
            coords[0] = dir{x, y}
            coords[1] = dir{x+1, y}
            coords[2] = dir{x, y+1}
            coords[3] = dir{x+1, y+1}
            val_count := 0
            for i, c := range(coords) {
                if !outsideMap(c.x, c.y, x_size, y_size) && garden[c.y][c.x] == shape {
                    val_count += 1
                    working_vals[i] = true
                }
            }
            if val_count == 1 || val_count == 3 {
                sides += 1
            } else if val_count == 2 {
                // If the two meet at a corner that would be cool and worth 2 corners.
                if (working_vals[0] && working_vals[3]) || (working_vals[1] && working_vals[2]) {
                    sides += 2
                }
            }
        }
    }
    return sides;
}

// Input: map of garden plots.
// Output: Total price of fencing with the discount (sums of area*side count for each region)
// (note that two regions with the same value but not connected are separate)
// Strategy: for each 0 find the number of 9s that can be reached.
func part2(garden [][]uint64) int {
    areas := make(map[uint64]int)
    side_counts := make(map[uint64]int)
    y_size := len(garden)
    x_size := len(garden[0])

    // Need to fill the garden values first so we can distinguish each of the areas.
    new_id := uint64(100)
    for y, row := range(garden) {
        for x, val := range(row) {
            if val < 100 {
                fill(garden, x, y, x_size, y_size, val, new_id)
                new_id += 1
            }
        }
    }
    
    for _, row := range(garden) {
        for _, val := range(row) {
            areas[val] += 1
            // Only compute this once.
            if side_counts[val] <= 0 {
                side_counts[val] = sideCount(garden, x_size, y_size, val)
            }
        }
    }
    
    totalCost := 0
    for k, _ := range(areas) {
        fmt.Println(k, areas[k], side_counts[k])
        totalCost += areas[k] * side_counts[k]
    }
    return totalCost
}

func main() {
    example := getInput("example.txt")
    garden := getInput("input.txt")

    fmt.Printf("example1: %v, should be 1930\n", part1(example))
    fmt.Printf("example2: %v, should be 1206\n", part2(example))

    fmt.Printf("part1: %v\n", part1(garden))
    fmt.Printf("part2: %v\n", part2(garden))
}
