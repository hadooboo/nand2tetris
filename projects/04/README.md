# Project 4

## docs

### spec

https://drive.google.com/file/d/1BXyxg3biv1-v_ziFdMd8eUFj7QXFCPKI/view

## todo

- [x] Mult
- [x] Fill

## memo

- Mult는 `for (int i = M[R1]; i > 0; i--) M[R2] += M[R0]` 을 실행하는 코드이다.
- Fill의 로직은 다음과 같다.

```
INPUT 라벨: 키보드로부터 입력이 있는지 확인
    입력이 있는 경우, 이미 f=-1 이었으면 무시하고 다음 INPUT 루프 진행, f=0 이었으면 f=-1로 업데이트 후 FILL 라벨로 이동
    입력이 없는 경우 CLEAR 라벨로 이동하여, 이미 f=0 이었으면 무시하고 다음 INPUT 루프 진행, f=-1 이었으면 f=0으로 업데이트 후 FILL 라벨로 이동
FILL 라벨: i=@SCREEN으로 업데이트 후 LOOP 라벨로 이동
LOOP 라벨: i==@KBD가 되면 루프를 종료하고 다시 INPUT 라벨로 이동, 그 전까지는 i가 가리키는 주소에 f 변수에 저장되어 있는 값을 넣고 i++한 뒤 다음 루프로 진행(M[*i] = f)
```
