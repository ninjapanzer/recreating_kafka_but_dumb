## Stuff thats Dumb

### Links
- https://github.com/jmhodges/levigo
- https://medium.com/@viktordev/build-a-concurrent-tcp-server-in-go-with-graceful-shutdown-include-unit-tests-2e2d63ee2161
- https://cbor.me/?diag=%22pii%22%20=%20(%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20age:%20int,%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20name:%20tstr,%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20employer:%20tstr,%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20)
- https://www.ietf.org/archive/id/draft-ietf-cbor-time-tag-05.html#name-table-of-contents-8
- https://datatracker.ietf.org/doc/html/rfc7049#section-7.2
- https://github.com/fxamacker/cbor?tab=readme-ov-file#struct-tags
- https://pkg.go.dev/github.com/ugorji/go/codec#NewEncoder

Producer Management:
Each producer is merged into a single event file writing concurrently with a mutex

Consumer Management:
Multiple Consumers is hard-ish
So when we have a single consumer we want to stream from the beginning if they are new and from a target offset if
they have been registered before.

So this means that we need a consistent way to persist a consumers offset both as a read position but also as a rejoin position
V1 impelements a per consumer stream offset tracking and makes it the consmers responsiblity
to provide its offset if starting mid-stream. In retrospect this is kinda right but we need to mirror a consumer group.
A consumer group being different from a new consumer. We can track multiple file handles and read offsets but we also
need a global position when joining a consumer group so consumer registration needs to hold a handle to a file but also must share
mutex for a consumer group so the moment a poll request is received we allocate
an offset to that consumer. This probably means we can allocate an offset in lines instead of bytes.

We also need a way to ack a consumer poll operation so we can mark items in the consumer group as consumed.