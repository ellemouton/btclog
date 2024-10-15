btclog
======

[![Build Status](http://img.shields.io/travis/btcsuite/btclog.svg)](https://travis-ci.org/btcsuite/btclog)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/btcsuite/btclog)

Package btclog defines a logger interface and provides a default implementation
of a subsystem-aware leveled logger implementing the same interface.

## Installation

```bash
$ go get github.com/btcsuite/btclog
```

## GPG Verification Key

All official release tags are signed by Conformal so users can ensure the code
has not been tampered with and is coming from the btcsuite developers.  To
verify the signature perform the following:

- Download the public key from the Conformal website at
  https://opensource.conformal.com/GIT-GPG-KEY-conformal.txt

- Import the public key into your GPG keyring:
  ```bash
  gpg --import GIT-GPG-KEY-conformal.txt
  ```

- Verify the release tag with the following command where `TAG_NAME` is a
  placeholder for the specific tag:
  ```bash
  git tag -v TAG_NAME
  ```

## Usage

`btclog.NewSLogger` can be used to construct a new `btclog.Logger` type which
can then be used for logging calls. The `NewSLogger` function expects to be
initialised with a type that implements the `btclog.Handler` interface which is
responsible for writing logging records to a backend writer. Callers may provide
their own `Handler` implementations (for example, the standard library `slog`
package provides some handler implementations such as a JSON Handler and a Text
Handler) or else they may use the default one provided with this package:
`DefaultHandler`.

Example usage:

```
	// Create a new DefaultHandler that writes to stdout and set the
	// logging level to Trace.
	handler := btclog.NewDefaultHandler(os.Stdout)
	handler.SetLevel(btclog.LevelTrace)

	// Use the above handler to construct a new logger.
	log := btclog.NewSLogger(handler)

	/*
		2024-09-18 11:53:03.287 [INF]: An info level log
	*/
	log.Info("An info level log")

	// Create a subsystem logger with no timestamps.
	handler = btclog.NewDefaultHandler(os.Stdout, btclog.WithNoTimestamp())
	log = btclog.NewSLogger(handler.SubSystem("SUBS"))

	/*
		[INF] SUBS: An info level log
	*/
	log.Info("An info level log")

	// Include log source.
	handler = btclog.NewDefaultHandler(
		os.Stdout,
		btclog.WithCallerFlags(btclog.Lshortfile),
		btclog.WithNoTimestamp(),
	)
	log = btclog.NewSLogger(handler.SubSystem("SUBS"))

	/*
		[INF] SUBS main.go:36: An info level log
	*/
	log.Info("An info level log")

	// Attach attributes to a context. This will result in log lines
	// including these attributes if the context is passed to them. Also
	// pass in an attribute at log time.
	log = btclog.NewSLogger(btclog.NewDefaultHandler(
		os.Stdout, btclog.WithNoTimestamp(),
	).SubSystem("SUBS"))
	ctx := btclog.WithCtx(context.Background(), "request_id", 5)

	/*
		[INF] SUBS: A log line with context values request_id=5 another_key=another_value
	*/
	log.InfoS(ctx, "A log line with context values", "another_key", "another_value")
```

## License

Package btclog is licensed under the [copyfree](http://copyfree.org) ISC
License.
