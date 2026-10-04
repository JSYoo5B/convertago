# Documentation policy

This policy covers the language and required contents of package READMEs, godoc,
and examples.

## Language

A messenger package's README and godoc use the language of the platform's
origin. For example, `kakaowork` is documented in Korean, while `slack` and
`googlechat` are documented in English.

- **Other-language READMEs:** these may be added as `README.<lang>.md`, such as
  `README.en.md`. List them at the top of the package `README.md`.
- **godoc:** written only in the platform language and never translated, because
  godoc has no language variants. This covers doc comments on declarations,
  including struct fields and interface methods.
- **Comments inside function, test, and example bodies:** always English.
- **Reference lines:** `// Reference: <url>` is a shared format, used unchanged
  in every language.
- **Project-wide documents:** the root README, `docs/`, and `benchmarks/` are
  written in English.

## Package README

Each package README contains:

- The supported top-level and nested roles, with their slots, default slots,
  and required slots.
- The UX defaults, with reasons.
- The [message rules](message-rules.md), with severity, verification status,
  and reference.
- A link to the platform reference.

## godoc

- Each native model type links its reference section with a `Reference:` line.
- Field docs state constraints and defaults that match the rule definitions.
- `ToMessage` and `RegisterValidation` describe what they check, consistently
  with the README.

## Examples

Each top-level role has a tag conversion example. Each native model has a
`MarshalJSON` example. Example input must satisfy every Fatal rule.

Use `json.MarshalIndent(value, "", "  ")` and multiline `Output` comments in
JSON examples so the native message structure is visible in godoc. Keep each
executable example beside the API it documents so godoc associates it with that
declaration.
