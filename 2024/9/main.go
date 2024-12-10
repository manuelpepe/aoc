package main

import (
	"fmt"
	"os"
	"slices"
)

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	for ix := range data {
		data[ix] -= '0'
	}

	fmt.Printf("Checksum: %d\n", sol1(slices.Clone(data)))
	fmt.Printf("Checksum defrag: %d\n", sol2(data))
}

func sol1(data []byte) int {
	head := 1 // skip fileid = 0
	tail := len(data) - 1

	acc := 0
	blockIx := int(data[0]) // skip blocks for fileid = 0

	for {
		// 	if in free space (head % 2 == 1)
		// 		for each free space
		//			tailId = floor(tail / 2)
		// 			acc += tailId * blockIx
		// 	if in used space (head % 2 == 0)
		//  	for each space
		//			curId = head / 2
		//  		acc += id * blockIx

		if head%2 == 1 {
			// free space, add to acc from tail
			for j := 0; j < int(data[head]); j++ {
				if data[tail] == 0 {
					tail -= 2 // advance tail if out of data
				}

				if tail <= head {
					break // stop if wrapped
				}

				tailId := int(tail / 2)
				acc += tailId * blockIx

				data[tail] -= 1

				blockIx++
			}
		} else {
			// used space, just add to acc
			for j := 0; j < int(data[head]); j++ {
				curId := int(head / 2)
				acc += curId * blockIx

				blockIx++
			}
		}

		if head >= tail {
			break
		}

		head++

	}

	return acc
}

func sol2(data []byte) int {
	for ix := range data {
		if ix%2 == 1 {
			continue
		}
		// original size encoding:
		// max value for elements in data is 9 (i.e 4 bytes), this value
		// is encoded past the first 4 bytes, which allows the following algorithm
		// to decrease the available slots (first 4bytes) while keeping the original value (next 4 bytes)
		// intact without using additional memory.
		data[ix] += data[ix] << 4
	}

	acc := 0
	blockIx := 0

	ix := 0

	for ix < len(data) {
		if ix%2 == 1 {
			// free space
			fileIx, found := firstFileOfSize(data, data[ix], ix)
			if !found {
				blockIx += int(data[ix] & 0b1111)
				ix++
				continue
			}

			fileId := int(fileIx / 2)

			for data[fileIx]&0b1111 > 0 {
				data[fileIx] -= 1
				data[ix] -= 1

				acc += fileId * blockIx
				blockIx++
			}

			if data[ix] <= 0 {
				ix++
			}

		} else {
			// used space, just add to acc
			allocated := data[ix] & 0b1111
			for j := 0; j < int(allocated); j++ {
				curId := int(ix / 2)
				acc += curId * blockIx

				blockIx++
			}

			if allocated == 0 {
				blocks := data[ix] >> 4
				for j := 0; j < int(blocks); j++ {
					blockIx++
				}
			}

			ix++
		}
	}

	return acc
}

func firstFileOfSize(data []byte, size byte, limit int) (int, bool) {
	for ix := len(data) - 1; ix > limit; ix -= 2 {
		if data[ix]&0b1111 > size || data[ix]&0b1111 == 0 {
			continue
		}

		return ix, true
	}

	return 0, false
}
