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

func getInput(f string) [][]byte {
    file, err := os.Open(f)
    check(err)
    defer file.Close()

    scanner := bufio.NewScanner(file)

    // Make sure to copy each line so that it doesn't break
    var antennaMap [][]byte
    for scanner.Scan() {
        line := make([]byte, len(scanner.Bytes()))
        copy(line, scanner.Bytes())
        antennaMap = append(antennaMap, line)
    }

    return antennaMap
}

type antenna struct {
    freq byte
    x int
    y int
}

func get_antinodes(a1 antenna, a2 antenna) (antenna, antenna) {
    dx := a1.x - a2.x
    dy := a1.y - a2.y
    an1 := antenna {a1.freq, a1.x + dx, a1.y + dy}
    an2 := antenna {a2.freq, a2.x - dx, a2.y - dy}
    return an1, an2
}

// Setting these up so we don't start at negative
func mod(a, b int) int {
    return (a % b + b) % b
}

func dist(a, b int) int {
    if a > b {
        return a-b
    } else {
        return b-a
    }
}

// New inputs for x_size and y_size so we don't leave the map
func get_all_antinodes(a1 antenna, a2 antenna, x_size int, y_size int) []antenna {
    dx := a1.x - a2.x
    dy := a1.y - a2.y
    // Just in case
    if dx == 0 && dy == 0 {
        return make([]antenna, 0)
    }

    start_x := a1.x
    start_y := a1.y
    for !outsideMap_coord(start_x, start_y, x_size, y_size) {
        start_x -= dx
        start_y -= dy
    }
    // Go back, we've left the map.
    start_x += dx
    start_y += dy

    // Now add each of the antinodes.
    x := start_x
    y := start_y
    var antinodes []antenna
    for !outsideMap_coord(x, y, x_size, y_size) {
        an := antenna {a1.freq, x, y}
        x += dx
        y += dy
        antinodes = append(antinodes, an)
    }

    return antinodes
}

func outsideMap_coord(x int, y int, x_size int, y_size int) bool {
    return x < 0 || y < 0 || x >= x_size || y >= y_size
}

func outsideMap(a antenna, x_size int, y_size int) bool {
    return a.x < 0 || a.y < 0 || a.x >= x_size || a.y >= y_size
}

// Input: map of antennas.
// Antennas with the same frequency cause antinodes at two collinear points (even at other antennas).
// Count how many unique locations have an antinode.
// All antinodes must be on the map.
func part1(antennaMap [][]byte) int {
    antennas_by_freq := make([][]antenna, 75)
    antinode_count := 0
    y_size := len(antennaMap)
    x_size := len(antennaMap[0])

    // Process map.
    for y, row := range(antennaMap) {
        for x, val := range(row) {
            if val != '.' {
                freq_id := val - '0'
                antennas_by_freq[freq_id] = append(antennas_by_freq[freq_id], antenna{freq_id, x, y})
            }
        }
    }

    // Set up antinodes.
    has_antinode := make([]bool, x_size*y_size)
    for _, antennas := range(antennas_by_freq) {
        // Go over each pair of antennas and set up the ones on the left and right.
        for i := 0; i < len(antennas); i++ {
            for j := i+1; j < len(antennas); j++ {
                a1 := antennas[i]
                a2 := antennas[j]
                an1, an2 := get_antinodes(a1, a2)
                if !outsideMap(an1, x_size, y_size) && has_antinode[an1.y*x_size+an1.x] == false {
                    antinode_count += 1
                    has_antinode[an1.y*x_size+an1.x] = true
                }
                if !outsideMap(an2, x_size, y_size) && has_antinode[an2.y*x_size+an2.x] == false {
                    antinode_count += 1
                    has_antinode[an2.y*x_size+an2.x] = true
                }
            }
        }
    }

    return antinode_count
}

// Input: map of antennas.
// Antennas with the same frequency cause antinodes at ALL collinear equidistant points
// (even at other antennas). This also includes the antenna locations themselves.
// Count how many unique locations have an antinode.
// All antinodes must be on the map.
func part2(antennaMap [][]byte) int {
    antennas_by_freq := make([][]antenna, 75)
    antinode_count := 0
    y_size := len(antennaMap)
    x_size := len(antennaMap[0])

    // Process map.
    for y, row := range(antennaMap) {
        for x, val := range(row) {
            if val != '.' {
                freq_id := val - '0'
                antennas_by_freq[freq_id] = append(antennas_by_freq[freq_id], antenna{freq_id, x, y})
            }
        }
    }

    // Set up antinodes.
    has_antinode := make([]bool, x_size*y_size)
    for _, antennas := range(antennas_by_freq) {
        // Go over each pair of antennas and set up the ones on the left and right.
        for i := 0; i < len(antennas); i++ {
            for j := i+1; j < len(antennas); j++ {
                a1 := antennas[i]
                a2 := antennas[j]
                all_antinodes := get_all_antinodes(a1, a2, x_size, y_size)
                for _, an := range(all_antinodes) {
                    if !outsideMap(an, x_size, y_size) && has_antinode[an.y*x_size+an.x] == false {
                        antinode_count += 1
                        has_antinode[an.y*x_size+an.x] = true
                    }
                }
            }
        }
    }

    return antinode_count
}

func main() {
    example := getInput("example.txt")
    antennaMap := getInput("input.txt")

    fmt.Printf("example1: %v, should be 14\n", part1(example))
    fmt.Printf("example2: %v, should be 34\n", part2(example))

    fmt.Printf("part1: %v\n", part1(antennaMap))
    fmt.Printf("part2: %v\n", part2(antennaMap))
}
