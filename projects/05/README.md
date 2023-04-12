# Project 5

## docs

### spec

https://drive.google.com/file/d/1CJ1ymH6xdC5Z-Da8G0tqowaoOXq1cdbU/view

## todo

- [x] Memory
- [ ] CPU
- [ ] Computer

## memo

- Memory는 RAM16K, Screen, Keyboard 모듈을 하나로 합하여 만든다. 다만, RAM을 만들었었던 때처럼 address에 따라 어떤 칩으로 load signal을 보내고, 어떤 out을 선택해야 할 지 결정해야 하는데 load는 RAM16K, Screen에만 존재하므로 DMux 하나만 사용하면 되고, out은 3개가 존재하므로 Mux16을 2번 사용하면 된다.
