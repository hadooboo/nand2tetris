# Project 3

## docs

### spec

https://drive.google.com/file/d/1ArUW8mkh4Kax-2TXGRpjPWuHf70u6_TJ/view

## todo

- [x] Bit
- [x] Register
- [x] RAM8
- [x] RAM64
- [x] RAM512
- [x] RAM4K
- [x] RAM16K
- [x] PC

## memo

- Bit는 load 여부를 판단하기 위한 Mux와 1 bit 정보를 저장하는 DFF로 구성한다.
- Register는 16개의 Bit로 만든다.
- RAM8은 8개의 Register와 각 Register에 load 정보를 뿌리는 DMux8Way, 각 Register로부터 out을 수합하는 Mux8Way16으로 구성한다.
- RAM64는 RAM8에서 Register 부분을 RAM8로 바꾼 것과 같다. Mux, DMux에서의 sel은 address의 상위 3개의 비트를 사용하고, 각 RAM8에서 address를 지정하기 위해서는 address의 하위 3개의 비트를 사용한다.
- RAM512, RAM4K, RAM16K는 큰 틀에서 RAM64를 구현한 방법과 동일하다. RAM512는 8개의 RAM8, RAM4K는 8개의 RAM512, RAM16K는 4개의 RAM4K를 이용하여 만든다.
- PC에서 가장 중요한 것은 load=1, inc=1, reset=1, else 경우에 대해 각각을 sel=00, 01, 10, 11으로 매핑하여 다음 input에 대한 Mux 처리를 하는 것이다. sel0이 1인 경우는 `(!reset & load) | (!reset & !load & !inc)` 이고, sel1이 1인 경우는 `(!reset & !load & inc) | (!reset & !load & !inc)` 이다. 이 논리를 Not, And, Or 게이트를 이용하여 sel input으로 만든 뒤 이에 Mux 칩을 이용하여 Register에 대한 다음 input을 결정한다.
