# Environment Variables

This document describes the environment variables used by `cuerator`.

| Name      | Optionality         | Description                  |
| --------- | ------------------- | ---------------------------- |
| [`DEBUG`] | defaults to `false` | enable verbose/debug logging |

⚠️ `cuerator` may consume other undocumented environment variables. This
document only shows variables declared using [Ferrite].

## Specification

All environment variables described below must meet the stated requirements.
Otherwise, `cuerator` prints usage information to `STDERR` then exits.
**Undefined** variables and **empty** values are equivalent.

The key words **MUST**, **MUST NOT**, **REQUIRED**, **SHALL**, **SHALL NOT**,
**SHOULD**, **SHOULD NOT**, **RECOMMENDED**, **MAY**, and **OPTIONAL** in this
document are to be interpreted as described in [RFC 2119].

### `DEBUG`

> enable verbose/debug logging

The `DEBUG` variable **MAY** be left undefined, in which case the default value
of `false` is used. Otherwise, the value **MUST** be either `true` or `false`.

```bash
export DEBUG=true
export DEBUG=false # (default)
```

<!-- references -->

[`debug`]: #DEBUG
[ferrite]: https://github.com/dogmatiq/ferrite
[rfc 2119]: https://www.rfc-editor.org/rfc/rfc2119.html
