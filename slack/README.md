# Slack

Use the `slack` tag key with `convertago.ToSlackMessage` to build
native messages. This package contains the native models and their conversion
implementation.

See the [tag conversion guide](../docs/tags.md) for grammar, ordering, groups,
and diagnostics, and the [generation guide](../docs/generation.md) for generated
accessors and cross compilation.

## Supported builders

Top-level roles are `header`, `section`, `rich_text`, `image`, `actions`,
`context`, `divider`, `markdown`, `video`, `text`.

| Role | Scalar inputs | Native children |
| --- | --- | --- |
| `text` | `text` (default) | Message fallback text at the root, inline text inside rich text |
| `header` | `text` (default), `block_id`, `level` | `plain_text` in `text` |
| `section` | `text` (default), `block_id`, `expand` | Text objects in `text` or repeated `fields`; one `button` or `image` in `accessory` |
| `rich_text` | `text` (default), `block_id` | Inlines in `text`, native rich text regions in `elements` |
| `image` | `url` (default), `alt`, `title`, `block_id` | `slack_file` in `slack_file`, `plain_text` in `title` |
| `actions` | `block_id` | `button` children in `elements` |
| `context` | `block_id` | `plain_text`, `mrkdwn`, or `image` children in `elements` |
| `markdown` | `text` (default) | — |
| `video` | `video_url` (default), `title`, `alt`, `thumbnail_url`, `title_url`, `description`, `author_name`, `provider_name`, `provider_icon_url`, `block_id` | `plain_text` in `title` or `description` |
| `divider` | `block_id` | Empty struct when no ID is needed |

`plain_text` and `mrkdwn` are nested text objects. Both accept `text` by default;
`plain_text` also accepts `emoji`, and `mrkdwn` accepts `verbatim`. Scalar section
text accepts `format=plain` or `format=mrkdwn`, with one format per text slot.
Text-object children cannot mix with scalar contributions in the same slot.
A section requires either `text` or `fields`. An image requires `alt` and exactly
one of `url` or a `slack_file` child; that child selects exactly one of `url` or
`id`. Nested images become ImageElement and reject `title` and `block_id`.

Nested `button` accepts `text` by default, `action_id`, `url`, `value`, `style`,
`accessibility_label`, `agent_prompt`, `agent_prompt_display`, repeated
`visible_to_user_ids`, and an optional `confirm` child. Confirmation dialogs
require `title`, `text` (default), `confirm`, and `deny`, and accept `style`. Plain-text
object children can supply labels; only the dialog's `text` supports mrkdwn.
Button and confirmation styles are `primary` or `danger` when supplied.

Rich text regions are `rich_text_section`, `rich_text_list`,
`rich_text_preformatted`, and `rich_text_quote`. Lists require a `style` of
`bullet` or `ordered` and repeated sections in `elements`, and accept `indent`,
`offset`, and `border`; an offset requires an ordered list. Sections and quotes
accept scalar parts and `text`, `link`, `user`, or `emoji` children in `text`.
Preformatted regions accept only text and links, plus `language` and `border`.
Quotes also accept `border`. A supplied border preserves an explicit zero.
Scalar parts in `rich_text` form implicit sections between explicit regions.

Inline `link` accepts `url` by default, `text`, `unsafe`, `from_llm`,
`is_slack_url`, and `truncated`. `user` accepts `user_id` by default and
`from_llm`. `emoji` accepts `name` by default and `unicode`. Rich text styles are
`bold`, `italic`, `strike`, `code`, `underline`, `highlight`, `client_highlight`,
and `unlink`; links and users exclude `code`, and emoji has no styles. Message
fallback text carries no rich text styles.

## Validation and limits

The converter enforces at most 50 blocks per message and a cumulative limit of
12000 characters across Markdown blocks. Modeled native constraints and
conflicting inputs return errors. App scopes and registered video unfurl domains
must be handled by the sending application.

For manually constructed native models, call `slack.RegisterValidation(v)` on a
new go-playground validator before using `v.Struct(message)`. Field tags enforce
scalar constraints; registered callbacks enforce contextual text limits, URL
schemes, and the cumulative Markdown limit. See `ExampleRegisterValidation`.

Conversion failures are `convertago.Diagnostic` values with source field paths.
Native validation returns `validator.ValidationErrors` with native field paths.
Converted messages already enforce the modeled constraints; validate again if
you subsequently modify their native fields.

The modeled rules follow [Slack Block Kit](https://docs.slack.dev/reference/block-kit/).
