[![Go Tests](https://github.com/bitstonks/leftright/actions/workflows/go.yml/badge.svg)](https://github.com/bitstonks/syndi/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bitstonks/leftright)](https://goreportcard.com/report/github.com/bitstonks/leftright)
[![Go Reference](https://pkg.go.dev/badge/github.com/bitstonks/leftright.svg)](https://pkg.go.dev/github.com/bitstonks/leftright)
[![Release](https://img.shields.io/github/release/bitstonks/leftright.svg)](https://github.com/bitstonks/leftright/releases/latest)

# Left Right Concurrency In Go

Left-Right concurrency primitive is a way to achieve arbitrary number wait-free reads from any data structure for the
price of only having a single writer, which is also lock-free, but not wait-free. Waiting time is at most two read
operations. However, writes can be batched, and we only have to pay this waiting cost once we publish the batch.

It does this by keeping two copies of the underlying data structure, one used for reading and one for writing. Whenever
we want (e.g. after every write), we atomically swap them. All readers can immediately start reading from the new side,
while the writer has to wait for outstanding reads to finish before it can write again.

See [example on the lock package][ex-link] for a demo of the lock API.

[ex-link]: https://pkg.go.dev/github.com/bitstonks/leftright/pkg/lock#example-package-Simple

# Specs
* Arbitrary number of concurrent readers.
* Readers are always lock-free and wait-free.
* Single writer (or mutex synchronized).
* Writing is lock-free and wait-free. HOWEVER:
   * There is a wait to publish the written changes to readers.
   * The wait time is at most two reads long.

# TODO
* Performance test to compare it with sync.Map.
* A more realistic example with goroutines.
