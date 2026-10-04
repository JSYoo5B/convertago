# Message rules

This guide defines where the rules for a valid message come from, how strongly
each rule is enforced, and which values convertago fills in on the user's
behalf. It applies to converted messages and to native models that applications
construct themselves.

## Sources of truth

Rules come from the following sources, in order of precedence:

1. **Platform reference.** The platform's official API reference and published
   policy documents define what a valid message is. A converted message must not
   fail when it is delivered, for example with an HTTP 400 response from a
   webhook.
2. **UX defaults.** Within the range the reference allows, convertago fills in
   properties the user did not specify, following the platform's UX guidance.
3. **Implementation.** Code, struct tags, and validator callbacks implement the
   two sources above. They never introduce constraints of their own.

When the implementation and the reference disagree, the implementation is wrong.
When a package README or godoc disagrees with the reference, the documentation
is wrong.

## Severities

Every rule has exactly one severity.

| Severity | Meaning | Default | `WithWarningAsError` |
| --- | --- | --- | --- |
| Fatal | The platform rejects the request. | Error | Error |
| Warning | The platform accepts the request, but the result may not be what the user intended. | Diagnostic | Error |
| Advisory | The request succeeds, but UX guidance or client limits may keep it from displaying as requested. | Diagnostic | Diagnostic |

Typical cases for each severity:

- Fatal: a missing required property, an out-of-range enum value, or a limit the
  reference states is rejected.
- Warning: the platform silently ignores a property, picks one of two
  conflicting values, or drops content beyond a limit.
- Advisory: the client truncates long text, or the request departs from a UX
  guideline without breaking anything.

Diagnostics are delivered through `WithDiagnostics` in source order. A rule
violation never changes the user's content. convertago does not truncate, drop,
or rewrite a value to satisfy a rule.

## Evidence

Each rule records the evidence for its severity: the reference URL and section,
or the date and outcome of a delivery test against the real platform.

Rules are classified as follows:

- **Request shape:** the reference defines the shape of a valid request. A
  missing required property, a value outside the allowed set, and a value whose
  type the reference does not allow in that position are **Fatal**.
- **Stated outcome:** a reference statement about the outcome verifies the rule
  at the severity that outcome implies. For example, text truncated with an
  ellipsis is **Advisory**, and a property the platform ignores is a
  **Warning**.
- **UX guidance:** recommendations from a UX guide are **Advisory**.
- **Unstated outcome:** a reference often states a limit without saying what
  happens when the limit is exceeded. Until a delivery test confirms the
  outcome, such a rule is a **Warning** and is marked as unverified. An empty
  string in a required property belongs here, because the property is still
  present in the request.

A verified rule keeps the severity its evidence supports.

## One definition per rule

Each platform package defines each rule once. Tag conversion and native
validation enforce the same definition, so the same native message receives the
same verdict and severity on both paths. A rule that one path checks and the
other does not is a defect.

go-playground validator support, through `RegisterValidation`, is an adapter
over these rule definitions. Validator tags and callbacks report only Fatal
rules, because the validator has no notion of severity.

## UX defaults

A UX default fills a property only when the user did not specify it. It never
overrides an explicit value, and it always produces a value the reference
accepts.

Defaults apply equally to converted messages and to native models that
applications construct themselves. A native model therefore needs to tell an
unspecified property apart from an explicit zero value whenever the two produce
different output. In that case, the field is a pointer or has an explicit
"unset" state.

Each package README lists its defaults with the reason for each one.

## Unavailable features

A recognized tag feature that the selected builder cannot express, such as a
style on a header without inline styles, is not a message rule. It remains an
error unless the tag is `optional`, as described in the
[tag conversion guide](tags.md#diagnostics-and-severities).
