package main

import (
  "fmt"
  "log"
  "bufio"
  "bytes"
  "os"
  "strings"
  "strconv"
  "regexp"
)

var jumpTable = map[string]string{
  "null": "000",
  "JGT" : "001",
  "JEQ" : "010",
  "JGE" : "011",
  "JLT" : "100",
  "JNE" : "101",
  "JLE" : "110",
  "JMP" : "111",
}

var destTable = map[string]string{
  "null": "000",
  "M"   : "001",
  "D"   : "010",
  "MD"  : "011",
  "A"   : "100",
  "AM"  : "101",
  "AD"  : "110",
  "AMD" : "111",
}

var compTable = map[string]string{
  "0"   : "101010",
  "1"   : "111111",
  "-1"  : "111010",
  "D"   : "001100",
  "A"   : "110000",
  "M"   : "110000",
  "!D"  : "001101",
  "!A"  : "110001",
  "!M"  : "110001",
  "-D"  : "001111",
  "-A"  : "110011",
  "-M"  : "110011",
  "D+1" : "011111",
  "1+D" : "011111",
  "1+A" : "110111",
  "A+1" : "110111",
  "1+M" : "110111",
  "M+1" : "110111",
  "1-D" : "001110",
  "D-1" : "001110",
  "1-A" : "110010",
  "A-1" : "110010",
  "1-M" : "110010",
  "M-1" : "110010",
  "D+A" : "000010",
  "A+D" : "000010",
  "D+M" : "000010",
  "M+D" : "000010",
  "A-D" : "000111",
  "D-A" : "000111",
  "M-D" : "010011",
  "D-M" : "010011",
  "D&A" : "000000",
  "A&D" : "000000",
  "D&M" : "000000",
  "M&D" : "000000",
  "D|A" : "010101",
  "A|D" : "010101",
  "D|M" : "010101",
  "M|D" : "010101",
}
var aSymbolRe = regexp.MustCompile("@")
var eqRe = regexp.MustCompile("=")
var jmpRe = regexp.MustCompile(";")
var aBitRe = regexp.MustCompile("M")
var jmperRe = regexp.MustCompile("JMP")
func main() {
  asmString := readFileAndPrepare("MaxL.asm")
  createProgram(asmString)
}

// remove comments and whitespace
func readFileAndPrepare(path string) string {
  var buffer bytes.Buffer
  re := regexp.MustCompile("^//.*")
  f, err := os.Open(path)
  if err != nil {
    log.Fatal("error opening file: %v\n",err)
  }

  scanner := bufio.NewScanner(f)
  for scanner.Scan() {
    line := scanner.Text()
    line = re.ReplaceAllString(line, "")
    if (len(line) > 0) {
      buffer.WriteString(line+"\n")
    }
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  return buffer.String()
}

func createProgram(asmString string) {
  //cTable := map[string]string{}
  //aTable := map[string]string{}
  scanner := bufio.NewScanner(bytes.NewBufferString(asmString))
  for scanner.Scan() {
    line := scanner.Text()
    if (len(line) > 0) {
      if (aSymbolRe.FindStringIndex(string(line)) == nil) {
        if (eqRe.FindStringIndex(string(line)) != nil && jmpRe.FindStringIndex(string(line)) != nil) {
          //aBit := getABit(string(line[1:]))
          //fmt.Printf("111%s%s%s%s\n",aBit,compTable[line[:1]],destTable["null"],jumpTable[line[2:]])
        } else if (eqRe.FindStringIndex(string(line)) == nil && jmpRe.FindStringIndex(string(line)) != nil) {
          aBit := getABit(string(line[1:]))
          fmt.Printf("111%s%s%s%s\n",aBit,compTable[line[:1]],destTable["null"],jumpTable[line[2:]])
        } else if (eqRe.FindStringIndex(string(line)) != nil && jmpRe.FindStringIndex(string(line)) == nil) {
          aBit := getABit(string(line[1:]))
          fmt.Printf("111%s%s%s%s\n",aBit,compTable[line[2:]],destTable[line[0:1]],jumpTable["null"])
        }
      } else if (aSymbolRe.FindStringIndex(string(line)) != nil) {
        fmt.Println(asBinaryString(string(line[1:])))
        //aTable[string(asmString[i])] = asBinaryString(string(asmString[i]))
      }
    }
  }
}

func getABit(aOrC string) string {
  if (len(aBitRe.FindString(aOrC)) > 0 && len(jmpRe.FindString(aOrC)) == 0) {
    return "1"
  } else {
    return "0"
  }
}
func asBinaryString(numb string) string {
  convNumb,_ := strconv.ParseInt(numb, 10, 64)
  a := 16-len(strconv.FormatInt(convNumb, 2))
  newBinary := strings.Repeat("0", a)+strconv.FormatInt(convNumb, 2)
  return newBinary
}
