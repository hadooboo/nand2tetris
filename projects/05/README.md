# Project 5

## docs

### spec

https://drive.google.com/file/d/1CJ1ymH6xdC5Z-Da8G0tqowaoOXq1cdbU/view

## todo

- [x] Memory
- [x] CPU
- [x] Computer

## memo

- Memory는 RAM16K, Screen, Keyboard 모듈을 하나로 합하여 만든다. 다만, RAM을 만들었었던 때처럼 address에 따라 어떤 칩으로 load signal을 보내고, 어떤 out을 선택해야 할 지 결정해야 하는데 load는 RAM16K, Screen에만 존재하므로 DMux 하나만 사용하면 되고, out은 3개가 존재하므로 Mux16을 2번 사용하면 된다.
- CPU, Computer 에 대한 설명은 다음 손그림으로 갈음한다.

    - CPU
<img width="845" alt="Screenshot 2023-04-13 at 5 33 37 PM" src="https://user-images.githubusercontent.com/32262002/231702477-94abc9f2-876c-43e4-8981-c44b17dea4be.png">

    - Computer
<img width="536" alt="Screenshot 2023-04-13 at 5 32 55 PM" src="https://user-images.githubusercontent.com/32262002/231702281-d2d58f95-42dd-45dd-a379-e8a82c00cd1f.png">
