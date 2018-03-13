# Change log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- New tool: [oas-expand](https://github.com/hypnoglow/oas2/tree/7678e995b788570a0483e667e030f8c7166a6681/cmd/oas-expand) expands all `$ref`s in spec to improve startup time of
the application.
- OpenAPI parameter with `type: string` and `format: uuid` is now supported.

### Changed

- Replaced self-written response recorder with [chi's](https://github.com/go-chi/chi/blob/master/middleware/wrap_writer18.go#L12) 
- The project now has more advanced CI checks.
