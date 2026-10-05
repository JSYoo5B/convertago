# Slack

Use the `slack` tag key with `convertago.ToSlackMessage` to build
native messages. This package contains the native models and their conversion
implementation.

See the [tag conversion guide](../docs/tags.md) for grammar, ordering, groups,
and diagnostics, the [message rules](../docs/message-rules.md) for rule sources
and severities, and the [generation guide](../docs/generation.md) for generated
accessors and cross compilation.

## Supported builders

Top-level roles are `header`, `section`, `rich_text`, `image`, `actions`,
`context`, `divider`, `markdown`, `video`, `text`.

| Role | Scalar inputs | Native children |
| --- | --- | --- |
| `text` | `text` (default) | Message fallback text at the root, inline text inside rich text |
| `header` | `text` (default), `block_id`, `level` | `plain_text` in `text` |
| `section` | `text` (default), `block_id`, `expand` | One `plain_text` or `mrkdwn` object in `text`, or scalar parts; text objects in repeated `fields`; one `button` or `image` in `accessory` |
| `rich_text` | `text` (default), `block_id` | Inlines in `text`, native rich text regions in `elements` |
| `image` | `url` (default), `alt`, `title`, `block_id` | `slack_file` in `slack_file`, `plain_text` in `title` |
| `actions` | `block_id` | `button` children in `elements` |
| `context` | `block_id` | `plain_text`, `mrkdwn`, or `image` children in `elements` |
| `markdown` | `text` (default) | None |
| `video` | `video_url` (default), `title`, `alt`, `thumbnail_url`, `title_url`, `description`, `author_name`, `provider_name`, `provider_icon_url`, `block_id` | `plain_text` in `title` or `description` |
| `divider` | `block_id` | Empty struct when no ID is needed |

`plain_text` and `mrkdwn` are nested text objects. Both accept `text` by default;
`plain_text` also accepts `emoji`, and `mrkdwn` accepts `verbatim`. Scalar text in
a section or in a confirmation dialog's `text` accepts `format=plain` or
`format=mrkdwn`, with one format per text slot. Scalar `format=mrkdwn` produces a
`mrkdwn` object with `verbatim: true`, while a `mrkdwn` child keeps `verbatim`
false unless it sets the slot. Text-object children cannot mix with scalar
contributions in the same slot. The `markdown` role accepts `format=markdown`,
and the other text slots accept `format=plain`.
Nested images become ImageElement and reject `title` and `block_id`.
Root `text` fields are joined without a separator into the message fallback text.

Nested `button` accepts `text` by default, `action_id`, `url`, `value`, `style`,
`accessibility_label`, `agent_prompt`, `agent_prompt_display`, repeated
`visible_to_user_ids`, and an optional `confirm` child. Confirmation dialogs
require `title`, `text` (default), `confirm`, and `deny`, and accept `style`. Plain-text
object children can supply labels; only the dialog's `text` supports mrkdwn.
Button and confirmation styles are `primary` or `danger` when supplied.

Rich text regions are `rich_text_section`, `rich_text_list`,
`rich_text_preformatted`, and `rich_text_quote`. Lists require a `style` of
`bullet` or `ordered` and repeated sections in `elements`, and accept `indent`,
`offset`, and `border`. Sections and quotes
accept scalar parts and `text`, `link`, `user`, or `emoji` children in `text`.
Preformatted regions accept only text and links, plus `language` and `border`.
Quotes also accept `border`. A supplied border preserves an explicit zero.
Scalar parts in `rich_text` form implicit sections between explicit regions.

Inline `link` accepts `url` by default, `text`, `unsafe`, `from_llm`,
`is_slack_url`, and `truncated`. `user` accepts `user_id` by default and
`from_llm`. `emoji` accepts `name` by default and `unicode`. Rich text styles are
`bold`, `italic`, `strike`, `code`, `underline`, `highlight`, `client_highlight`,
and `unlink`. Emoji has no styles, and message fallback text carries no rich text
styles.

## Rules

Converted messages and manually constructed messages follow the same rules.
A violation never changes the message content. The severity decides what
happens, as described in the [message rules](../docs/message-rules.md).

| Severity | Default | `WithWarningAsError` |
| --- | --- | --- |
| Fatal | Error | Error |
| Warning | Diagnostic | Error |
| Advisory | Diagnostic | Diagnostic |

`chat.postMessage` rejects blocks that do not match Block Kit with
`invalid_blocks`, so the Block Kit field constraints are verified Fatal rules.
Unverified rules describe constraints whose outcome the reference does not state,
such as URL formats the reference does not define.

### Message

| Rule | Severity | Verified | Description | Reference |
| --- | --- | --- | --- | --- |
| `message.content.required` | Fatal | Yes | A message requires text or blocks. | [chat.postMessage][post] |
| `message.blocks.count` | Fatal | Yes | A message holds at most 50 blocks. | [Blocks][blocks] |
| `message.text.fallback` | Advisory | Yes | Text is recommended as the fallback for blocks. | [chat.postMessage][post] |
| `message.text.length` | Advisory | Yes | Text over 40,000 characters is truncated. | [chat.postMessage][post] |
| `markdown.text.total` | Warning | No | Markdown blocks together hold at most 12,000 characters. | [Markdown block][markdown] |
| `markdown.text.required` | Warning | No | Markdown text is required. | [Markdown block][markdown] |
| `element.required` | Fatal | Yes | Block and element lists contain no missing entries. | [Blocks][blocks] |
| `block_id.length` | Fatal | Yes | A `block_id` holds at most 255 characters. | [Blocks][blocks] |

### Text objects and layout blocks

| Rule | Severity | Verified | Description | Reference |
| --- | --- | --- | --- | --- |
| `text.text.required` | Fatal | Yes | Text object text is required. | [Text object][text] |
| `text.text.length` | Fatal | Yes | Text object text holds at most 3,000 characters. | [Text object][text] |
| `header.text.length` | Fatal | Yes | Header text holds at most 150 characters. | [Header block][header] |
| `header.level.value` | Fatal | Yes | A header level is 1 to 4. | [Header block][header] |
| `section.content.required` | Fatal | Yes | A section requires `text` or `fields`. | [Section block][section] |
| `section.fields.count` | Fatal | Yes | A section holds at most 10 fields. | [Section block][section] |
| `section.fields.length` | Fatal | Yes | Each section field holds at most 2,000 characters. | [Section block][section] |
| `image.source.required` | Fatal | Yes | An image requires exactly one of `image_url` or `slack_file`. | [Image block][image] |
| `image.image_url.length` | Fatal | Yes | `image_url` holds at most 3,000 characters. | [Image block][image] |
| `image.image_url.format` | Warning | No | `image_url` is an absolute HTTP or HTTPS URL. | [Image block][image] |
| `image.alt_text.required` | Warning | No | `alt_text` is required. | [Image block][image] |
| `image.alt_text.length` | Fatal | Yes | `alt_text` holds at most 2,000 characters. | [Image block][image] |
| `image.title.length` | Fatal | Yes | An image title holds at most 2,000 characters. | [Image block][image] |
| `slack_file.source.required` | Fatal | Yes | A Slack file requires exactly one of `url` or `id`. | [Slack file object][file] |
| `slack_file.url.format` | Warning | No | A Slack file `url` is an absolute HTTP or HTTPS URL. | [Slack file object][file] |
| `actions.elements.count` | Fatal | Yes | An actions block holds at most 25 elements. | [Actions block][actions] |
| `actions.elements.required` | Warning | No | An actions block requires elements. | [Actions block][actions] |
| `context.elements.count` | Fatal | Yes | A context block holds at most 10 elements. | [Context block][context] |
| `context.elements.required` | Warning | No | A context block requires elements. | [Context block][context] |

The image rules apply to both ImageBlock and ImageElement.

### Buttons and confirmation dialogs

| Rule | Severity | Verified | Description | Reference |
| --- | --- | --- | --- | --- |
| `button.text.length` | Fatal | Yes | Button text holds at most 75 characters. | [Button element][button] |
| `button.text.truncation` | Advisory | Yes | Button text may truncate after about 30 characters. | [Button element][button] |
| `button.action_id.length` | Fatal | Yes | `action_id` holds at most 255 characters. | [Button element][button] |
| `button.url.length` | Fatal | Yes | A button `url` holds at most 3,000 characters. | [Button element][button] |
| `button.url.format` | Warning | No | A button `url` is an absolute URI. | [Button element][button] |
| `button.value.length` | Fatal | Yes | A button `value` holds at most 2,000 characters. | [Button element][button] |
| `button.style.value` | Fatal | Yes | A button `style` is `primary` or `danger`. | [Button element][button] |
| `button.accessibility_label.length` | Fatal | Yes | `accessibility_label` holds at most 75 characters. | [Button element][button] |
| `button.agent_prompt.length` | Fatal | Yes | `agent_prompt` holds at most 4,000 characters. | [Button element][button] |
| `button.agent_prompt_display.prompt` | Warning | Yes | `agent_prompt_display` takes effect only with `agent_prompt`. | [Button element][button] |
| `button.visible_to_user_ids.required` | Warning | No | `visible_to_user_ids` contains no empty user IDs. | [Button element][button] |
| `confirm.title.length` | Fatal | Yes | A confirmation title holds at most 100 characters. | [Confirmation dialog][confirm] |
| `confirm.text.required` | Fatal | Yes | Confirmation text is required. | [Confirmation dialog][confirm] |
| `confirm.text.length` | Fatal | Yes | Confirmation text holds at most 300 characters. | [Confirmation dialog][confirm] |
| `confirm.confirm.length` | Fatal | Yes | The confirm label holds at most 30 characters. | [Confirmation dialog][confirm] |
| `confirm.deny.length` | Fatal | Yes | The deny label holds at most 30 characters. | [Confirmation dialog][confirm] |
| `confirm.style.value` | Fatal | Yes | A confirmation `style` is `primary` or `danger`. | [Confirmation dialog][confirm] |

### Video

| Rule | Severity | Verified | Description | Reference |
| --- | --- | --- | --- | --- |
| `video.title.length` | Fatal | Yes | A video title is shorter than 200 characters. | [Video block][video] |
| `video.description.length` | Fatal | Yes | A video description is shorter than 200 characters. | [Video block][video] |
| `video.author_name.length` | Fatal | Yes | `author_name` is shorter than 50 characters. | [Video block][video] |
| `video.video_url.https` | Fatal | Yes | `video_url` is an HTTPS URL. | [Video block][video] |
| `video.title_url.https` | Fatal | Yes | `title_url` is an HTTPS URL. | [Video block][video] |
| `video.alt_text.required` | Warning | No | Video `alt_text` is required. | [Video block][video] |
| `video.thumbnail_url.format` | Warning | No | `thumbnail_url` is an absolute HTTP or HTTPS URL. | [Video block][video] |
| `video.provider_icon_url.format` | Warning | No | `provider_icon_url` is an absolute HTTP or HTTPS URL. | [Video block][video] |

### Rich text

| Rule | Severity | Verified | Description | Reference |
| --- | --- | --- | --- | --- |
| `rich_text.elements.required` | Warning | No | A rich text block requires elements. | [Rich text block][rich] |
| `rich_text_section.elements.required` | Warning | No | Sections, list items, preformatted regions, and quotes require elements. | [Rich text section][richsection] |
| `rich_text_list.style.value` | Fatal | Yes | A list `style` is `bullet` or `ordered`. | [Rich text list][list] |
| `rich_text_list.offset.ordered` | Warning | No | `offset` numbers only an ordered list. | [Rich text list][list] |
| `rich_text_list.number.value` | Warning | No | `indent` and `offset` are not negative. | [Rich text list][list] |
| `rich_text.border.value` | Warning | No | `border` is 0 or 1. | [Rich text list][list] |
| `text_inline.text.required` | Warning | No | Rich text inline text is required. | [Text element][textelement] |
| `link.url.format` | Warning | No | A link `url` is an absolute URI. | [Link element][link] |
| `link.style.code` | Fatal | Yes | A link style does not accept `code`. | [Link element][link] |
| `user.user_id.required` | Warning | No | `user_id` is required. | [User element][user] |
| `user.style.code` | Fatal | Yes | A user style does not accept `code`. | [User element][user] |
| `emoji.name.required` | Warning | No | An emoji `name` is required. | [Emoji element][emoji] |

App scopes and registered video unfurl domains are not checked; the sending
application handles them.

## Validation

`convertago.ToSlackMessage` checks these rules while it builds a message.
Diagnostic paths are source field paths, such as `$.Items[1].Name`.

`slack.Validate(message, options...)` checks a manually constructed message with
the same rules and options. Its diagnostic paths follow the sent JSON, such as
`$.blocks[1].text`. A converted message needs no further validation unless you
modify it.

With go-playground validator, register `slack.RegisterValidation(v)` and call
`v.Struct(message)`. Validator errors carry no severity, so it reports only Fatal
rules, and each error's tag is the rule ID. See `ExampleValidate` and
`ExampleRegisterValidation`.

[post]: https://docs.slack.dev/reference/methods/chat.postMessage/
[blocks]: https://docs.slack.dev/reference/block-kit/blocks/
[markdown]: https://docs.slack.dev/reference/block-kit/blocks/markdown-block/
[text]: https://docs.slack.dev/reference/block-kit/composition-objects/text-object/
[header]: https://docs.slack.dev/reference/block-kit/blocks/header-block/
[section]: https://docs.slack.dev/reference/block-kit/blocks/section-block/
[image]: https://docs.slack.dev/reference/block-kit/blocks/image-block/
[file]: https://docs.slack.dev/reference/block-kit/composition-objects/slack-file-object/
[actions]: https://docs.slack.dev/reference/block-kit/blocks/actions-block/
[context]: https://docs.slack.dev/reference/block-kit/blocks/context-block/
[button]: https://docs.slack.dev/reference/block-kit/block-elements/button-element/
[confirm]: https://docs.slack.dev/reference/block-kit/composition-objects/confirmation-dialog-object/
[video]: https://docs.slack.dev/reference/block-kit/blocks/video-block/
[rich]: https://docs.slack.dev/reference/block-kit/blocks/rich-text-block/
[richsection]: https://docs.slack.dev/reference/block-kit/block-elements/rich-text-section-element/
[list]: https://docs.slack.dev/reference/block-kit/block-elements/rich-text-list-element/
[textelement]: https://docs.slack.dev/reference/block-kit/block-elements/text-element/
[link]: https://docs.slack.dev/reference/block-kit/block-elements/link-element/
[user]: https://docs.slack.dev/reference/block-kit/block-elements/user-element/
[emoji]: https://docs.slack.dev/reference/block-kit/block-elements/emoji-element/
