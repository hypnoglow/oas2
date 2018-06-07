# Change log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

That's **a lot** of changes in this upcoming release :)

### Added

- New `oas.Wrap()` option for `oas.NewRouter()` applies a middleware that wraps
the router. That means it executes before the actual routing, so you can modify
its behaviour, e.g. introducing new routes (handling 404 errors) or methods (handling
OPTIONS method for CORS). See `oas.Wrap()` [documentation](https://github.com/hypnoglow/oas2/blob/b0d734259c9ebab2bb7196b49a48e3e3c0ada79a/router.go#L141)
for more information. Note that those middleware also applies to a spec served 
by router, when enabled via `oas.ServeSpec()` option.
- `DefaultBaseRouter()` got exposed for informational purposes.
- `DefaultExtractorFunc()` was added, so users don't need to import [chi](https://github.com/go-chi/chi)
package if they use default router. Thus, you can just `oas.Use(oas.PathParameterExtractor(oas.DefaultExtractorFunc))`

### Changed

- Middleware order got adjusted. Now, when you pass middleware to `oas.NewRouter()` using
`oas.Use()` option, they get applied exactly in the same order. See `oas.Use()`
option [documentation](https://github.com/hypnoglow/oas2/blob/b0d734259c9ebab2bb7196b49a48e3e3c0ada79a/router.go#L167)
and [this](https://github.com/hypnoglow/oas2/blob/b0d734259c9ebab2bb7196b49a48e3e3c0ada79a/e2e/middleware_order/main_test.go#L32)
test for an example. 
- `oas` package now introduces wrapper types `oas.Document` and `oas.Operation`. All functions
that previously exposed parameters from other libraries from [go-openapi](https://github.com/go-openapi)
now use these types. 

    For example, `oas.LoadFile()` now returns `*oas.Document` instead of `*loads.Document`.
    Thus, `oas.NewRouter()` accepts `*oas.Document` instead of `*loads.Document`.

    The main purpose for this change is that most users only need those types to
    use within this library, and they had to import [go-openapi](https://github.com/go-openapi)
    libraries just to pass variables around. Also, you still can access underlying types from [go-openapi](https://github.com/go-openapi)
    if you need.
    
    Functions that had their signatures changed:
    
    - `oas.LoadFile()`
    - `oas.NewRouter()`
    - `oas.WithOperation()`
    - `oas.GetOperation()`
    - `oas.MustOperation()`
     
- Adjusted string representation of variables of type `oas.ValidationError`.
- All package middlewares that require `oas.Operation` to work now **panic**
if operation is not found in the request context. Previously, if middleware
cannot find operation in request context, it would just silently skip the validation.
That behavior is undesirable because validation is important. So, the change 
is done to actually notify package user that he is doing something wrong. 

    Affected middlewares:
    
    - `oas.PathParameterExtractor()`
    - `oas.QueryValidator()`
    - `oas.BodyValidator()`
    - `oas.ResponseBodyValidator()`

## [0.6.1] - 2018-05-28

### Fixed

- Fixed a bug when `oas.LoadFile` was not returning expanded spec.

## [0.6.0] - 2018-05-22

### Added

- New feature: It is now possible to pass additional options to middlewares.
In particular, the first available option is to provide Content-Type selectors
to validators.

### Changed

- **Breaking Change** Router now accepts a document instead of a spec. This 
allows, among other things, to access original spec, which is used when serving 
router spec.
- **Breaking Change** `LoadSpec` function is replaced with `LoadFile`, which 
works almost the same way, but provides more options, for example to set
spec host when loading.
- **Breaking Change** `DecodeQuery` now accepts request and extracts operation 
spec from it, so developer no longer need to extract operation params from request
only to use them in `DecodeQuery`. Old function moved to `DecodeQueryParams`.

### Fixed

- Fixed a bug about type conversion for `default` values of query parameters.

## [0.5.1] - 2018-03-14

### Fixed

- Fixed a bug that led to serving spec with duplicate parameters.

## [0.5.0] - 2018-03-13

### Added

- New feature: add option [`ServeSpec`](https://github.com/hypnoglow/oas2/blob/4b7ce7cc55bdd7cbb66e94e8af94f3dd08e8fc01/router.go#L127) for router to serve its OpenAPI spec under the base path.
- New tool: [oas-expand](https://github.com/hypnoglow/oas2/tree/7678e995b788570a0483e667e030f8c7166a6681/cmd/oas-expand) expands all `$ref`s in spec to improve startup time of
the application.
- OpenAPI parameter with `type: string` and `format: uuid` is now supported.

### Changed

- **Breaking Change** `JsonError` renamed to`JSONError`.
- Replaced internally-used self-written response recorder with [chi's](https://github.com/go-chi/chi/blob/master/middleware/wrap_writer18.go#L12) 
- The project now has more advanced CI checks.
