# Changelog

All notable changes to this project are documented here. The project follows
[Semantic Versioning](https://semver.org/). Before v1.0.0, minor releases may
change the public API.

## v0.1.0

The first release.

### Added

- Conversion of tagged Go structs into Kakao Work Block Kit, Slack Block Kit,
  and Google Chat Cards v2 messages through `ToKakaoworkMessage`,
  `ToSlackMessage`, and `ToGoogleChatMessage`, with a separate tag namespace for
  each messenger. See the [tag conversion guide](docs/tags.md).
- Native message models for each messenger, usable without tags.
- Message rules taken from each platform's reference, each with a Fatal,
  Warning, or Advisory severity, its reference, and its verification status.
  Conversion and the new `Validate` function of each package apply the same
  rules. `WithWarningAsError` turns Warnings into errors, and `WithDiagnostics`
  receives the rest. See the [message rules](docs/message-rules.md).
- `RegisterValidation` adapters for go-playground validator, which report the
  Fatal rules.
- Fixed tag values for enum and boolean slots, such as
  `kakaowork:"styled;color=red;bold"`, with fields able to override them.
- A generator, `cmd/convertago`, that writes direct field accessors for known
  struct types. See the [generation guide](docs/generation.md).

### Requirements

- Go 1.26 or later.

### Known limitations

- Many limits are unverified Warnings, because the references do not state what
  happens when they are exceeded and no delivery tests have been recorded yet.
- A Kakao Work message without `preview` text reports `message.text.required`;
  no UX default fills it.
- Rule checks add CPU time and allocations to conversion, as recorded in the
  [benchmark results](benchmarks/README.md).
- Generated accessors keep the tags they were generated from. Regenerate them
  after changing tags or upgrading convertago.
- `golang.org/x/crypto`, required through go-playground validator, has a known
  issue without a fix (GO-2026-5932) in a package that convertago does not use.
