#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

extern uint64_t stack_len;
extern uint64_t call_stack_len;

void printc(int64_t c) {
  fputc(c, stdout);
}

void printi(int64_t i) {
  printf("%d", (int) i);
}

int64_t readc() {
  return fgetc(stdout);
}

int64_t readi() {
  int i;
  scanf("%d", &i);
  return i;
}

void flush() {
  fflush(stdout);
}

void check_stack(uint64_t n) {
  if (stack_len < n) {
    fputs("stack underflow\n", stdout);
    fflush(stdout);
    exit(1);
  }
}

void check_call_stack() {
  if (call_stack_len < 1) {
    fputs("call stack underflow\n", stdout);
    fflush(stdout);
    exit(1);
  }
}