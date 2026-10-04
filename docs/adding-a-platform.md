# Adding a platform

This guide lists what a new messenger package provides and the order in which
to build it. The package follows the [message rules](message-rules.md) and the
[documentation policy](documentation-policy.md).

## Native models

- Native models mirror the reference JSON. Field names follow the reference
  property names.
- Exported fields expose only properties the user can set. Fixed properties,
  such as a block's `type`, are emitted by `MarshalJSON`.
- Interfaces that group blocks, elements, or actions are sealed, so that only
  the package's own types implement them.

## Tests

- Each rule has a passing case and a failing case at its declared severity.
- Conversion and native validation agree on the verdict for the same message.
- Each UX default has a test, for converted output and for native output.
- `WithWarningAsError` turns Warnings into errors and leaves Advisories as diagnostics.

Place the tests as described in the
[development guide](development.md#repository-layout).

## Steps

1. Collect the reference sections the package models, and record their URLs.
2. Define the native models, following [native models](#native-models).
3. Define the rules, with severity and evidence for each.
4. Define the UX defaults, and the reasons for each.
5. Add the conversion profile and the shared rule enforcement for tag
   conversion and native validation.
6. Write the README, godoc, and examples, following the
   [documentation policy](documentation-policy.md).
7. Add the tests listed under [tests](#tests).
