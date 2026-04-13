---
# Code Style Guidelines
---

## Use `x` package collection helpers — never write collection boilerplate

Business and top-level code must be readable at a glance. Hand-written loops for mapping, filtering, converting, or grouping collections are boilerplate — use the generic helpers from `github.com/agurinov/gopl/x` instead.

### Available helpers

| Helper | Purpose |
|---|---|
| `x.SliceConvert[T1, T2](in, func(T1) T2)` | Map slice elements to a different type |
| `x.SliceConvertError[T1, T2](in, func(T1) (T2, error))` | Map with error propagation |
| `x.SliceFilter[T](in, func(T) bool)` | Keep elements matching predicate |
| `x.SliceToMap[K, V, E](in, func(E) (K, V))` | Slice → map |
| `x.MapToSlice[K, V, E](in, func(K, V) E)` | Map → slice |
| `x.MapKeys[K, V](in)` | Extract map keys |
| `x.MapConvert[K1,K2,V1,V2](in, func(K1,V1)(K2,V2))` | Convert map keys/values |
| `x.MapFilter[K, V](in, func(K, V) bool)` | Filter map entries |
| `x.MapClone[K, V](in)` | Shallow-clone a map |
| `x.Unique[T](in)` | Deduplicate slice |
| `x.FilterOutEmpty[T](in)` | Remove zero values from slice |
| `x.First[E](s)` / `x.Last[E](s)` | Safe first/last element |
| `x.Flatten[T](in)` | Flatten `[][]T` → `[]T` |
| `x.SliceBatch[T](in, size)` | Split slice into fixed-size batches |
| `x.GroupBy[T, K](in, func(T) K)` | Group slice elements by key |
| `x.FlattenChans[T](chs...)` | Drain and merge buffered channels |
| `x.FlattenErrors(chs...)` | Drain error channels into a single joined error |

### Examples

```go
// BAD — boilerplate loop
ids := make([]string, 0, len(users))
for _, u := range users {
    ids = append(ids, u.ID)
}

// GOOD
ids := x.SliceConvert(users, func(u User) string { return u.ID })
```

```go
// BAD
active := make([]User, 0)
for _, u := range users {
    if u.Active {
        active = append(active, u)
    }
}

// GOOD
active := x.SliceFilter(users, func(u User) bool { return u.Active })
```

```go
// BAD
index := make(map[string]User, len(users))
for _, u := range users {
    index[u.ID] = u
}

// GOOD
index := x.SliceToMap(users, func(u User) (string, User) { return u.ID, u })
```
