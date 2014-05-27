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
  "M-D" : "000111",
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

var predefinedTable = map[string]string{
  "SP"  : "0",
  "LCL" : "1",
  "ARG" : "2",
  "THIS": "3",
  "THAT": "4",
  "R0"  : "0",
  "R1"  : "1",
  "R2"  : "2",
  "R3"  : "3",
  "R4"  : "4",
  "R5"  : "5",
  "R6"  : "6",
  "R7"  : "7",
  "R8"  : "8",
  "R9"  : "9",
  "R10" : "10",
  "R11" : "11",
  "R12" : "12",
  "R13" : "13",
  "R14" : "14",
  "R15" : "15",
  "SCREEN": "16384",
  "KBD" : "24576",
}

var aSymbolRe = regexp.MustCompile("@")
var lSymbolRe = regexp.MustCompile("^\\(.*.\\)")
var lSymbolReTrimmer = regexp.MustCompile("\\(|\\)")
var eqRe = regexp.MustCompile("=")
var jmpRe = regexp.MustCompile(";")
var aBitRe = regexp.MustCompile("\\=.*M")
var jmperRe = regexp.MustCompile("JMP")

func main() {
  if (len(os.Args) < 2) {
    fmt.Println("USAGE: ...")
  } else {
    asmString := readFileAndPrepare(os.Args[1])
    createProgram(asmString)
  }
}

// remove comments
func readFileAndPrepare(path string) string {
  var buffer bytes.Buffer
  re := regexp.MustCompile("//.*")
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
  lTable := map[string]string{}
  aTable := map[string]string{}
  scannerL := bufio.NewScanner(bytes.NewBufferString(asmString))
  scannerAC := bufio.NewScanner(bytes.NewBufferString(asmString))
  lineCounter := 0
  // scanning for labels
  for scannerL.Scan() {
    line := strings.TrimSpace(string(scannerL.Text()))
    lineCounter++
    if (lSymbolRe.FindStringIndex(line) != nil) {
      tempL := lSymbolReTrimmer.ReplaceAllString(line, "")
      if (lTable[tempL] == "") {
        lineCounter--
        lTable[tempL] = asBinaryString(strconv.Itoa(lineCounter))
      }
    }
  }
  // scanning the whole program
  for scannerAC.Scan() {
    line := strings.TrimSpace(string(scannerAC.Text()))
    if (len(line) > 0) {
      if (aSymbolRe.FindStringIndex(line) == nil) {
        if (eqRe.FindStringIndex(line) != nil && jmpRe.FindStringIndex(line) != nil) {
          cmp := eqRe.FindStringIndex(line)[1]
          jmp := jmpRe.FindStringIndex(line)[1]
          jpm := jmpRe.FindStringIndex(line)[0]
          des := eqRe.FindStringIndex(line)[0]
          aBit := getABit(line[eqRe.FindStringIndex(line)[0]:])
          fmt.Printf("111%s%s%s%s\n",aBit,compTable[line[cmp:jpm]],destTable[line[0:des]],jumpTable[line[jmp:]])
        } else if (eqRe.FindStringIndex(line) == nil && jmpRe.FindStringIndex(line) != nil) {
          cmp := jmpRe.FindStringIndex(line)[0]
          jmp := jmpRe.FindStringIndex(line)[1]
          aBit := getABit(line[jmpRe.FindStringIndex(line)[0]:])
          fmt.Printf("111%s%s%s%s\n",aBit,compTable[line[:cmp]],destTable["null"],jumpTable[line[jmp:]])
        } else if (eqRe.FindStringIndex(line) != nil && jmpRe.FindStringIndex(line) == nil) {
          cmp := eqRe.FindStringIndex(line)[1]
          des := eqRe.FindStringIndex(line)[0]
          aBit := getABit(line[eqRe.FindStringIndex(line)[0]:])
          fmt.Printf("111%s%s%s%s\n",aBit,compTable[line[cmp:]],destTable[line[0:des]],jumpTable["null"])
        }
      } else if (aSymbolRe.FindStringIndex(line) != nil) {
        sym := line[aSymbolRe.FindStringIndex(line)[1]:]
        if (predefinedTable[sym]!= "") {
          fmt.Println(asBinaryString(predefinedTable[sym]))
        } else if (lTable[sym] != "") {
          fmt.Println(lTable[sym])
        } else if isInt(sym) {
          fmt.Println(asBinaryString(sym))
        } else {
          if (aTable[sym] == "") {
            aTable[sym] = asBinaryString(strconv.Itoa(len(aTable)+16))
          }
          fmt.Println(aTable[sym])
        }
      }
    }
  }
}

func addOrGetFromLabelTable(table map[string]string, line string) string {
  return "TEST"
}

func isInt(strOrNum string) bool {
  _, err := strconv.ParseInt(strOrNum, 10, 64)
  if err != nil {
    return false
  } else {
    return true
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
