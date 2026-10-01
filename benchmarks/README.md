# 리플렉션 캐시와 생성 접근자 벤치마크

Apple M3 Pro, darwin/arm64, Go 1.27.1, GOMAXPROCS=1. 2026-10-01에 커밋 `73375ce`를 측정했다. race 옵션 없이 각 항목을 3회 실행한 중앙값이다.

CPU 시간은 `os.wait4`가 반환한 user·system 시간의 합이다. 프로세스 시작, 벤치마크 준비, GC 비용을 포함하며, 경과 시간인 `ns/op`와 별도로 측정했다. 생성과 컴파일은 측정 전에 완료했다.

## 생성 접근자 비교의 범위

이 생성기는 AST와 타입 정보를 분석해 `ConvertagoFields` 접근자를 생성한다. 생성된 코드는 알려진 필드를 직접 읽고 태그 메타데이터를 구성한다. 입력 확인과 접근자 선택은 두 경로가 같은 엔트리포인트를 거치며, 메시지 조합·검증·네이티브 모델 생성도 공통 구현을 사용한다. 아래 Generated 항목은 이러한 생성 접근자의 런타임 비용을 측정한다. AST 분석·코드 생성·Go 컴파일에 걸리는 시간은 포함하지 않는다.

Reflection은 같은 필드 배치와 태그를 가진 별도의 정의 타입으로 접근자 메서드를 제거해 리플렉션을 선택한다. 두 입력은 같은 값을 가리킨다. Flat은 헤더·그룹 텍스트·정수, Nested는 이미지·목록·포인터를 포함한다. Dynamic은 동일한 Nested 값을 `any` 필드에 담은 입력으로, 생성된 코드에서도 해당 내용은 캐시 리플렉션으로 읽는다. 각 경로의 선택, 소스 필드와 최종 JSON의 일치를 테스트했다.

정적 입력인 Flat·Nested의 필드 읽기 CPU 비용은 생성 접근자에서 약 24~29% 줄었다. 중첩 입력의 전체 변환 CPU 감소는 카카오워크 5.3%, Slack 4.8%, Google Chat 3.2%였다. 소스 읽기 이후의 작업을 공유하므로 전체 변환에서의 절감 폭은 작아졌다.

생성 접근자의 할당 바이트 감소는 Flat 48 B/op, Nested 40 B/op로 작았다. Nested의 할당 횟수는 1회 늘었다. 생성 코드에는 호출 시 구성되는 스타일 슬라이스와 소스 값이 남아 있다. Dynamic의 전체 변환 CPU 차이는 약 0.4~1.2%로 반복 측정의 변동과 비슷한 수준이므로 이 입력의 큰 개선을 뒷받침하지 않는다.

## 필드 읽기: 캐시 리플렉션과 생성 접근자

각 항목을 독립 프로세스에서 500,000회 실행했다. `BenchmarkSourceFields`의 Reflection과 Generated를 같은 `conversion.Fields` 엔트리포인트로 비교한다. 첫 호출에서 필요한 캐시를 준비하며, 현재 값을 읽는 비용을 포함하고 메시지 조합과 네이티브 렌더링은 제외한다.

| 메신저 | 입력 | 경로 | ns/op | CPU 초 / 50만 회 | B/op | allocs/op | 최대 RSS MiB |
| --- | --- | --- | ---: | ---: | ---: | ---: | ---: |
| 카카오워크 | Flat | Reflection | 353.0 | 0.1771 | 1,091 | 3 | 11.52 |
| 카카오워크 | Flat | Generated | 253.6 | 0.1270 | 1,043 | 3 | 11.62 |
| 카카오워크 | Nested | Reflection | 1,243.0 | 0.6134 | 3,432 | 10 | 11.73 |
| 카카오워크 | Nested | Generated | 932.4 | 0.4584 | 3,392 | 11 | 11.64 |
| 카카오워크 | Dynamic | Reflection | 1,423.0 | 0.7021 | 3,608 | 11 | 11.91 |
| 카카오워크 | Dynamic | Generated | 1,361.0 | 0.6716 | 3,592 | 10 | 11.66 |
| Slack | Flat | Reflection | 346.2 | 0.1734 | 1,091 | 3 | 11.64 |
| Slack | Flat | Generated | 254.7 | 0.1282 | 1,043 | 3 | 11.56 |
| Slack | Nested | Reflection | 1,286.0 | 0.6359 | 3,656 | 10 | 11.95 |
| Slack | Nested | Generated | 956.4 | 0.4694 | 3,616 | 11 | 11.81 |
| Slack | Dynamic | Reflection | 1,527.0 | 0.7542 | 3,832 | 11 | 11.70 |
| Slack | Dynamic | Generated | 1,432.0 | 0.7063 | 3,816 | 10 | 12.08 |
| Google Chat | Flat | Reflection | 348.9 | 0.1753 | 1,091 | 3 | 11.88 |
| Google Chat | Flat | Generated | 248.2 | 0.1241 | 1,043 | 3 | 11.58 |
| Google Chat | Nested | Reflection | 1,289.0 | 0.6358 | 3,656 | 10 | 11.67 |
| Google Chat | Nested | Generated | 980.4 | 0.4806 | 3,616 | 11 | 11.94 |
| Google Chat | Dynamic | Reflection | 1,483.0 | 0.7326 | 3,832 | 11 | 11.80 |
| Google Chat | Dynamic | Generated | 1,409.0 | 0.6949 | 3,816 | 10 | 11.80 |

## 메신저별 전체 변환

각 항목을 독립 프로세스에서 100,000회 실행했다. 두 경로에 같은 입력을 사용하고, 변환 전 캐시를 준비한다. 값 읽기·조합·검증·네이티브 메시지 생성 비용을 포함한다. 호출자의 JSON 직렬화는 제외하며, 변환기 내부의 직렬화는 포함한다.

| 메신저 | 입력 | 경로 | ns/op | CPU 초 / 10만 회 | B/op | allocs/op | 최대 RSS MiB |
| --- | --- | --- | ---: | ---: | ---: | ---: | ---: |
| 카카오워크 | Flat | Reflection | 2,499.0 | 0.2515 | 3,488 | 46 | 12.08 |
| 카카오워크 | Flat | Generated | 2,399.0 | 0.2409 | 3,440 | 46 | 12.22 |
| 카카오워크 | Nested | Reflection | 7,777.0 | 0.7733 | 11,056 | 123 | 12.41 |
| 카카오워크 | Nested | Generated | 7,368.0 | 0.7325 | 11,016 | 124 | 12.66 |
| 카카오워크 | Dynamic | Reflection | 10,612.0 | 1.0562 | 12,456 | 143 | 12.89 |
| 카카오워크 | Dynamic | Generated | 10,519.0 | 1.0467 | 12,440 | 142 | 12.78 |
| Slack | Flat | Reflection | 2,879.0 | 0.2894 | 3,784 | 56 | 12.34 |
| Slack | Flat | Generated | 2,747.0 | 0.2755 | 3,736 | 56 | 12.41 |
| Slack | Nested | Reflection | 9,962.0 | 0.9906 | 12,808 | 153 | 12.66 |
| Slack | Nested | Generated | 9,484.0 | 0.9427 | 12,768 | 154 | 12.45 |
| Slack | Dynamic | Reflection | 13,086.0 | 1.3019 | 14,272 | 175 | 13.08 |
| Slack | Dynamic | Generated | 13,030.0 | 1.2968 | 14,256 | 174 | 13.06 |
| Google Chat | Flat | Reflection | 5,546.0 | 0.5543 | 5,960 | 75 | 13.00 |
| Google Chat | Flat | Generated | 5,445.0 | 0.5437 | 5,912 | 75 | 12.78 |
| Google Chat | Nested | Reflection | 16,555.0 | 1.6474 | 17,890 | 193 | 13.34 |
| Google Chat | Nested | Generated | 16,034.0 | 1.5946 | 17,850 | 194 | 13.19 |
| Google Chat | Dynamic | Reflection | 20,058.0 | 1.9961 | 19,355 | 215 | 13.44 |
| Google Chat | Dynamic | Generated | 19,812.0 | 1.9711 | 19,339 | 214 | 13.58 |

## 타입 메타데이터 캐시

각 항목을 독립 프로세스에서 2,000,000회 실행했다. 내부 공통 테스트 프로파일의 입력을 사용한다. Cached는 이미 만들어진 타입 계획을 조회하고, Uncached는 매번 태그를 해석하고 구조를 검증한다. Uncached에는 캐시 삽입·삭제가 포함되지 않는다. 값 읽기와 메시지 생성은 제외한다.

| 입력 | 캐시 | ns/op | CPU 초 / 200만 회 | B/op | allocs/op | 최대 RSS MiB |
| --- | --- | ---: | ---: | ---: | ---: | ---: |
| Flat | Cached | 25.1 | 0.0551 | 0 | 0 | 6.56 |
| Flat | Uncached | 2,482.0 | 4.8672 | 5,224 | 27 | 11.39 |
| Nested | Cached | 25.3 | 0.0536 | 0 | 0 | 6.62 |
| Nested | Uncached | 3,267.0 | 6.4695 | 5,008 | 44 | 11.41 |

## 현재 값 읽기: Cold와 Warm

내부 공통 테스트 프로파일의 입력을 100ms씩 3회 측정했다. Cold는 해당 입력의 root 및 dynamic 타입 계획을 매번 지운 뒤 `Fields`를 호출한다. 삭제 비용은 타이머와 할당량 측정에서 제외한다. 타이머 중단·재개 비용이 프로세스 CPU 시간을 왜곡하므로 이 표에서는 Go 벤치마크가 측정한 경과 시간과 할당량을 비교한다. 위 메신저별 SourceFields와는 입력 구조가 다르다.

| 입력 | 캐시 | ns/op | B/op | allocs/op |
| --- | --- | ---: | ---: | ---: |
| Flat | Cold | 3,834.0 | 6,552 | 33 |
| Flat | Warm | 373.6 | 1,219 | 3 |
| Nested | Cold | 5,828.0 | 9,360 | 62 |
| Nested | Warm | 1,374.0 | 4,225 | 15 |
| Dynamic | Cold | 5,568.0 | 10,208 | 67 |
| Dynamic | Warm | 1,565.0 | 4,417 | 16 |

## 메모리와 재현 방법

`B/op`와 `allocs/op`는 호출당 할당량이다. 최대 RSS는 Go 런타임과 테스트 초기화를 포함한 프로세스 전체의 최고값이다. 이 값들은 캐시가 보관하는 메모리 크기와 별개다. 캐시 적중과 생성 접근자 경로에서도 현재 값 읽기와 메시지 생성에 필요한 할당은 발생한다.

138개 원본 측정값은 [measurements.csv](measurements.csv)에 기록했다. Fields 항목은 타이머 중단·재개 때문에 CPU와 RSS 열을 비워 두었다. [측정 도구](measure.py)는 macOS와 Linux에서 바이너리를 먼저 빌드하고 두 경로를 같은 횟수로 실행한다. CSV는 모든 측정이 성공한 뒤 갱신한다.

```sh
go generate ./internal/benchmarksource
python3 benchmarks/measure.py --output benchmarks/measurements.csv
```

일반 Go 벤치마크 명령은 [루트 README](../README.md#benchmarks)를 참고한다. 전체 테스트, race를 적용한 벤치마크 1회 실행, go vet, git diff --check를 통과했다. 성능 수치는 이 환경의 관측값이며 테스트 통과 조건으로 사용하지 않는다.
