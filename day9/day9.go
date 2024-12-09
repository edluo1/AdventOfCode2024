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

func getInput(f string) string {
    file, err := os.Open(f)
    check(err)
    defer file.Close()

    scanner := bufio.NewScanner(file)

    // One-line input.
    for scanner.Scan() {
        return scanner.Text()
    }

    return ""
}

func sameValSlice(size, val int) []int {
    s := make([]int, size)
    for i := range s {
        s[i] = val
    }
    return s
}

// Input: disk map. Alternates between file size and empty space.
// Output: filesystem checksum after moving things into the empty space.
func part1(diskMap string) int {
    var bytemap []int

    checksum := int(0)
    
    // Iterate over string.
    for idx, c := range(diskMap) {
        id := idx/2
        val := int(c - '0')
        var toAdd []int
        if idx%2 == 0 {
            // File.
            toAdd = sameValSlice(val, id)
        } else {
            // Empty space.
            toAdd = sameValSlice(val, -1)
        }
        bytemap = append(bytemap, toAdd...)
    }

    j := 0
    // Iterate until we've gotten through everything.
    for i := len(bytemap)-1; i > j; i-- {
        if bytemap[i] != -1 {
            // Search for value that's -1.
            for j < i && bytemap[j] != -1 {
                if bytemap[j] != -1 {
                    checksum += j * bytemap[j] // Do this now
                }
                j += 1
            }
            if bytemap[j] == -1 {
                // Move it.
                bytemap[j] = bytemap[i]
                bytemap[i] = -1
            }
        }
    }
    // Continue with the rest.
    for j < len(bytemap) && bytemap[j] != -1 {
        checksum += j * bytemap[j] // Continue adding
        j += 1
    }

    return checksum
}

type file struct {
    id int // -1 if free space, id otherwise
    size int
}

type Node struct {
    data file
    next *Node
    prev *Node
}

func connectNode(a *Node, b *Node) {
    if a != nil {
        a.next = b
    }
    if b != nil {
        b.prev = a
    }
}

func moveFile(withinEmptySpace *Node, toMove *Node) {
    beforeNode := toMove.prev
    afterNode := toMove.next
    // Merge the two before and previous nodes if they are empty
    if afterNode != nil && (beforeNode.data.id == afterNode.data.id) {
        // Consider afterNode removed from the list, merged empty space
        connectNode(beforeNode, afterNode.next)
        beforeNode.data.size += afterNode.data.size + toMove.data.size
    } else if afterNode != nil && beforeNode.data.id != -1 && afterNode.data.id == -1 {
        // Empty space before but not after
        connectNode(beforeNode, afterNode)
        afterNode.data.size += toMove.data.size
    } else if afterNode != nil && beforeNode.data.id == -1 && afterNode.data.id != -1 {
        // Empty space after but not before
        connectNode(beforeNode, afterNode)
        beforeNode.data.size += toMove.data.size
    } else if afterNode != nil && beforeNode.data.id != -1 && afterNode.data.id != -1 {
        // No empty space in either way
        newNode := Node { file { -1, toMove.data.size }, nil, nil }
        connectNode(beforeNode, &newNode)
        connectNode(&newNode, afterNode)
    } else if afterNode == nil {
        // Move tail
        connectNode(beforeNode, nil)
    }
    // Connect back piece to front
    connectNode(withinEmptySpace.prev, toMove)

    // Connect front piece to back
    connectNode(toMove, withinEmptySpace)
    withinEmptySpace.data.size -= toMove.data.size
    // Check if link is now empty.
    if withinEmptySpace.data.size <= 0 {
        // Set toMove's next node since an empty space has to be removed
        connectNode(toMove, withinEmptySpace.next)
    }
}

// Input: disk map. Alternates between file size and empty space. Only move if the whole file fits in the space.
// Output: filesystem checksum after moving things into the empty space.
func part2(diskMap string) int {
    var head *Node
    var tail *Node
    
    var currentNode *Node
    // Iterate over string.
    for idx, c := range(diskMap) {
        id := idx/2
        val := int(c - '0')
        var toAdd Node
        if val != 0 {
            if idx%2 == 0 {
                // File.
                toAdd = Node{ file { id, val }, nil, nil }
                if head == nil {
                    head = &toAdd
                }
            } else {
                // Empty space.
                toAdd = Node{ file { -1, val }, nil, nil }
                if tail == nil {
                    tail = &toAdd
                }
            }
            if currentNode != nil {
                if currentNode.data.id == toAdd.data.id {
                    currentNode.data.size += toAdd.data.size
                } else {
                    currentNode.next = &toAdd
                    toAdd.prev = currentNode
                    // Go to next
                    currentNode = &toAdd
                }
            } else {
                currentNode = &toAdd
            }
        }
    }
    tail = currentNode

    nodeToMove := tail
    for nodeToMove.data.id == -1 {
        nodeToMove = nodeToMove.prev
    }

    minFileMoved := 2147483647 // Used to indicate when we should skip
    for nodeToMove != nil {
        previousNode := nodeToMove.prev
        if nodeToMove.data.id != -1 && nodeToMove.data.id < minFileMoved {
            // fmt.Println("moving node", nodeToMove.data.id, "size", nodeToMove.data.size)
            // Move to free space.
            freeArea := head
            spaceExists := true
            // Keep going until we find one
            for freeArea != nil && !(freeArea.data.id == -1 && freeArea.data.size >= nodeToMove.data.size) {
                if freeArea == nil || freeArea.data.id == nodeToMove.data.id {
                    spaceExists = false
                    break
                }
                freeArea = freeArea.next
            }
            if spaceExists {
                if nodeToMove == tail {
                    tail = nodeToMove.prev
                }
                // fmt.Println("found area with size", freeArea.data.size)
                moveFile(freeArea, nodeToMove)
            }
            if minFileMoved > nodeToMove.data.id {
                minFileMoved = nodeToMove.data.id
            }
        }
        nodeToMove = previousNode
    }

    // Iterate starting from the head
    startValue := 0
    endValue := 0
    checksum := int(0)
    currentNode = head
    for currentNode != nil {
        startValue = endValue 
        endValue = startValue + currentNode.data.size
        // Add sum of startValue-endValue (end exclusive), and multiply by the id.
        // fmt.Println(currentNode.data.id, startValue, endValue)
        if currentNode.data.id != -1 {
            totalSum := 0
            for i := startValue; i < endValue; i++ {
                totalSum += i
            }

            checksum += currentNode.data.id * totalSum
        }
        currentNode = currentNode.next
    }

    return checksum
}

func main() {
    example := getInput("example.txt")
    diskMap := getInput("input.txt")

    fmt.Printf("example1: %v, should be 1928\n", part1(example))
    fmt.Printf("example2: %v, should be 2858\n", part2(example))

    fmt.Printf("part1: %v\n", part1(diskMap))
    fmt.Printf("part2: %v\n", part2(diskMap))
}
