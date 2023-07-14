*go <-> py*
- https://blog.marlin.org/cgo-referencing-c-library-in-go
- https://github.com/DataDog/go-python3/blob/03d93fb21b679303a8622e983080c22d0126e74e/dict_test.go
- https://www.datadoghq.com/blog/engineering/cgo-and-python/
- https://github.com/sbinet/go-python/blob/5167a7b23e1d6ac26f3529a48ae01b9ffb62e1ed/none.go
- https://github.com/jshiwam/cpy3x/blob/e05a8296b36246dedf4b9a472c4719fe81f4df98/pycore/dict.god
- https://poweruser.blog/embedding-python-in-go-338c0399f3d5
- https://github.com/tebeka/talks/blob/master/embed-py-fosdem/glue.c
- https://pythonextensionpatterns.readthedocs.io/en/latest/index.html
- https://github.com/ardanlabs/python-go/blob/fea5e5f973a0d9dbeb5aea80f6d07bbd3e6c1758/py-in-mem/outliers.go#L94
- https://www.ardanlabs.com/blog/2020/09/using-python-memory.html
- https://github.com/ardanlabs/python-go/blob/master/py-in-mem/Makefile
- https://www.netlify.com/blog/2021/03/18/tracking-down-a-cgo-crash-in-production/
- https://docs.python.org/3.8/c-api/init.html?highlight=py_finalize#c.Py_FinalizeEx
- https://python-list.python.narkive.com/rJT7Xh3Q/when-embedding-python-how-do-you-redirect-stdout-stderr
- https://www.datadoghq.com/blog/engineering/cgo-and-python/#the-dreadful-global-interpreter-lock
- https://dev.pippi.im/writing/cgo-and-python/
- https://github.com/go-python/cpy3
- https://github.com/TykTechnologies/tyk/blob/master/dlpython/binding.go
- https://github.com/kluctl/go-embed-python
- https://github.com/go-python/gpython/tree/main/examples/embedding
- https://github.com/aadog/py3-go
- https://github.com/go-python/gopy

*Codegen tools:*
- https://pkg.go.dev/github.com/drone/sqlgen
- https://github.com/cohesivestack/valgo
- https://preslav.me/2023/03/07/reasons-against-sqlc/
- https://traefik.io/blog/mocktail-the-mock-generator-for-strongly-typed-mocks/

*DIP:*
- https://go.dev/blog/wire
- https://github.com/google/wire/blob/main/_tutorial/README.md
- https://nohowtech.com/posts/dependency-injection-2023-04-18/

*Migrations and SQL:*
- https://github.com/pressly/goose
- https://github.com/amacneil/dbmate

*JSON/Protobuf:*
- https://www.cockroachlabs.com/blog/high-performance-json-parsing/
- https://dave.cheney.net/high-performance-json.html
- https://vincent.bernat.ch/en/blog/2023-dynamic-protobuf-golang

*Perf:*
- https://go.dev/blog/pgo-preview
- https://dave.cheney.net/2020/04/25/inlining-optimisations-in-go

*ldflags/compile:*
- https://levelup.gitconnected.com/a-better-way-than-ldflags-to-add-a-build-version-to-your-go-binaries-2258ce419d2d

*finite state machine:*
- https://kyleshevlin.com/guidelines-for-state-machines-and-xstate
- https://github.com/looplab/fsm
- https://tproger.ru/translations/finite-state-machines-theory-and-implementation/
- [] Think about proper StateMap data structure

*ratelimiter:*
- https://blog.bytebytego.com/p/rate-limiting-fundamentals?utm_source=substack&utm_medium=email

*Patterns / SAGA:*
- https://dormoshe.io/trending-news/saga-pattern-made-easy-4j42-62197?utm_source=twitter&utm_campaign=twitter

*testing tools:*
- https://evilmartians.com/chronicles/go-integration-testing-with-courage-and-coverage
- https://github.com/jinzhu/now

*kafka:*
- https://github.com/segmentio/kafka-go
- https://github.com/twmb/franz-go
- https://github.com/Shopify/sarama
- github.com/confluentinc/confluent-kafka-go/kafka
