# SSA-Test

I was bored and wanted to write a basic SSA optimiser, so I did.

Please see [the example](./example/README.md).

## The language

Programs must follow the structure:

```ssa
<assignments>

use <names>
```

Assignments can be simple or use a mathematical operator from `+-/*`:

```ssa
a = 1
a = a + 1
a = 2 / a
```

A somewhat 'complex' program might be:

```ssa
a = 1
b = a + 1
b = a * 2
c = b + 1
a = c * c

use a
```
