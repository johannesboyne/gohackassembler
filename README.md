#GO Hack Assembler

This is a golang hack assembler, build for the [CSOS Class](https://www.wi.uni-muenster.de/student-affairs/bachelor-master-lectures/168427) @ [WWU](http://uni-muenster.de) by [Johannes Boyne](https://github.com/johannesboyne).

####Purpose and Code-Structure

The code is written functional at most but sometimes with a combination of OO principles.
This provides two major opportunities: avoiding states on the one hand but relying on some on the other usefull OO concepts.

#####Design

The program loads a file as string and removes commets. Next it scans the file using golangs `bufio.NewScanner` this provides an easy way of stepping through the string line by line. This process has to be done twice:

 1. build symbol table (as proposed in *Chapter 6 - Nand2Tetris*)
 2. build binary translation

###Usage

```shell
$ go run main.go Add.asm

0000000000000010
1110110000010000
0000000000000011
1110000010010000
0000000000000000
1110001100001000
```

or use the build

```shell
$ ./gohack Add.asm

0000000000000010
1110110000010000
0000000000000011
1110000010010000
0000000000000000
1110001100001000
```

###Testing

Use the `test.sh` it will go through some of the sample `.asm` files with the assembler and compares them using the UNIX `diff` tool with the provided `.hack` counterparts.

```shell
$ ./test.sh
```

###LICENSE

MIT

###Copyright

Johannes Boyne (2014)
