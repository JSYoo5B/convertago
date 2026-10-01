# Google Chat

Use the `googlechat` tag key with `convertago.ToGoogleChatMessage` to build
native messages. This package contains the native models and their conversion
implementation.

See the [shared conversion guide](../README.md) for tag grammar, ordering,
groups, diagnostics, and generated accessors.

## Supported builders

Top-level roles are `header`, `textParagraph`, `image`, `decoratedText`,
`buttonList`, `divider`, `columns`, `grid`, `carousel`, `chipList`, `section`,
`card`, `cardWithId`, `text`, `fallbackText`.

| Role | Scalar inputs | Native children |
| --- | --- | --- |
| `header` | `text` (default), `subtitle`, `url`, `alt`, `imageType` | — |
| `textParagraph` | `text` (default), `maxLines` | — |
| `image` | `url` (default), `alt` | A click action in `onClick` |
| `decoratedText` | `text` (default), `topLabel`, `bottomLabel`, `wrapText`, `startIconVerticalAlignment` | Icons in `startIcon` or `endIcon`; paragraphs in `topLabelText`, `contentText`, or `bottomLabelText`; click action in `onClick`; `button` or `switchControl` trailing control |
| `buttonList` | — | `button` children in `buttons` |
| `columns` | — | One or two `column` children in `columnItems` |
| `grid` | `title`, `columnCount` | `gridItem` in `items`, `borderStyle`, click action in `onClick` |
| `carousel` | — | `carouselCard` children in `carouselCards` |
| `chipList` | `layout` | `chip` children in `chips` |
| `section` | `header` (default), `collapsible`, `uncollapsibleWidgetsCount` | Widgets in `widgets`, optional `collapseControl` |
| `card` | `sectionDividerStyle` | `header`, `section` in `sections`, or widgets in `widgets` |
| `cardWithId` | `cardId` | Required `card` |
| `text`, `fallbackText` | `text` (default) | Message text and card fallback text |
| `divider` | — | Empty struct |

Top-level widgets automatically form one card and one section. Explicit sections
split that card into sections; explicit card builders select separate cards.
A header must precede its card's sections or widgets. Multiple cards require
nonempty, distinct `cardId` values, using `cardWithId` wrappers. Consecutive
widgets around explicit sections preserve their order as implicit sections.
Widgets also accept `horizontalAlignment` of `START`, `CENTER`, or `END` at the
top widget level; column and carousel wrappers reject that property.

A `column` accepts `horizontalSizeStyle`, `horizontalAlignment`,
`verticalAlignment`, and widgets in `widgets`. Its allowed children are
`textParagraph`, `image`, `decoratedText`, `buttonList`, and `chipList`.
A `carouselCard` accepts `textParagraph`, `image`, and `buttonList` in `widgets`
or `footerWidgets`. A `chip` accepts `label` by default, `disabled`, `altText`,
an `icon`, and a click action in `onClick`. A `gridItem` accepts `title` by
default, `subtitle`, `id`, `layout`, and `imageComponent` in `image`.
An image component accepts `imageUri` by default, `altText`, `imageCropStyle` in
`cropStyle`, and `borderStyle`. Crop styles accept `type` and `aspectRatio`;
a custom rectangle requires a positive ratio. Borders accept `type`,
`cornerRadius`, and `color` in `strokeColor`.

Click children are `onClick`, `action`, `openLink`, or `overflowMenu`. A direct
child constructs its OnClick wrapper automatically. An explicit `onClick`
requires exactly one of those three actions. `openLink` accepts `url` by default.
An `action` accepts `function` by default, repeated `actionParameter` children in
`parameters`, `loadIndicator`, `persistValues`, `interaction`, repeated
`requiredWidgets`, and `allWidgetsAreRequired`. Parameters require `key` and
`value` with distinct keys. An overflow menu requires `overflowMenuItem` children
in `items`; each requires `text` and `onClick`, and accepts `startIcon` and
`disabled`. Nested overflow menus are rejected.

A `button` accepts `text` by default, `disabled`, `altText`, `type`, `icon`,
`color`, and required `onClick`. Buttons and chips require text or an icon.
A `color` accepts `red`, `green`, and `blue` between zero and one. An `icon`
selects exactly one of `knownIcon`, `iconUrl`, or a `materialIcon` child, and
accepts `altText` and `imageType`. Material icons accept `name` by default,
`fill`, `weight`, and `grade`. A `switchControl` accepts `name` by default,
`value`, `selected`, `controlType`, and an `action` in `onChangeAction`.
Decorated text accepts at most one of `button`, `switchControl`, and `endIcon`.
A `collapseControl` requires buttons in `expandButton` and `collapseButton`,
and accepts `horizontalAlignment`; collapse properties require `collapsible`.
Enum properties use the exact uppercase values from the native types.

Strings are literal by default. Google Chat paragraphs and decorated text escape
plain text into HTML and convert newlines into `<br>`. Explicit `html` preserves
markup; paragraphs also accept `markdown`. HTML can include escaped plain
contributions; Markdown cannot share a paragraph with plain or HTML parts.
`bold`, `italic`, `strike`, `code`, and `underline` apply to HTML text; Markdown
excludes underline. Section headers accept plain text or explicit HTML.
Image widget, header image, and icon URLs require HTTPS.

## Validation and limits

The converter enforces at most 100 widgets per card and a 32 KB limit on the
serialized cards payload. Modeled native constraints and conflicting inputs
return errors. App scopes must be handled by the sending application.

For manually constructed native models, call `googlechat.RegisterValidation(v)`
on a new go-playground validator before using `v.Struct(message)`. Field tags
enforce scalar constraints; registered callbacks enforce HTTPS image URLs,
widget totals, card identifiers, and the serialized payload limit.
See `ExampleRegisterValidation`.

Conversion failures are `convertago.Diagnostic` values with source field paths.
Native validation returns `validator.ValidationErrors` with native field paths.
Converted messages already enforce the modeled constraints; validate again if
you subsequently modify their native fields.

The modeled rules follow [Google Chat Cards v2](https://developers.google.com/workspace/chat/api/reference/rest/v1/cards).
