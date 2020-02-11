package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/andrewarchi/graph"
	"github.com/andrewarchi/nebula/analysis"
	"github.com/andrewarchi/nebula/codegen"
	"github.com/andrewarchi/nebula/ir"
	"github.com/andrewarchi/nebula/ws"
	"llvm.org/llvm/bindings/go/llvm"
)

var (
	mode            string
	maxStackLen     uint
	maxCallStackLen uint
	maxHeapBound    uint
	noFold          bool
	packed          bool

	modeActions = map[string]func(*ws.Program){
		"":       emitIR,
		"ws":     emitWS,
		"wsa":    emitWSA,
		"ir":     emitIR,
		"llvm":   emitLLVM,
		"dot":    printDOT,
		"matrix": printMatrix,
	}
)

const usageHeader = `Nebula is a compiler for stack-based languages targeting LLVM IR.

Usage:

	%s [options] <program>

Options:

`

const usageFooter = `
Examples:

	%s -mode=ir programs/pi.out.ws > pi.nir
	%s -mode=llvm programs/ascii4.out.ws > ascii4.ll
	%s -mode=llvm -heap=400000 programs/interpret.out.ws > interpret.ll
	%s -mode=dot programs/interpret.out.ws | dot -Tpng > graph.png

`

const modeUsage = `Output mode:
* ws      emit Whitespace syntax
* wsa     emit Whitespace AST
* ir      emit Nebula IR (default)
* llvm    emit LLVM IR
* dot     print control flow graph as Graphviz DOT digraph
* matrix  print control flow graph as Unicode matrix`

func init() {
	flag.Usage = usage
	flag.StringVar(&mode, "mode", "", modeUsage)
	flag.UintVar(&maxStackLen, "stack", codegen.DefaultMaxStackLen, "Maximum stack length for LLVM codegen")
	flag.UintVar(&maxCallStackLen, "calls", codegen.DefaultMaxCallStackLen, "Maximum call stack length for LLVM codegen")
	flag.UintVar(&maxHeapBound, "heap", codegen.DefaultMaxHeapBound, "Maximum heap address bound for LLVM codegen")
	flag.BoolVar(&noFold, "nofold", false, "Disable constant folding")
	flag.BoolVar(&packed, "packed", false, "Enable bit packed format for input file")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "No program provided.")
		incorrectUsage()
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Too many arguments provided.")
		incorrectUsage()
	}

	modeAction, ok := modeActions[mode]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unrecognized mode: %s\n", mode)
		incorrectUsage()
	}

	filename := args[0]
	program, err := ws.LexProgram(filename, packed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lex error: %v\n", err)
		os.Exit(1)
	}
	modeAction(program)
}

func usage() {
	cmd := os.Args[0]
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, usageHeader, cmd)
	flag.PrintDefaults()
	fmt.Fprintf(w, usageFooter, cmd, cmd, cmd, cmd)
}

func incorrectUsage() {
	fmt.Fprintf(os.Stderr, "Run %s -help for usage.\n", os.Args[0])
	os.Exit(2)
}

func convertSSA(p *ws.Program) *ir.Program {
	program, err := p.ConvertSSA()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if _, ok := err.(*ir.ErrorRetUnderflow); !ok {
			os.Exit(1)
		}
	}
	if !noFold {
		analysis.FoldConstArith(program)
	}
	return program
}

func emitWS(p *ws.Program) {
	fmt.Print(p.DumpWS())
}

func emitWSA(p *ws.Program) {
	fmt.Print(p.Dump("    "))
}

func emitIR(p *ws.Program) {
	fmt.Print(convertSSA(p).String())
}

func emitLLVM(p *ws.Program) {
	conf := codegen.Config{
		MaxStackLen:     maxStackLen,
		MaxCallStackLen: maxCallStackLen,
		MaxHeapBound:    maxHeapBound,
	}
	program := convertSSA(p)
	mod := codegen.EmitLLVMIR(program, conf)
	if err := llvm.VerifyModule(mod, llvm.PrintMessageAction); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Print(mod.String())
}

func printDOT(p *ws.Program) {
	fmt.Print(convertSSA(p).DotDigraph())
}

func printMatrix(p *ws.Program) {
	fmt.Print(graph.FormatMatrix(analysis.ControlFlowGraph(convertSSA(p))))
}
