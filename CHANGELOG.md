# Change log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- New feature: It is now possible to pass additional options to middlewares.
In particular, the first available option is to provide Content-Type selectors
to validators.

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
