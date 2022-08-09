# go-mongo
- `go.mongodb.org/mongo-driver` 를 기반으로 실제 프로젝트에서 유용하게 사용될 수 있는 wrapper

## 개발 배경
- 대부분의 api 서버 개발에서 통용되고 있는 MVC 구조에서 어떻게 하면 효율적으로 코드를 작성할 수 있을까? 에서 시작
- 기본적으로 go에서는 err를 throw(`panic`)를 하지 않고, 직접 에러를 리턴하고, 상위 함수에서는 해당 에러를 바탕으로 적절히 로그를 남기거나 로직을 타야한다.
- 하지만 하위 레이어에서는 본인이 내뱉을 에러 혹은 리턴 값이 어떻게 사용될 지 알 수 없고, 실제 에러인지 알 수 없다.
  - 예를 들어, db 레이어(이하 manager)에서 `userManager`가 있다고 하면, find user를 하였을 때, 해당 user가 없는 경우 이것이 발생하면 안되는 상황인지 or 정상적인 상황인지는 상위 레어이(이하 service)에서 결정된다.
  - 그렇다면 로그는 어디서 어떻게 남겨야하는거? 에 대한 고민
  - manager에서는 로그 레벨을 정할 수 없고, 이를 info, warn으로 남길 것인가?
  - 위의 예시에서 `(userDoc, error)`의 리턴 형식을 가질 때, user not found의 경우 error를 리턴해야하나? or `userDoc`을 `nil`, error도 `nil`로 해야하나?
  - manager 레이어에서 다양한 에러(not found, duplicate key, timeout, decode 등)가 발생할 수 있는데, service 레이어에서는 이를 어떻게 판별할 수 있을까?

## Basic concept
- manager 레이어에서는 로그를 남기지 않고, db 쿼리를 실행한 정보와 에러를 wrapping하여 리턴하고, service 레이어에서는 필요한 에러를 핸들링하고 로그 레벨을 판별하여 로그를 찍는다.
- `error`가 `nil`일 경우, 리턴되는 값(document 데이터)은 `nil`이 아님과 쿼리 성공을 보장한다.
- manager 레이어에서는 db error를 한번 wrapping하여, service 레이어에서는 mongo driver에서 내뱉는 에러를 핸들링하는 것이 아닌, 내부에서 정의한 에러를 핸들링힌다.
- `singleResult`, `Cursor` 등 매번 Decode하는 중복 코드를 없앤다. 