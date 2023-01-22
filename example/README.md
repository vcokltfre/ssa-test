# Example

In the given example:

```ssa
a = 1
a = 2
a = 3
b = a + a

use a b
```

We assign pointlessly to 'a' twice, since those values never get used. If we run `go run main.go optimise examples/example.ssa examples/optimised.ssa` we get:

```ssa
a = 3
b = a + a
use a b
```

Which has removed the pointless assignments.

Setting `SSA_DUMP=1` will also give the SSA representation of the original code, which was used to perform the optimisation:

```ssa
a-0 = 1
a-1 = 2
a-2 = 3
b-0 = a-2 + a-2
```
