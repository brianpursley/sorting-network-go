# sorting-network-go
A Go utility for checking and rendering sorting networks

Adapted from [this Python version](https://github.com/brianpursley/sorting-network)

## Usage

You specify a comparison network as a comma-separated list of comparators, where each comparator is formated like `X:Y` where `X` and `Y` are zero-based input indices. For example: `0:1,2:3,0:2,1:3,1:2`

### Load a 16-input comparison network from a file and check if it is a sorting network 
```bash
$ go run sorting-network.go -input examples/16-input.cn -check
```
Output:
```
It is a sorting network!  
```

### Check a 4-input comparison network from stdin
```bash
echo 0:1,2:3,0:2,1:3,1:2 | go run sorting-network.go -check
```
Output:
```
It is a sorting network!  
```

### Check a 4-input comparison network from stdin (Not a sorting network)
```bash
echo 0:1,2:3,0:2,1:3 | go run sorting-network.go -check
```
Output:
```
It is not a sorting network.
```

### Load a 4-input comparison network from a file and render it as an SVG
```bash
go run sorting-network.go -input examples/4-input.cn -svg > examples/4-input.svg
```

### Load a 4-input comparison network from a file and render it as a PNG

You can use rsvg-convert to convert the output from SVG to some other format, like PNG.  rsvg-convert can be installed by using `sudo apt-get install librsvg2-bin` on Ubuntu.

```bash
go run sorting-network.go -input examples/4-input.cn -svg | rsvg-convert > examples/4-input.png
```

## Example sorting networks

### 4-Input

```text
0:1,2:3
0:2,1:3
1:2
```

![4-Input Sorting Network](https://github.com/brianpursley/sorting-network-go/blob/master/examples/4-input.png)

### 5-Input

```text
0:1,3:4
2:4
2:3,1:4
0:3
0:2,1:3
1:2
```

![5-Input Sorting Network](https://github.com/brianpursley/sorting-network-go/blob/master/examples/5-input.png)

### 16-Input

```text
0:1,2:3,4:5,6:7,8:9,10:11,12:13,14:15
0:2,1:3,4:6,5:7,8:10,9:11,12:14,13:15
0:4,1:5,2:6,3:7,8:12,9:13,10:14,11:15
0:8,1:9,2:10,3:11,4:12,5:13,6:14,7:15
5:10,6:9,3:12,13:14,7:11,1:2,4:8
1:4,7:13,2:8,11:14
2:4,5:6,9:10,11:13,3:8,7:12
6:8,10:12,3:5,7:9
3:4,5:6,7:8,9:10,11:12
6:7,8:9
```

![16-Input Sorting Network](https://github.com/brianpursley/sorting-network-go/blob/master/examples/16-input.png)

