package main

import (
	// "flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// var testflag = flag.String("testflag", "no value", "test aja dulu")

var cpuprof_file *os.File
var memprof_file *os.File
var cpuprof_filename string = "cpuprof.prof"
var memprof_filename string = "memprof.prof"

// func heavyTask() {
// 	end := time.Now().Add(1 * time.Minute)

// 	// Busy loop to simulate CPU work until the target time is reached.
// 	for time.Now().Before(end) {
// 		// Perform some work
// 		for i := 0; i < 1000000; i++ {
// 			_ = i * i // Simple computation to keep the CPU busy
// 		}
// 	}
// }

func SimpleMemoryAllocation(numAllocations int, sizePerAllocation int) {

	var data []byte

	for i := 0; i < numAllocations; i++ {
		data = make([]byte, sizePerAllocation)
		data[0] = byte(i % 256) // Initialize one element to ensure allocation

		// Optionally print memory usage every 1000 allocations
		// if (i+1)%1000 == 0 {
		// 	var memStats runtime.MemStats
		// 	runtime.ReadMemStats(&memStats)
		// 	fmt.Printf("Allocations: %d, Current memory usage: %d KB\n", i+1, memStats.Alloc/1024)
		// }
	}
}

func HeavyRecursiveFibonacci(n int) int {

	if n <= 1 {
		return n
	}
	return HeavyRecursiveFibonacci(n-1) + HeavyRecursiveFibonacci(n-2)
}

func start_pprof() {

	cpuprof_file, _ = os.Create(cpuprof_filename)
	memprof_file, _ = os.Create(memprof_filename)

	numAllocations := 100        // Number of allocations to make
	sizePerAllocation := 1 << 20 // Size of each allocation in bytes (1 KB)

	pprof.StartCPUProfile(cpuprof_file)

	// heavyTask()

	_ = HeavyRecursiveFibonacci(40)

	SimpleMemoryAllocation(numAllocations, sizePerAllocation)
}

type DataBlock struct {
	ID   int
	Data []byte
}

// AllocateMemory function demonstrates extensive memory allocation
func AllocateMemory(blockCount, blockSize int) []DataBlock {
	allocatedBlocks := make([]DataBlock, 0, blockCount)

	for i := 0; i < blockCount; i++ {
		block := DataBlock{
			ID: i,
			// Data: make([]byte, blockSize),
			Data: make([]byte, 0, blockSize),
		}
		// Filling data with a pattern
		for j := 0; j < blockSize; j++ {

			// block.Data[j] = byte((i + j) % 256)
			block.Data = append(block.Data, byte((i+j)%256))

		}
		allocatedBlocks = append(allocatedBlocks, block)

		// Optional: Print memory usage periodically for monitoring
		if (i+1)%1000 == 0 {
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			fmt.Printf("Allocated blocks: %d, Memory usage: %d KB\n", i+1, memStats.Alloc/1024)
		}
	}
	return allocatedBlocks
}

func main() {

	start_pprof()

	defer pprof.StopCPUProfile()

	defer cpuprof_file.Close()
	defer memprof_file.Close()

	memprof_file, _ := os.Create(memprof_filename)

	blockCount := 10000    // Number of blocks to allocate
	blockSize := 1024 * 10 // Size of each block in bytes (10KB)

	fmt.Println("Starting memory allocation...")
	blocks := AllocateMemory(blockCount, blockSize)
	fmt.Printf("Allocated %d blocks of %d bytes each.\n", len(blocks), blockSize)

	pprof.WriteHeapProfile(memprof_file)

	// Simulating some processing time
	time.Sleep(5 * time.Second)

	// Manually clearing allocated memory (for testing purposes)
	blocks = nil
	runtime.GC() // Forcing garbage collection

	// Printing memory usage after deallocation
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("Memory usage after deallocation: %d KB\n", memStats.Alloc/1024)
}
