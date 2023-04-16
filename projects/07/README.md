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
