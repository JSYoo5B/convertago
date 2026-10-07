# 카카오워크

`kakaowork` 태그를 붙인 구조체를 `convertago.ToKakaoworkMessage`로 변환합니다.
이 패키지는 카카오워크 Block Kit 구조체와 해당 구조체로 변환하는 구현을 제공합니다.

태그 문법, 필드 순서, 그룹, 진단은 [태그 변환 가이드](../docs/tags.md)를 참고하세요.
규칙의 근거와 심각도 기준은 [메시지 규칙](../docs/message-rules.md)에서,
코드 생성과 교차 빌드는 [생성 가이드](../docs/generation.md)에서 설명합니다.

## 지원하는 빌더

최상위에서 사용할 수 있는 역할은 `header`, `text`, `image_link`, `button`,
`action`, `divider`, `description`, `section`, `context`, `preview`입니다.
표에서 굵게 표시한 슬롯은 태그에서 반드시 지정해야 합니다.

| 역할 | 스칼라 입력 | 중첩 요소 |
| --- | --- | --- |
| `header` | **`text`** (기본 슬롯), `style` | 없음 |
| `text` | **`text`** (기본 슬롯) | `text` 슬롯에 `styled`, `link`, `mention` |
| `styled` | **`text`** (기본 슬롯), `color`, `bold`, `italic`, `strike` | 없음 |
| `link` | **`text`** (기본 슬롯), **`url`** | 없음 |
| `mention` | **`text`** (기본 슬롯), **`user_id`** | 없음 |
| `image_link` | **`url`** (기본 슬롯) | 없음 |
| `button` | **`text`** (기본 슬롯), `style` | **`action`** 슬롯에 버튼 액션 하나 |
| `action` | 없음 | **`elements`** 슬롯에 `button` |
| `description` | **`term`**, `accent` | **`content`** 슬롯에 `text` |
| `section` | 없음 | **`content`** 슬롯에 `text`, `accessory` 슬롯에 `image_link`, `action` 슬롯에 버튼 액션 |
| `context` | 없음 | **`content`** 슬롯에 `text`, **`image`** 슬롯에 `image_link` |
| `preview` | **`text`** (기본 슬롯) | 없음 |
| `divider` | 없음 | 빈 구조체 |

`preview` 필드가 여러 개이면 선언 순서대로 구분자 없이 이어 붙여 메시지의 `text`가 됩니다.
`text`와 `styled`는 `bold`, `italic`, `strike` 스타일과 `format=plain`을 받습니다.
버튼 액션과 인라인 역할은 해당 요소를 받을 수 있는 상위 빌더 안에서만 사용합니다.

### 태그 고정값

아래 슬롯은 필드 대신 태그에 값을 적어 고정할 수 있습니다. enum 슬롯은
`이름=값`으로, bool 슬롯은 이름만 적어 `true`로 고정하며 `이름=false`도 쓸 수 있습니다.
같은 슬롯을 필드로도 넘기면 필드 값이 우선하고, 문법과 해석 순서는
[태그 변환 가이드](../docs/tags.md#fixed-slot-values)에 있습니다.

| 역할 | enum 슬롯 | bool 슬롯 |
| --- | --- | --- |
| `header` | `style` (`white`, `blue`, `red`, `yellow`) | 없음 |
| `styled` | `color` (`default`, `red`, `blue`, `grey`) | `bold`, `italic`, `strike` |
| `button` | `style` (`default`, `primary`, `danger`) | 없음 |
| `description` | 없음 | `accent` |
| `open_inapp_browser` | 없음 | `standalone` |

```go
type Notice struct {
	Title string `kakaowork:"header;style=blue"`
	Alert struct {
		Word string `kakaowork:"part"`
	} `kakaowork:"styled;color=red;bold"`
}
```

`header;style=blue`의 `style`은 값이 텍스트 서식 어휘가 아니므로 헤더 배경으로 해석합니다.
`text;style=bold`처럼 텍스트 서식 어휘만 적으면 기존과 같이 텍스트 서식입니다.

## 버튼 액션

중첩 버튼 액션은 `open_system_browser`, `open_inapp_browser`,
`open_external_app`, `submit_action`, `call_modal`, `exclusive`를 지원합니다.
앞의 다섯 액션은 `value`를 기본 입력으로 받고, 선택적으로 `name`을 받습니다.
`submit_action`은 `name`도 지정해야 합니다.

`open_inapp_browser`는 `standalone`, `width`, `height`도 받습니다.
너비와 높이는 독립 창(`standalone`)에서만 적용됩니다.
`open_external_app`의 `value`는 `ios`, `aos` 키에 앱 URI를 담은 쿼리 문자열입니다.
예를 들면 `ios=kakaomap%3A%2F%2Flook&aos=kakaomap%3A%2F%2Flook`과 같습니다.

`exclusive`는 `default` 슬롯에 기본 액션을 받고, `pc`, `mobile`, `windows`,
`macos`, `android`, `ios` 슬롯에 환경별 액션을 선택적으로 받습니다.
기본 자식 슬롯은 `default`이며, 각 슬롯에는 앞의 다섯 액션 중 하나를 넣습니다.
카카오워크는 OS, 플랫폼, `default` 순서로 적용할 액션을 고릅니다.

## 값과 기본값

헤더 배경은 `white`, `blue`, `red`, `yellow`를 받습니다.
버튼 스타일은 `default`, `primary`, `danger`를 받습니다.
인라인 색상은 `default`, `red`, `blue`, `grey`를 받습니다.

사용자가 지정하지 않은 속성은 JSON으로 직렬화할 때 아래 값으로 채웁니다.
변환 결과와 직접 만든 구조체에 똑같이 적용합니다.

| 속성 | 기본값 | 이유 |
| --- | --- | --- |
| 헤더 `style` | `white` | 레퍼런스의 기본값입니다. |
| 버튼 `style` | `default` | 레퍼런스에서 필수 속성이므로 기본 스타일로 채웁니다. |

## 규칙

변환 결과와 직접 만든 메시지에 같은 규칙을 적용합니다. 규칙을 위반해도 내용을
잘라내거나 고치지 않습니다. 심각도에 따른 처리는 다음과 같습니다.

| 심각도 | 기본 동작 | `WithWarningAsError` |
| --- | --- | --- |
| Fatal | 오류 | 오류 |
| Warning | 진단 | 오류 |
| Advisory | 진단 | 진단 |

검증 열의 "확인"은 위반했을 때의 결과를 레퍼런스가 밝히고 있다는 뜻입니다.
"미확인" 규칙은 레퍼런스가 제한만 적고 결과를 밝히지 않아 Warning으로 둡니다.

### 메시지

| 규칙 | 심각도 | 검증 | 내용 | 레퍼런스 |
| --- | --- | --- | --- | --- |
| `message.text.required` | Warning | 미확인 | 메시지 `text`(`preview`)가 필요합니다. | [메시지 API][messages] |
| `message.text.length` | Warning | 미확인 | 메시지 `text`는 10,000자 이하입니다. | [메시지 API][messages] |
| `message.blocks.type` | Fatal | 확인 | `blocks`에는 말풍선 블록만 넣습니다. | [Block Kit][blockkit] |
| `header.position` | Warning | 미확인 | 헤더 블록은 첫 블록 하나만 둘 수 있습니다. | [헤더 블록][header] |
| `message.buttons.count` | Advisory | 확인 | 버튼을 네 개 이상 사용하지 않습니다. | [UX 가이드][ux] |
| `message.image_link.after_button` | Advisory | 확인 | 버튼 블록 아래에 이미지 링크 블록을 두지 않습니다. | [UX 가이드][ux] |
| `message.image_link.consecutive` | Advisory | 확인 | 이미지 링크 블록을 연속으로 두지 않습니다. | [UX 가이드][ux] |
| `message.divider.edge` | Advisory | 확인 | 구분선 블록을 처음이나 끝에 두지 않습니다. | [UX 가이드][ux] |
| `message.divider.only` | Advisory | 확인 | 구분선 블록만으로 메시지를 구성하지 않습니다. | [UX 가이드][ux] |

### 블록

| 규칙 | 심각도 | 검증 | 내용 | 레퍼런스 |
| --- | --- | --- | --- | --- |
| `header.text.required` | Warning | 미확인 | 헤더 텍스트가 필요합니다. | [헤더 블록][header] |
| `header.text.length` | Advisory | 확인 | 20자를 넘는 헤더 텍스트는 말줄임 처리됩니다. | [헤더 블록][header] |
| `header.text.line_break` | Warning | 미확인 | 헤더 텍스트는 줄바꿈을 지원하지 않습니다. | [헤더 블록][header] |
| `header.style.value` | Fatal | 확인 | 헤더 배경은 정해진 값 중 하나입니다. | [헤더 블록][header] |
| `text.text.required` | Warning | 미확인 | 텍스트 블록의 `text`가 필요합니다. | [텍스트 블록][text] |
| `text.text.length` | Warning | 미확인 | `text`와 인라인 텍스트의 합은 각각 500자 이하입니다. | [텍스트 블록][text] |
| `text.inlines.mismatch` | Warning | 확인 | `text`가 인라인과 다르면 인라인이 우선 적용됩니다. | [텍스트 블록][text] |
| `text.inlines.required` | Fatal | 확인 | 인라인 목록에 빈 요소를 넣지 않습니다. | [텍스트 블록][text] |
| `text.links.count` | Advisory | 확인 | 텍스트 블록 하나에 링크를 두 개 이상 넣지 않습니다. | [UX 가이드][ux] |
| `styled.text.required` | Warning | 미확인 | `styled` 인라인의 텍스트가 필요합니다. | [텍스트 블록][text] |
| `styled.color.value` | Fatal | 확인 | 인라인 색상은 정해진 값 중 하나입니다. | [텍스트 블록][text] |
| `link.text.required` | Warning | 미확인 | 링크 텍스트가 필요합니다. | [텍스트 블록][text] |
| `link.url.scheme` | Warning | 미확인 | 링크 주소는 HTTP, HTTPS, mailto, tel 절대 URI입니다. | [텍스트 블록][text] |
| `mention.text.required` | Warning | 미확인 | 멘션 텍스트가 필요합니다. | [텍스트 블록][text] |
| `mention.user_id.value` | Advisory | 확인 | 멘션은 채팅방 참여자에게만 동작하므로 사용자 ID는 양수입니다. | [텍스트 블록][text] |
| `image_link.url.format` | Warning | 미확인 | 이미지 주소는 HTTP 또는 HTTPS 절대 URL입니다. | [이미지 링크 블록][image] |
| `button.text.required` | Warning | 미확인 | 버튼 텍스트가 필요합니다. | [버튼 블록][button] |
| `button.text.length` | Advisory | 확인 | 20자를 넘는 버튼 텍스트는 말줄임 처리됩니다. | [버튼 블록][button] |
| `button.style.value` | Fatal | 확인 | 버튼 스타일은 정해진 값 중 하나입니다. | [버튼 블록][button] |
| `button.action.required` | Fatal | 확인 | 버튼 액션이 필요합니다. | [버튼 블록][button] |
| `action.elements.count` | Warning | 미확인 | 액션 블록에는 버튼을 두 개 또는 세 개 넣습니다. | [액션 블록][action] |
| `description.term.required` | Warning | 미확인 | 설명 블록의 `term`이 필요합니다. | [설명 블록][description] |
| `description.term.length` | Warning | 미확인 | 설명 블록의 `term`은 10자 이하입니다. | [설명 블록][description] |

### 버튼 액션

| 규칙 | 심각도 | 검증 | 내용 | 레퍼런스 |
| --- | --- | --- | --- | --- |
| `open_system_browser.value.format` | Warning | 미확인 | 브라우저 주소는 HTTP 또는 HTTPS 절대 URL입니다. | [버튼 블록][button] |
| `open_inapp_browser.value.format` | Warning | 미확인 | 브라우저 주소는 HTTP 또는 HTTPS 절대 URL입니다. | [버튼 블록][button] |
| `open_inapp_browser.size.standalone` | Warning | 미확인 | 너비와 높이는 독립 창에서만 지정할 수 있습니다. | [버튼 블록][button] |
| `open_inapp_browser.size.value` | Warning | 미확인 | 너비와 높이는 양수 픽셀 값입니다. | [버튼 블록][button] |
| `open_external_app.value.format` | Warning | 미확인 | `value`는 `ios`, `aos` 키에 앱 URI를 담은 쿼리 문자열입니다. | [버튼 블록][button] |
| `submit_action.name.required` | Warning | 미확인 | 제출 액션의 `name`이 필요합니다. | [버튼 블록][button] |
| `submit_action.value.required` | Warning | 미확인 | 제출 액션의 `value`가 필요합니다. | [버튼 블록][button] |
| `call_modal.value.required` | Warning | 미확인 | 모달 호출 액션의 `value`가 필요합니다. | [버튼 블록][button] |
| `exclusive.default.required` | Fatal | 확인 | 환경별 액션에는 `default` 액션이 필요합니다. | [버튼 블록][button] |
| `exclusive.action.type` | Fatal | 확인 | 환경별 액션의 각 슬롯에는 `exclusive`를 넣을 수 없습니다. | [버튼 블록][button] |

이모지를 쓰지 않는 버튼 이름처럼 UX 가이드에 있지만 기계적으로 판단할 수 없는
권고는 검사하지 않습니다.

## 검증 방법

`convertago.ToKakaoworkMessage`는 메시지를 만들면서 위 규칙을 검사합니다.
진단 경로는 원본 필드 경로(`$.Items[1].Name`)입니다.

직접 만든 메시지는 `kakaowork.Validate(message, options...)`로 검사합니다.
변환과 같은 규칙과 옵션을 사용하며, 진단 경로는 전송되는 JSON의 경로(`$.blocks[1].text`)입니다.
변환한 메시지를 수정하지 않았다면 다시 검사할 필요는 없습니다.

go-playground validator를 사용한다면 `kakaowork.RegisterValidation(v)`를 등록한 뒤
`v.Struct(message)`로 검사할 수 있습니다. validator 오류에는 심각도가 없으므로 Fatal
규칙만 보고하며, 오류의 태그는 규칙 이름입니다.
사용 예시는 `ExampleValidate`와 `ExampleRegisterValidation`에 있습니다.

[blockkit]: https://docs.kakaoi.ai/kakao_work/blockkit/
[messages]: https://docs.kakaoi.ai/kakao_work/webapireference/messages/
[ux]: https://docs.kakaoi.ai/kakao_work/blockkit/uxguide/
[header]: https://docs.kakaoi.ai/kakao_work/blockkit/headerblock/
[text]: https://docs.kakaoi.ai/kakao_work/blockkit/textblock/
[image]: https://docs.kakaoi.ai/kakao_work/blockkit/imagelinkblock/
[button]: https://docs.kakaoi.ai/kakao_work/blockkit/buttonblock/
[action]: https://docs.kakaoi.ai/kakao_work/blockkit/actionblock/
[description]: https://docs.kakaoi.ai/kakao_work/blockkit/descriptionblock/
