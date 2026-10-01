# 카카오워크

`kakaowork` 태그를 붙인 구조체를 `convertago.ToKakaoworkMessage`로 변환합니다.
이 패키지는 카카오워크 Block Kit 구조체와 해당 구조체로 변환하는 구현을 제공합니다.

태그 문법, 필드 순서, 그룹, 진단, 코드 생성 방법은
[공통 변환 가이드](../README.md)를 참고하세요.

## 지원하는 빌더

최상위에서 사용할 수 있는 역할은 `header`, `text`, `image_link`, `button`,
`action`, `divider`, `description`, `section`, `context`, `preview`입니다.

| 역할 | 스칼라 입력 | 중첩 요소 |
| --- | --- | --- |
| `header` | `text` (기본 슬롯), `style` | — |
| `text` | `text` (기본 슬롯) | `text` 슬롯에 `styled`, `link`, `mention` |
| `styled` | `text` (기본 슬롯), `color`, `bold`, `italic`, `strike` | — |
| `link` | `text` (기본 슬롯), `url` | — |
| `mention` | `text` (기본 슬롯), `user_id` | — |
| `image_link` | `url` (기본 슬롯) | — |
| `button` | `text` (기본 슬롯), `style` | `action` 슬롯에 버튼 액션 하나 |
| `action` | — | `elements` 슬롯에 `button` 두 개 또는 세 개 |
| `description` | `term` (최대 10자), `accent` | `content` 슬롯에 필수 `text` |
| `section` | — | `content` 슬롯에 필수 `text`, `accessory` 슬롯에 선택적 `image_link`, `action` 슬롯에 선택적 버튼 액션 |
| `context` | — | `content` 슬롯에 필수 `text`, `image` 슬롯에 필수 `image_link` |
| `preview` | `text` (기본 슬롯) | — |
| `divider` | — | 빈 구조체 |

## 버튼 액션

중첩 버튼 액션은 `open_system_browser`, `open_inapp_browser`,
`open_external_app`, `submit_action`, `call_modal`, `exclusive`를 지원합니다.
앞의 다섯 액션은 `value`를 기본 입력으로 받고, 선택적으로 `name`을 받습니다.
`submit_action`은 `name`이 필수입니다.

`open_inapp_browser`는 `standalone`, `width`, `height`도 받습니다.
너비나 높이를 지정하려면 독립 창을 사용하도록 `standalone`을 설정해야 합니다.

`exclusive`는 `default` 슬롯에 기본 액션을 필수로 받고,
`pc`, `mobile`, `windows`, `macos`, `android`, `ios` 슬롯에
환경별 액션을 선택적으로 받습니다. 각 슬롯에는 앞의 다섯 액션 중 하나를 넣습니다.
버튼 액션과 인라인 역할은 해당 요소를 받을 수 있는 상위 빌더 안에서만 사용합니다.

## 스타일과 URL

헤더 배경은 `white`, `blue`, `red`, `yellow`를 지원합니다.
버튼 스타일은 `default`, `primary`, `danger`를 지원합니다.
인라인 색상은 `default`, `red`, `blue`, `grey`를 지원합니다.
`text`와 `styled`는 `bold`, `italic`, `strike`를 지원합니다.

링크는 HTTP, HTTPS, mailto, tel 주소를 받습니다.
이미지와 브라우저 액션은 HTTP 또는 HTTPS 주소가 필요합니다.
외부 앱 액션은 절대 애플리케이션 URI를 받습니다.

## 유효성 검사와 제한

변환기는 헤더 위치와 블록별 텍스트 길이 등 구현된 제약과 입력 충돌을 검사합니다.
제약에 맞지 않는 내용은 잘라내거나 임의로 생략하지 않고 오류를 반환합니다.

네이티브 구조체를 직접 만들 때는 `validator.New()`로 만든 검증기에
`kakaowork.RegisterValidation(v)`를 먼저 호출하고 `v.Struct(message)`로 검사합니다.
필드 태그의 길이·URL 제한과 함께 헤더 위치, 인라인 텍스트의 합산 길이,
버튼 액션 조합을 검사합니다. `ActionBlock`은 버튼 두 개 또는 세 개가 필요하며
각 버튼의 필드도 검사합니다. 검증기 등록 예시는 `ExampleRegisterValidation`에 있습니다.

변환 오류는 원본 필드 경로를 담은 `convertago.Diagnostic`이며, 직접 검증한
오류는 네이티브 필드 경로를 담은 `validator.ValidationErrors`입니다.
변환 결과를 다시 `v.Struct`로 검사할 필요는 없습니다. 메시지를 변환한 뒤 직접
수정했다면 전송 전에 다시 검사할 수 있습니다.

지원하는 규칙은
[카카오워크 Block Kit 문서](https://docs.kakaoi.ai/kakao_work/blockkit/)를 기준으로 합니다.
