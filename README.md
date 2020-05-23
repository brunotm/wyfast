# wyfast

wyfast is a small 64-bit non-cryptographic hash function.

It is simpler and faster variation of wyhash v5  https://github.com/wangyi-fudan/wyhash

## Performance

* Full [smhasher log](./doc/smhasher.txt)

```
[[[ Speed Tests ]]]

Bulk speed test - 262144-byte keys
Alignment  7 -  8.690 bytes/cycle - 24861.65 MiB/sec @ 3 ghz
Alignment  6 -  8.528 bytes/cycle - 24399.15 MiB/sec @ 3 ghz
Alignment  5 -  8.690 bytes/cycle - 24862.28 MiB/sec @ 3 ghz
Alignment  4 -  8.567 bytes/cycle - 24510.46 MiB/sec @ 3 ghz
Alignment  3 -  8.690 bytes/cycle - 24863.57 MiB/sec @ 3 ghz
Alignment  2 -  8.547 bytes/cycle - 24453.99 MiB/sec @ 3 ghz
Alignment  1 -  8.691 bytes/cycle - 24864.79 MiB/sec @ 3 ghz
Alignment  0 -  8.730 bytes/cycle - 24978.13 MiB/sec @ 3 ghz
Average      -  8.642 bytes/cycle - 24724.25 MiB/sec @ 3 ghz

Small key speed test -    1-byte keys -    14.20 cycles/hash
Small key speed test -    2-byte keys -    15.00 cycles/hash
Small key speed test -    3-byte keys -    15.00 cycles/hash
Small key speed test -    4-byte keys -    12.00 cycles/hash
Small key speed test -    5-byte keys -    19.00 cycles/hash
Small key speed test -    6-byte keys -    19.00 cycles/hash
Small key speed test -    7-byte keys -    19.00 cycles/hash
Small key speed test -    8-byte keys -    12.00 cycles/hash
Small key speed test -    9-byte keys -    20.15 cycles/hash
Small key speed test -   10-byte keys -    20.00 cycles/hash
Small key speed test -   11-byte keys -    20.00 cycles/hash
Small key speed test -   12-byte keys -    20.00 cycles/hash
Small key speed test -   13-byte keys -    20.00 cycles/hash
Small key speed test -   14-byte keys -    20.00 cycles/hash
Small key speed test -   15-byte keys -    20.00 cycles/hash
Small key speed test -   16-byte keys -    20.00 cycles/hash
Small key speed test -   17-byte keys -    20.89 cycles/hash
Small key speed test -   18-byte keys -    21.08 cycles/hash
Small key speed test -   19-byte keys -    20.88 cycles/hash
Small key speed test -   20-byte keys -    20.00 cycles/hash
Small key speed test -   21-byte keys -    20.24 cycles/hash
Small key speed test -   22-byte keys -    21.14 cycles/hash
Small key speed test -   23-byte keys -    20.00 cycles/hash
Small key speed test -   24-byte keys -    20.00 cycles/hash
Small key speed test -   25-byte keys -    20.00 cycles/hash
Small key speed test -   26-byte keys -    20.00 cycles/hash
Small key speed test -   27-byte keys -    20.00 cycles/hash
Small key speed test -   28-byte keys -    20.00 cycles/hash
Small key speed test -   29-byte keys -    20.00 cycles/hash
Small key speed test -   30-byte keys -    20.00 cycles/hash
Small key speed test -   31-byte keys -    20.00 cycles/hash
Average                                    19.019 cycles/hash

[[[ 'Hashmap' Speed Tests ]]]

std::unordered_map
Init std HashMapTest:     956.344 cycles/op (235886 inserts, 1% deletions)
Running std HashMapTest:  270.628 cycles/op (12.0 stdv)

greg7mdp/parallel-hashmap
Init fast HashMapTest:    320.245 cycles/op (235886 inserts, 1% deletions)
Running fast HashMapTest: 241.417 cycles/op (19.8 stdv)  ....... PASS
```