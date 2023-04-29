# Project 8

## docs

### spec

https://drive.google.com/file/d/1F2cYb2cIPFG0B_GybMcnNUPtc5mq8mHY/view

## todo

- [x] branching commands translator
- [ ] function commands translator

## memo

- 프로그램 흐름 제어 명령 VM번역기: label, goto, if-goto 문법은 assembly 문법에서와 거의 닮아 있기 때문에 구현이 간단하다. label의 경우 (label), goto의 경우 0;JMP, if-goto의 경우 stack pop 후 D;JNE 명령만 실행하면 된다.
