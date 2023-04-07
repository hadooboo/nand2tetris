# Project 2

## docs

### spec

https://drive.google.com/file/d/17SzlbKXl0kc5BHsKsKMrOlx-EEpWvq7g/view

### HDL guide

https://drive.google.com/file/d/1dPj4XNby9iuAs-47U9k3xtYy9hJ-ET0T/view

## todo

- [x] HalfAdder
- [x] FullAdder
- [x] Add16
- [x] Inc16
- [x] ALU

## memo

- HalfAdder는 sum을 위한 1개의 Xor 게이트, carry를 위한 1개의 And 게이트를 이용하여 만든다.
- FullAdder는 2개의 HalfAdder과 1개의 Or 게이트를 이용하여 만든다. sum은 2개의 HalfAdder를 사용하면 얻을 수 있고, carry는 각 HalfAdder 중에 한 번이라도 발생했으면 1이 되어야 하기 때문에 Or 게이트를 사용한다.
- Add16, Inc16은 16개의 FullAdder로 만든다.
- ALU는 주어진 로직에 따라 구현하면 되는데, zx, nx 등의 flag input을 받을 때는 이를 Mux 게이트의 sel input으로 지정하여 원하는 값이 output stream으로 연결되도록 구성한다. 또한, zr output bit를 계산하기 위해 16 bit 길이의 out output가 모두 0으로만 이루어져 있는지 확인해야 하는데 이는 Or8Way 게이트 2개와 Or, Not 게이트를 순차적으로 연결하여 구현할 수 있다. 최종 16 bit 길이의 out output을 이용하여 out, zr, ng 모두 계산하는데, out을 먼저 계산해 놓은 뒤 다른 bit를 계산하는 것은 internal bus를 다음 input에서 슬라이싱해서 쓸 수 없다는 점 때문에 문제가 된다. 이렇게 사용할 internal bus들을 위해서는 이가 만들어지는 output 쪽에서 여러 개의 bus를 만들어서 제공해야 한다는 것이 HDL 스펙 가장 마지막에 나와 있고, 이 방식을 적용하였다.
