# Project 7

## docs

### spec

https://drive.google.com/file/d/1CITliwTJzq19ibBF5EeuNBZ3MJ01dKoI/view

## todo

- [x] stack arithmetic commands translator
- [ ] memory access commands translator

## memo

- 스택 산술 명령 VM번역기: 9개의 arithmetic 계산과 push constant 명령여에 대한 번역을 완료하였고, SimpleAdd.vm, StackTest.vm 파일을 각각 asm 파일로 번역하는 것까지 테스트하였다. 대부분의 연산은 스택 계산 개념을 이용하여 직관적으로 구현이 가능하나, arithmetic 계산 중 lt, gt는 고려해야 할 점이 있다. 일반적으로는 스택에서 2개의 원소를 뺀 후 x-y를 D 레지스터에 담은 다음에 jump를 이용해서 구현할 것이다. 그러나 이렇게 했을 때의 문제는 x-y에서 오버 또는 언더플로우가 일어났을 때이다. 그러한 경우는 다음과 같다.
  - x < 0 && y >= 0 일 때, 오버플로우가 일어나면 x-y > 0이 될 수 있다.
  - x >= 0 && y < 0 일 때, 언더플로우가 일어나면 x-y < 0이 될 수 있다.
- 이러한 경우에 대처하기 위해, (x >= 0, x < 0) X (y >= 0, y < 0) 인 4가지 케이스에 대해 먼저 label, jump를 이용하여 분기하고, 자명하게 true, false인 경우와 x-y를 이용하여 확인해야 하는 경우를 처리한다.

- 위와 같은 경우를 원래는 translator 코드에 포함시키려고 하였으나, 다음 링크를 보니 **구현의 간단명료함**을 위해서 그런 경우를 고려하지 않는다고 말하였다.

http://nand2tetris-questions-and-answers-forum.52.s1.nabble.com/Project-7-gt-and-lt-behavior-not-clearly-specified-for-signed-operands-td4036926.html


  - 또한, sign bit로 case를 나누는 코드에서도 0 > -32768 과 같은 상황을 제대로 테스트하지 못한다. 왜냐 하면 0 - (-32768) = 32768 이어야 하지만 이 경우 overflow가 일어나 다시 -32768이 되어 제대로 된 대소 비교를 하지 못한다. 이 경우까지 분기문으로 포함시키는 것은 너무나도 큰 더티 코드가 될 것이다.
  - 이런 경우를 해결하기 가장 쉬운 방법은 ALU에서 논리 칩을 이용하여 직접 overflow 여부를 판단할 수 있게 해주는 signal을 제공하는 것일 것이다. 이를 assembly 코드에서 모두 처리하는 것은 좋지 않아 보인다. 이것은 간단한 컴퓨터를 만드려고 하는 nand2tetris 프로젝트의 특성 상 포기해야 하는 부분으로 보인다.
  - 따라서 대소 비교를 위해 sign bit로 case를 나눠 처리했던 코드는 여기에 남기고 코드의 간단명료함을 위해 더 이상 overflow, underflow를 고려하지 않도록 수정하였다.

```
@SP
AM=M-1
D=M
@GT.BGE.%v
D;JGE
@GT.BLT.%v
0;JMP
(GT.BGE.%v)
@SP
AM=M-1
M=M-D
D=D+M
@GT.BGE.AGE.%v
D;JGE
@GT.BGE.ALT.%v
0;JMP
(GT.BGE.AGE.%v)
@GT.CMPR.%v
0;JMP
(GT.BGE.ALT.%v)
@GT.FALSE.%v
0;JMP
(GT.BLT.%v)
@SP
AM=M-1
M=M-D
D=D+M
@GT.BLT.AGE.%v
D;JGE
@GT.BLT.ALT.%v
0;JMP
(GT.BLT.AGE.%v)
@GT.TRUE.%v
0;JMP
(GT.BLT.ALT.%v)
@GT.CMPR.%v
0;JMP
(GT.CMPR.%v)
@SP
A=M
D=M
@GT.TRUE.%v
D;JGT
@GT.FALSE.%v
0;JMP
(GT.TRUE.%v)
@SP
A=M
M=-1
@GT.CONT.%v
0;JMP
(GT.FALSE.%v)
@SP
A=M
M=0
@GT.CONT.%v
0;JMP
(GT.CONT.%v)
@SP
M=M+1
```
