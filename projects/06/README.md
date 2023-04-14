# Project 6

## docs

### spec

https://drive.google.com/file/d/1CITliwTJzq19ibBF5EeuNBZ3MJ01dKoI/view

## todo

- [x] Assembler w/o symbolic references
- [x] Assembler w/ symbol handling capabilities

## memo

- 기호를 처리하지 않는 어셈블러: go로 구현하였다. 책에 있는 가이드대로 parser, code 모듈을 만들고 Add.asm, MaxL.asm, RectL.asm, PongL.asm에 대해 테스트를 완료하였다. go에서는 bit 타입은 없기 때문에 code interface에서 3bit 반환 같은 메소드는 만들지 못하고 나머지 자리가 0으로 채워진 2byte(=16bit) 배열을 반환하도록 하여 마지막에 or 처리만으로 결과값을 알 수 있게 하였다.
- 기호를 처리하는 어셈블러: 마찬가지로 go로 구현하였고, 기호를 처리하지 않는 어셈블러에 symboltable 모듈을 추가하여 구현하였다. Max.asm, Rect.asm, Pong.asm에 대해 테스트를 완료하였다. symboltable에서 label pc 주소를 찾아내기 위한 1st pass, 실제 machine code로 변환하기 위한 2nd pass를 거치도록 구현하였다.
