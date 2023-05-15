# Project 11

## docs

### spec

https://drive.google.com/file/d/1O-129lGOVNQ8XU7J4z0SGgbp7gPUv0sj/view

## todo

- [x] Symbol Table
- [x] Code Generation

## memo

- Symbol Table
  - 기존 구현에서 identifier를 출력하던 자리에서 symboltable 모듈에 저장하는 부분으로 변경만 하면 된다. 또한, 현재 위치가 어디인지에 따라서, 앞에 있는 keyword가 무엇인지에 따라서 명확하게 정의할 수 있다.
- Code Generation
  - 처음에는 굉장히 어려울 것이라고 생각하였다. 그러나 언어의 정의가 상당히 쉽고 **결정적**이기 때문에 논리적 흐름만 잘 따라가면 된다.
  - method: argument 0번으로 자기 자신 인스턴스를 넘겨주어야 한다. 따라서 symbol table에 변수들을 정의할 때 자기 자신을 첫번째 ARG에 넣는 것이 인덱스를 결정할 때 헷갈리지 않을 것이다. 다만 그 변수의 이름은 인스턴스 내부에 정의된 이름이 아니어야 하므로 "this" predefined keyword를 사용하여 겹치지 않도록 하였다. 또한, method를 호출하기에 앞서 가장 먼저 해야 할 것은 `push pointer 0` 를 통해 this에 있는 인스턴스 주소를 스택의 argument 0번 자리에 올리는 것이다.
  - array handling: 책에 나와 있기로는 `a[1]` 같은 배열 인덱싱을 할 때 배열 주소, 인덱스를 스택에 올리고 `add` 연산을 통해 더한 뒤 그 주소를 `pop pointer 1` 를 통해 `that` segment에 올리는 방식으로 사용하라고 한다. 좋은 방법이다. 그러나 문제는 "ComplexArrays" example에서 알 수 있듯이 `let` statement에서 `let a[1] = b[2]` 와 같이 왼쪽, 오른쪽 모두에 배열 접근이 등장할 때 일어난다. `a[1]` 주소를 `that` segment에 저장해 둔 뒤 `b[2]` 를 계산하다가 `that` segment의 값이 덮어씌워지기 때문이다. 이를 위해 스택에 `a`, `1` 을 더하지 않고 쌓아둔 뒤 `b[2]` 와 같은 오른쪽 expression을 모두 계산하고, 그 값을 `temp0` segment에 잠시 넣어 두었다가 `a[1]` 주소를 계산하여 `that` segment에 저장하고, 마지막으로 `push temp 0`, `pop that 0` 를 이용하여 배열에 값을 저장한다.
