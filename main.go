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
  "./predefines"
)

var jumpTable       = predefines.JumpTable
var destTable       = predefines.DestTable
var compTable       = predefines.CompTable
var predefinedTable = predefines.PredefinedTable

var aSymbolRe         = regexp.MustCompile("@")
var lSymbolRe         = regexp.MustCompile("^\\(.*.\\)")
var lSymbolReTrimmer  = regexp.MustCompile("\\(|\\)")
var eqRe              = regexp.MustCompile("=")
var jmpRe             = regexp.MustCompile(";")
var aBitRe            = regexp.MustCompile("\\=.*M")
var jmperRe           = regexp.MustCompile("JMP")

func main() {
  if (len(os.Args) < 2) {
    usage :=  "The Go Hack Assembler\n" +
              "---------------------\n\n" +
              "USAGE: $ ./gohack Add.asm"
    fmt.Println(usage)
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
