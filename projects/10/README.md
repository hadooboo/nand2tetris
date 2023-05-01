# Project 10

## docs

### spec

https://drive.google.com/file/d/1O1nTS24VM2kp_ilTZCrBZOryhTK1e0qN/view

## todo

- [x] JackTokenizer
- [ ] CompilationEngine

## memo

- 토큰화 모듈
  - 메모리 효율성과 최적화를 생각하기보다는 최대한 간단하게 구현하려고 하였다. 파일의 개수가 적고 라인 수도 적은 경우 효율성이 그렇게 중요한 척도는 아니기 때문이다.
  - 현재는 전체 파일 내용을 모두 다 읽어들여 string 변수로 저장한 후 계속해서 토근화하면서 남아 있는 string이 없을 때까지 반복하는 식으로 구현하였다.
  - symbol 중에서 '<', '>', '&' 는 XML encoding을 위해 각각 `&lt;`, `&gt;`, `&amp;`로 변환해서 반환한다. 이것은 XML encoding을 하는 쪽에서 변환해서 사용하는 것이 좋다고 생각하지만, 현재 프로젝트에서 XML 외에 다른 포맷으로 변환하지도 않기도 하고 깔끔한 구현을 위해 이렇게 선택하였다.
  - XML 출력을 위해 `encoding/xml` 모듈을 사용하지는 않았고 raw string을 이용하여 출력하였다. 이렇게 한 것은 현재 출력하는 tokens XML이 완전한 XML 문법을 지키는 것이 아니기 때문에 오히려 복잡해지기 때문이다. 같은 계층에서 여러 같은 키가 순서가 중요한 채로 출력되어야 하는 것이 공식 문법과 다른 점이다.
