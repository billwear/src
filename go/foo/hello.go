package main

import (
    "fmt"
)

func main() {
    // Basic for loop
    for i := 0; i < 5; i++ {
        fmt.Println("Basic loop iteration:", i)
    }
    // Range-based for loop
    fruits := []string{"apple", "banana", "cherry"}
    for index, fruit := range fruits {
        fmt.Printf("Fruit at index %d is %s\n", index, fruit)
    }
}
