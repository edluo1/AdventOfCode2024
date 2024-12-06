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

const UP byte = 1
const RIGHT byte = 2
const DOWN byte = 4
const LEFT byte = 8

type guard struct {
    facing byte
    x int
    y int
}

func getGuardMap(f string) [][]byte {
    file, err := os.Open(f)
    check(err)
    defer file.Close()

    scanner := bufio.NewScanner(file)

    // Make sure to copy each line so that it doesn't break
    var guardMap [][]byte
    for scanner.Scan() {
        line := make([]byte, len(scanner.Bytes()))
        copy(line, scanner.Bytes())
        guardMap = append(guardMap, line)
    }

    return guardMap
}

func copyMap(m [][]byte) [][]byte {
    c := make ([][]byte, len(m))
    for i := range m {
        c[i] = make([]byte, len(m[i]))
        copy(c[i], m[i])
    }
    return c
}

func outsideMap(g guard, x_size int, y_size int) bool {
    return g.x < 0 || g.y < 0 || g.x >= x_size || g.y >= y_size
}

func outsideMap_c(x int, y int, x_size int, y_size int) bool {
    return x < 0 || y < 0 || x >= x_size || y >= y_size
}

func checkAhead(guardMap [][]byte, g *guard) bool {
    y_size := len(guardMap)
    x_size := len(guardMap[0])
    x := g.x
    y := g.y
    switch g.facing {
    case UP:
        y -= 1
    case DOWN:
        y += 1
    case LEFT:
        x -= 1
    case RIGHT:
        x += 1
    }
    if outsideMap_c(x, y, x_size, y_size) {
        return true
    }
    if guardMap[y][x] == '#' {
        return false
    }
    return true
}

// Turn the guard until no longer facing a wall.
func turn_guard(guardMap [][]byte, g *guard) {
    safe_ahead := checkAhead(guardMap, g)
    for !safe_ahead {
        switch g.facing {
        case UP:
            g.facing = RIGHT
        case DOWN:
            g.facing = LEFT
        case LEFT:
            g.facing = UP
        case RIGHT:
            g.facing = DOWN
        }
        safe_ahead = checkAhead(guardMap, g)
    }
}

// Move the guard forward.
func move_guard(g *guard, forward bool) {
    mult := 1
    if !forward {
        mult = -1
    }
    switch g.facing {
    case UP:
        g.y -= 1 * mult
    case DOWN:
        g.y += 1 * mult
    case LEFT:
        g.x -= 1 * mult
    case RIGHT:
        g.x += 1 * mult
    }
}

// Input: guard's map
// Guard starts facing north, and you need to find every square he visits until he leaves.
func part1(guardMap [][]byte) int {
    y_size := len(guardMap)
    x_size := len(guardMap[0])
    // Search for guard
    g := guard{UP, -1, -1}
    for y, row := range(guardMap) {
        for x, val := range(row) {
            if val == '^' {
                g.x = x
                g.y = y
                break
            }
        }
        if g.x != -1 {
            break
        }
    }

    // Simulate the guard's walk.
    visit := 0
    for !outsideMap(g, x_size, y_size) {
        // Mark this area if it's not visited.
        if guardMap[g.y][g.x] == '^' || guardMap[g.y][g.x] == '.' {
            visit += 1
        }
        guardMap[g.y][g.x] = 'X'

        turn_guard(guardMap, &g)
        move_guard(&g, true)
    }

    return visit
}

type obstacle struct {
    facing byte // Set facing to the way the guard should face when he first comes across it.
    x int
    y int
}

// Input: guard's map
// Guard starts facing north. Figure out routes that would force him to circle around.
// You only get one obstacle.
// Iterate over all areas the guard walked onto and detect the loop by calculating which way it went.
func part2(guardMap [][]byte) int {
    clearMap := copyMap(guardMap) // Using this to reset
    y_size := len(guardMap)
    x_size := len(guardMap[0])
    // Search for guard
    g := guard{UP, -1, -1}
    for y, row := range(guardMap) {
        for x, val := range(row) {
            if val == '^' {
                g.x = x
                g.y = y
                break
            }
        }
        if g.x != -1 {
            break
        }
    }

    start_x := g.x
    start_y := g.y

    // Simulate the guard's walk first.
    obstacles := make([]obstacle, 0)
    for !outsideMap(g, x_size, y_size) {
        // Mark this area with a "facing" byte.

        if guardMap[g.y][g.x] == '^' || guardMap[g.y][g.x] == '.' {
            guardMap[g.y][g.x] = 0
            // Add a possible mark here if it isn't the start.
            if (!(g.x == start_x && g.y == start_y)) {
                obstacles = append(obstacles, obstacle{g.facing, g.x, g.y})
            }
        }

        guardMap[g.y][g.x] = 'X'

        turn_guard(guardMap, &g)
        move_guard(&g, true)
    }

    // Cleanup map and go over list of obstacles to see which cause loops
    loop_obstacles := 0
    for _, obs := range(obstacles) {
        // Cleanup map
        guardMap = copyMap(clearMap)

        // Start guard just behind the obstacle.
        // Since we always move forward the guard is always in the right place.
        guardMap[obs.y][obs.x] = '#'
        g.facing = obs.facing
        g.x = obs.x
        g.y = obs.y
        move_guard(&g, false)
        loop_found := false

        for !outsideMap(g, x_size, y_size) && !loop_found {
            // Mark this area with a "facing" byte.
            empty_space := guardMap[g.y][g.x] == '^' || guardMap[g.y][g.x] == '.'
            if empty_space {
                guardMap[g.y][g.x] = 0
            } else if !empty_space && guardMap[g.y][g.x] & g.facing != 0 {
                loop_found = true
                break
            }

            guardMap[g.y][g.x] |= g.facing

            turn_guard(guardMap, &g)
            move_guard(&g, true)
        }
        if loop_found {
            loop_obstacles += 1
        }
    }

    return loop_obstacles
}

func main() {
    example := getGuardMap("example.txt")
    guardMap := getGuardMap("input.txt")

    fmt.Printf("example1: %v, should be 41\n", part1(copyMap(example)))
    fmt.Printf("example2: %v, should be 6\n", part2(copyMap(example)))

    fmt.Printf("part1: %v\n", part1(copyMap(guardMap)))
    fmt.Printf("part2: %v\n", part2(copyMap(guardMap)))
}
