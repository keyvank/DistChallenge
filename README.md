# Fly.io Distributed Challenges

<https://fly.io/dist-sys>

## Unique ID

Suggested solution:

Use combination of timestamp(milliseconds precision) + node id + counter.

In this method if two nodes try to create an ID at the same time the
node id will be different so the ID will be unique.

If a node creates more than one ID in a milliseconds the counter
will be different, because we increament the counter on every ID
generation.

<details>
  <summary>Learning about time</summary>

Using `time.Now().String()` as Id would also work and the tests run fine.
But why?

This is the output if we print both `time.Now.UnixNano()` and `time.Now().String()`.

```
2023-06-02 18:05:01.361262 +0200 CEST m=+0.000106168
1685721901361262000
```

The timestamp can be same for two calls to id generation function,
but the `.String()` method provides both
[wall-clock and monotonic clock](https://pkg.go.dev/time#hdr-Monotonic_Clocks).
The result is that the first part of both strings are the same,
but the value after m makes the time returned from this functino monotoinc,
which means no two calls can have the same time.

</details>
