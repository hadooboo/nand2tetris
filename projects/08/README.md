# Project 8

## docs

### spec

https://drive.google.com/file/d/1F2cYb2cIPFG0B_GybMcnNUPtc5mq8mHY/view

## todo

- [x] branching commands translator
- [x] function commands translator

## memo

- 프로그램 흐름 제어 명령 VM번역기: label, goto, if-goto 문법은 assembly 문법에서와 거의 닮아 있기 때문에 구현이 간단하다. label의 경우 (label), goto의 경우 0;JMP, if-goto의 경우 stack pop 후 D;JNE 명령만 실행하면 된다.
- 함수 호출 명령 VM번역기: 직관적으로 생각해내기는 쉽지 않고 책에 있는 설명대로 작성하면 된다. 그 내용을 수도코드로 옮긴 것은 다음과 같다.

- call

```
push return-address // @return-address 명령문을 사용한 다음에 A 레지스터 값을 스택에 올린다.
push LCL // LCL, ARG, THIS, THAT 가상 레지스터에 있던 값들을 저장해둔다.
push ARG
push THIS
push THAT
ARG = SP-n-5 // ARG는 이미 스택에 푸시해놓은 값들이다. argument 세그먼트는 위에서 push한 5개의 값을 제외한 것들이기 때문에 시작 주소는 5 + (argument의 개수) 지점이다.
LCL = SP // function이 시작할 때 LCL에 local variable 개수만큼 0을 푸시하고 시작하기 때문에 call하는 시점에서는 SP와 LCL이 같다.
goto f // function 명령문 시작 주소로 이동한다.
(return-address) // function call이 끝나고 돌아올 위치를 지정한다.
```

- function

```
(f) // function이 시작할 위치를 알려준다.
repeat k times:
push 0 // 사용하는 local variable 개수만큼 0을 푸시한다. call 시점에서는 local variable을 몇 개 사용할 지 모르기 때문에 function 쪽에서 정의되는 것이 자연스럽다.
```

- return

```
FRAME = LCL // function call을 하면서 5개의 변수를 저장해두었었다. 각각의 위치를 알려주는 것은 LCL-1 ~ LCL-5 이다. 그러나 LCL 또한 복구되어야 하는 변수 중 하나이므로 값이 변해버린다. 따라서 임시로 FRAME이라는 변수를 만들어 사용한다.
RET = *(FRAME-5) // LCL-5 위치의 값이 return-address이다. argument 세그먼트의 크기가 0이었던 경우 return value를 작성하면서 return-address 자리가 덮어씌워질 수 있다. 따라서 임시변수를 선언하여 return-address 값을 저장해둔다.
*ARG = pop() // argument 세그먼트 첫번째 자리에 return value를 저장한다. caller 쪽 함수에서 스택 최상단에 있는 이 값을 return value로 생각할 것이다.
SP = ARG+1 // SP는 값이 존재하는 스택 최상단 한 칸 위를 가리키고 있어야 하기 때문에 방금 return value를 넣은 ARG 주소 + 1의 값을 가져야 한다.
THAT = *(FRAME-1) // THAT, THIS, ARG, LCL 값을 복구한다.
THIS = *(FRAME-2)
ARG = *(FRAME-3)
LCL = *(FRAME-4)
goto RET // return-address 위치로 이동한다.
```
