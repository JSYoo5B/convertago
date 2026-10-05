# Google Chat

Use the `googlechat` tag key with `convertago.ToGoogleChatMessage` to build
native messages. This package contains the native models and their conversion
implementation.

See the [tag conversion guide](../docs/tags.md) for grammar, ordering, groups,
and diagnostics, the [message rules](../docs/message-rules.md) for rule sources
and severities, and the [generation guide](../docs/generation.md) for generated
accessors and cross compilation.

## Supported builders

Top-level roles are `header`, `textParagraph`, `image`, `decoratedText`,
`buttonList`, `divider`, `columns`, `grid`, `carousel`, `chipList`, `section`,
`card`, `cardWithId`, `text`, `fallbackText`.

| Role | Scalar inputs | Native children |
| --- | --- | --- |
| `header` | `text` (default), `subtitle`, `url`, `alt`, `imageType` | None |
| `textParagraph` | `text` (default), `maxLines` | None |
| `image` | `url` (default), `alt` | A click action in `onClick` |
| `decoratedText` | `text` (default), `topLabel`, `bottomLabel`, `wrapText`, `startIconVerticalAlignment` | Icons in `startIcon` or `endIcon`; paragraphs in `topLabelText`, `contentText`, or `bottomLabelText`; click action in `onClick`; `button` or `switchControl` trailing control |
| `buttonList` | None | `button` children in `buttons` |
| `columns` | None | One or two `column` children in `columnItems` |
| `grid` | `title`, `columnCount` | `gridItem` in `items`, `borderStyle`, click action in `onClick` |
| `carousel` | None | `carouselCard` children in `carouselCards` |
| `chipList` | `layout` | `chip` children in `chips` |
| `section` | `header` (default), `collapsible`, `uncollapsibleWidgetsCount` | Widgets in `widgets`, optional `collapseControl` |
| `card` | `sectionDividerStyle` | `header`, `section` in `sections`, or widgets in `widgets` |
| `cardWithId` | `cardId` | Required `card` |
| `text`, `fallbackText` | `text` (default) | Message text and card fallback text |
| `divider` | None | Empty struct, plus `horizontalAlignment` at the top level |

Contiguous top-level widgets form an implicit card with one section. Explicit
sections split that card into sections; explicit card builders select separate
cards. A header must precede its card's sections or widgets and occur once.
Widgets on both sides of an explicit card form two implicit cards without a
`cardId`, which `message.card_id.required` rejects; use `cardWithId` wrappers for
multiple cards. Consecutive widgets around explicit sections preserve their order
as implicit sections. Top-level widgets also accept `horizontalAlignment` of
`START`, `CENTER`, or `END`; widgets nested in a column or carousel card reject it.

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
`value`. An overflow menu requires `overflowMenuItem` children
in `items`; each requires `text` and `onClick`, and accepts `startIcon` and
`disabled`. A nested overflow menu is dropped by Google Chat.

A `button` accepts `text` by default, `disabled`, `altText`, `type`, `icon`,
`color`, and required `onClick`. Buttons and chips require text or an icon.
A `color` accepts `red`, `green`, and `blue` between zero and one. An `icon`
selects exactly one of `knownIcon`, `iconUrl`, or a `materialIcon` child, and
accepts `altText` and `imageType`. Material icons accept `name` by default,
`fill`, `weight`, and `grade`. A `switchControl` accepts `name` by default,
`value`, `selected`, `controlType`, and an `action` in `onChangeAction`.
Decorated text accepts at most one of `button`, `switchControl`, and `endIcon`.
A `collapseControl` requires buttons in `expandButton` and `collapseButton`,
and accepts `horizontalAlignment`; collapse properties apply only with `collapsible`.
Enum properties use the exact uppercase values from the native types.

Strings are literal by default. Google Chat paragraphs and decorated text escape
plain text into HTML and convert newlines into `<br>`. Explicit `html` preserves
markup; paragraphs also accept `markdown`. HTML can include escaped plain
contributions; Markdown cannot share a paragraph with plain or HTML parts.
`bold`, `italic`, `strike`, `code`, and `underline` apply to HTML text; Markdown
excludes underline. Section headers accept plain text or explicit HTML.
Image widget, header image, and icon URLs use HTTPS.

## Rules

Converted messages and manually constructed messages follow the same rules.
A violation never changes the message content. The severity decides what
happens, as described in the [message rules](../docs/message-rules.md).

| Severity | Default | `WithWarningAsError` |
| --- | --- | --- |
| Fatal | Error | Error |
| Warning | Diagnostic | Error |

Fatal rules cover the request shape: required fields, union fields that accept
exactly one value, and enum values. Verified Warnings describe outcomes the
reference states, such as ignored sections or fields. Unverified Warnings
describe constraints whose outcome the reference does not state. All links point
to the [Cards v2 reference][cards] unless noted.

### Message, card, and section

| Rule | Severity | Verified | Description | Reference |
| --- | --- | --- | --- | --- |
| `message.cards.size` | Warning | No | The cards hold at most 32 KB of JSON. | [Message][message] |
| `message.card_id.required` | Fatal | Yes | Each card has a `cardId` when a message has several cards. | [CardWithId][cardwithid] |
| `message.card_id.unique` | Warning | No | `cardId` values are unique. | [CardWithId][cardwithid] |
| `card.content.required` | Fatal | Yes | A card requires a header or sections. | [Card][cards-card] |
| `card.widgets.count` | Warning | Yes | Sections that push a card over 100 widgets, counting nested widgets, are ignored. | [Card][cards-card] |
| `card.section_divider_style.value` | Fatal | Yes | `sectionDividerStyle` is a known value. | [Card][cards-card] |
| `header.title.required` | Warning | No | A header title is required. | [CardHeader][cards-header] |
| `header.image_url.https` | Warning | No | A header `imageUrl` is an HTTPS URL. | [CardHeader][cards-header] |
| `header.image.url_required` | Warning | No | `imageType` and `imageAltText` apply only with `imageUrl`. | [CardHeader][cards-header] |
| `image_type.value` | Fatal | Yes | `imageType` is `SQUARE` or `CIRCLE`. | [ImageType][cards-imagetype] |
| `section.widgets.required` | Warning | No | A section requires widgets. | [Section][cards-section] |
| `section.collapse.collapsible` | Warning | Yes | Collapse properties apply only to a collapsible section. | [Section][cards-section] |
| `section.uncollapsible.count` | Warning | No | `uncollapsibleWidgetsCount` is between 0 and the section's widget count. | [Section][cards-section] |
| `widget.content.required` | Fatal | Yes | A widget wrapper requires content. | [Widget][cards-widget] |
| `horizontal_alignment.value` | Fatal | Yes | `horizontalAlignment` is a known value. | [HorizontalAlignment][cards-halign] |
| `vertical_alignment.value` | Fatal | Yes | `startIconVerticalAlignment` is a known value. | [VerticalAlignment][cards-valign] |

### Widgets

| Rule | Severity | Verified | Description | Reference |
| --- | --- | --- | --- | --- |
| `text_paragraph.max_lines.value` | Warning | No | `maxLines` is not negative. | [TextParagraph][cards-paragraph] |
| `text_paragraph.text_syntax.value` | Fatal | Yes | `textSyntax` is `HTML` or `MARKDOWN`. | [TextParagraph][cards-paragraph] |
| `image.image_url.https` | Warning | No | An image `imageUrl` is an HTTPS URL. | [Image][cards-image] |
| `decorated_text.text.required` | Warning | No | Decorated text requires `text`. | [DecoratedText][cards-decorated] |
| `decorated_text.control.count` | Fatal | Yes | Decorated text holds at most one of `button`, `switchControl`, or `endIcon`. | [DecoratedText][cards-decorated] |
| `switch_control.name.required` | Warning | No | A switch control `name` is required. | [SwitchControl][cards-switch] |
| `switch_control.control_type.value` | Fatal | Yes | `controlType` is `SWITCH` or `CHECK_BOX`. | [SwitchControl][cards-switch] |
| `button_list.buttons.required` | Warning | No | A button list requires buttons. | [ButtonList][cards-buttonlist] |
| `button.content.required` | Warning | No | A button requires text or an icon. | [Button][cards-button] |
| `button.type.value` | Fatal | Yes | A button `type` is a known value. | [Button][cards-button] |
| `button.type.color` | Warning | Yes | A `color` forces the `FILLED` type, so another `type` is ignored. | [Button][cards-button] |
| `color.component.range` | Warning | No | Color components are between 0 and 1. | [Color][cards-color] |
| `icon.source.count` | Fatal | Yes | An icon has exactly one of `knownIcon`, `iconUrl`, or `materialIcon`. | [Icon][cards-icon] |
| `icon.icon_url.https` | Warning | No | `iconUrl` is an HTTPS URL. | [Icon][cards-icon] |
| `material_icon.name.required` | Warning | No | A Material icon `name` is required. | [MaterialIcon][cards-material] |
| `material_icon.weight.value` | Warning | Yes | An unsupported `weight` uses the default. | [MaterialIcon][cards-material] |
| `material_icon.grade.value` | Warning | Yes | An unsupported `grade` uses the default. | [MaterialIcon][cards-material] |

### Click actions

| Rule | Severity | Verified | Description | Reference |
| --- | --- | --- | --- | --- |
| `on_click.action.count` | Fatal | Yes | `onClick` has exactly one of `action`, `openLink`, or `overflowMenu`. | [OnClick][cards-onclick] |
| `action.function.required` | Warning | No | An action `function` is required. | [Action][cards-action] |
| `action.load_indicator.value` | Fatal | Yes | `loadIndicator` is `SPINNER` or `NONE`. | [Action][cards-action] |
| `action.interaction.value` | Fatal | Yes | `interaction` is `OPEN_DIALOG`. | [Action][cards-action] |
| `action.required_widgets.conflict` | Warning | No | `requiredWidgets` and `allWidgetsAreRequired` are not combined. | [Action][cards-action] |
| `action.required_widgets.required` | Warning | No | `requiredWidgets` contains no empty names. | [Action][cards-action] |
| `action.parameters.key_required` | Warning | No | An action parameter `key` is required. | [ActionParameter][cards-parameter] |
| `action.parameters.key_unique` | Warning | No | Action parameter keys are unique. | [ActionParameter][cards-parameter] |
| `open_link.url.format` | Warning | No | An `openLink` URL is an absolute URI. | [OpenLink][cards-openlink] |
| `overflow_menu.items.required` | Warning | No | An overflow menu requires items. | [OverflowMenu][cards-overflow] |
| `overflow_menu_item.text.required` | Warning | No | An overflow menu item requires `text`. | [OverflowMenuItem][cards-overflowitem] |
| `overflow_menu_item.on_click.overflow` | Warning | Yes | A nested overflow menu is dropped and disables the item. | [OverflowMenuItem][cards-overflowitem] |

### Layouts

| Rule | Severity | Verified | Description | Reference |
| --- | --- | --- | --- | --- |
| `columns.column_items.count` | Warning | No | Columns display one or two columns. | [Columns][cards-columns] |
| `column.widgets.required` | Warning | No | A column requires widgets. | [Column][cards-column] |
| `column.horizontal_size_style.value` | Fatal | Yes | `horizontalSizeStyle` is a known value. | [Column][cards-column] |
| `column.vertical_alignment.value` | Fatal | Yes | A column `verticalAlignment` is a known value. | [Column][cards-column] |
| `grid.items.required` | Warning | No | A grid requires items. | [Grid][cards-grid] |
| `grid.column_count.value` | Warning | No | `columnCount` is not negative. | [Grid][cards-grid] |
| `grid_item.content.required` | Warning | No | A grid item requires a title, subtitle, or image. | [GridItem][cards-griditem] |
| `grid_item.layout.value` | Fatal | Yes | A grid item `layout` is a known value. | [GridItem][cards-griditem] |
| `image_component.image_uri.format` | Warning | No | `imageUri` is an absolute HTTP or HTTPS URL. | [ImageComponent][cards-imagecomponent] |
| `image_crop.type.value` | Fatal | Yes | A crop `type` is a known value. | [ImageCropStyle][cards-crop] |
| `image_crop.aspect_ratio.value` | Warning | No | `RECTANGLE_CUSTOM` requires a positive, finite `aspectRatio`. | [ImageCropStyle][cards-crop] |
| `image_crop.aspect_ratio.custom` | Warning | No | `aspectRatio` applies only to `RECTANGLE_CUSTOM`. | [ImageCropStyle][cards-crop] |
| `border_style.type.value` | Fatal | Yes | A border `type` is `NO_BORDER` or `STROKE`. | [BorderStyle][cards-border] |
| `border_style.stroke_color.stroke` | Warning | No | `strokeColor` applies only to a `STROKE` border. | [BorderStyle][cards-border] |
| `border_style.corner_radius.value` | Warning | No | `cornerRadius` is not negative. | [BorderStyle][cards-border] |
| `carousel.cards.required` | Warning | No | A carousel requires cards. | [Carousel][cards-carousel] |
| `carousel_card.widgets.required` | Warning | No | A carousel card requires widgets. | [CarouselCard][cards-carouselcard] |
| `chip_list.chips.required` | Warning | No | A chip list requires chips. | [ChipList][cards-chiplist] |
| `chip_list.layout.value` | Fatal | Yes | A chip list `layout` is a known value. | [ChipList][cards-chiplist] |
| `chip.content.required` | Warning | No | A chip requires a label or an icon. | [Chip][cards-chip] |

App scopes are not checked; the sending application handles them.

## Validation

`convertago.ToGoogleChatMessage` checks these rules while it builds a message.
Diagnostic paths are source field paths, such as `$.Items[1].Name`.

`googlechat.Validate(message, options...)` checks a manually constructed message
with the same rules and options. Its diagnostic paths follow the sent JSON, such
as `$.cardsV2[0].card.sections[0].widgets[1].image.imageUrl`. A converted message
needs no further validation unless you modify it.

With go-playground validator, register `googlechat.RegisterValidation(v)` and
call `v.Struct(message)`. Validator errors carry no severity, so it reports only
Fatal rules, and each error's tag is the rule ID. See `ExampleValidate` and
`ExampleRegisterValidation`.

[cards]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards
[message]: https://developers.google.com/workspace/chat/api/reference/rest/v1/spaces.messages#Message
[cardwithid]: https://developers.google.com/workspace/chat/api/reference/rest/v1/spaces.messages#CardWithId
[cards-card]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Card
[cards-header]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#CardHeader
[cards-imagetype]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#ImageType
[cards-section]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Section
[cards-widget]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Widget
[cards-halign]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#HorizontalAlignment
[cards-valign]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#VerticalAlignment
[cards-paragraph]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#TextParagraph
[cards-image]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Image
[cards-decorated]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#DecoratedText
[cards-switch]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#SwitchControl
[cards-buttonlist]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#ButtonList
[cards-button]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Button
[cards-color]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Color
[cards-icon]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Icon
[cards-material]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#MaterialIcon
[cards-onclick]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#OnClick
[cards-action]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Action
[cards-parameter]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#ActionParameter
[cards-openlink]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#OpenLink
[cards-overflow]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#OverflowMenu
[cards-overflowitem]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#OverflowMenuItem
[cards-columns]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Columns
[cards-column]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Column
[cards-grid]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Grid
[cards-griditem]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#GridItem
[cards-imagecomponent]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#ImageComponent
[cards-crop]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#ImageCropStyle
[cards-border]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#BorderStyle
[cards-carousel]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Carousel
[cards-carouselcard]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#CarouselCard
[cards-chiplist]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#ChipList
[cards-chip]: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Chip
