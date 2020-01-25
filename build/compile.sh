#!/bin/bash

program="$1"
out="$2"

../nebula "$program" llvm > "$out.ll"
clang -S -emit-llvm ../codegen/ext/ext.c
llvm-link -o "$out.o" "$out.ll" ext.ll
llc "$out.o"
clang -o "$out" "$out.o.s"
