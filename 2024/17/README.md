### Solution 2:

This part was interesting but I haven't managed to get it done.

After experimenting a bit and looking at some tips on the subreddit I managed to reach these hypothesis:

##### 1. 

There seems to be some relation between each octet in the input number and the output it produces.

Here are some inputs and the outputs they produce

```
o30 = 3,0
o300 = 5,3,0
o3002 = 5,5,3,0
o30020 = 3,5,5,3,0
o300202 = 0,3,5,5,3,0
o3002511 = 6,0,3,5,5,3,0
o30025110 = 1,6,0,3,5,5,3,0
o300251105 = 3,1,6,0,3,5,5,3,0
```

Each octet adds a digit to the output.

Also, with this small sample, it seems that each octet could largely remain the same after being when shifted 4 octets (3002, 2511)
altho the second group of 4 octets changed mid calculation (from 202x to 2511) which surely will cause troubles.

I tried to generate each 4 octet group alone to & them together so:

```
$ go run main.go input-1.txt 5,5,3,0
output for A=1538 (o3002) (b11000000010): 5,5,3,0
$ go run main.go input-1.txt 1,6,0,3
output for A=2632 (o5110) (b101001001000): 1,6,0,3
$ go run main.go input-1.txt 7,5,4,3
output for A=0 (o0) (b0): 7,5,4,3  # err
$ go run main.go input-1.txt 2,4,1,5
output for A=2203 (o4233) (b100010011011): 2,4,1,5
```

I couldn't generate the third octet group...
I still tried to & the numbers with some empty data for the third:

```
3002 5110 0000 4233 octal = 105734712854683 decimal
```

which generates

```
$ go run main.go input-1.txt 105734712854683
output for A=105734712854683 (o3002511000004233) (b11000000010101001001000000000000000100010011011): 2,4,1,5,3,3,3,1,1,6,0,3,5,5,3,0
```

and that's almost the input, wrong by the 4 instructions corresponding to the third octet.

I would like to bruteforce the third octet but, as mentioned already:

```
$ go run main.go input-1.txt 7,5,4,3
output for A=0 (o0) (b0): 7,5,4,3  # err
```

so I'm kinda stuck here.

I tried to generate the last two groups together without luck:

```
$ go run main.go input-1.txt 2,4,1,5,7,5,4,3
output for A=0 (o0) (b0): 2,4,1,5,7,5,4,3
```

and smaller subsets:

```
$ go run main.go input-1.txt 4,3
output for A=0 (o0) (b0): 4,3
$ go run main.go input-1.txt 5,4,3
output for A=0 (o0) (b0): 5,4,3
$ go run main.go input-1.txt 7,5,4,3
output for A=0 (o0) (b0): 7,5,4,3
$ go run main.go input-1.txt 5,7,5,4,3
output for A=0 (o0) (b0): 5,7,5,4,3
$ go run main.go input-1.txt 1,5,7,5,4,3
output for A=0 (o0) (b0): 1,5,7,5,4,3
$ go run main.go input-1.txt 4,1,5,7,5,4,3
output for A=0 (o0) (b0): 4,1,5,7,5,4,3
```

Lastly i tried to recalculate the smallest input for a four octal output when failing to calculate the next one, this confirmed that there are multiple ways
to generate each output. This are the numbers that I found:

```
o3002 (1538)
o30025110 (6302280)
o30025113 (6302283)
o30025114 (6302284)
o30025367 (6302455)
o3004 (1540)
o3006 (1542)
o3007 (1543)
o3042 (1570)
o3044 (1572)
o3047 (1575)
o3062 (1586)
o3064 (1588)
o3072 (1594)
o3102 (1602)
o3104 (1604)
o3106 (1606)
o3107 (1607)
o3710 (1992)
```

I managed to get some variations up to the second octal but failed to go past it and eventually ran out of options for the first one.


##### *time passes...*

I came back after a couple of days and tried to do a simpler search. 
Starting from the leftmost octet increase it by 1 until the generated rightmost instruction set matches the what's expected. 
This is the final code:

```go
func sol2(m Machine) int {
	// get all instructions
	nums := make([]byte, 0)
	for i := range m.intructions {
		nums = append(nums, byte(m.intructions[i].opcode))
		nums = append(nums, byte(m.intructions[i].operand))
	}

	octIx := 15            // hardcoded for the 16 instruction input
	n := 01000000000000000 // start with leftmost octet at 1

	for {
		if octIx < 0 {
			break
		}

		exp := join(nums[octIx:])

		wa := with_a(m, n)
		spl := strings.Split(wa, ",")[octIx:]
		res := strings.Join(spl, ",")

		// fmt.Printf("0o%o (%d)-> %+v == %+v\n", n, octIx, res, exp)

		if res == exp {
			octIx -= 1
			continue
		}

		n += 01 << (3 * octIx) // add 1 to the leftmost unmatching octal
	}

	return n
}
```

which seems to work quite well.