<a name="v0.0.1"></a>

## [v0.0.1](https://github.com/lunarianss/Hurricane.git/compare/v0.1.1...v0.0.1) (2024-10-14)

### Bug Fixes

- **make gen:** support generate bizCodes which by template codes to analysis ast tree and runtime types

### Features

- fork pkg/errors as Go Err Handling
- install go-gitlint
- complete make functions like release, swagger, image deploy
- **errors:** support err detail with wrap stack, biz code ,message and aggregate

<a name="v0.1.1"></a>

## [v0.1.1](https://github.com/lunarianss/Hurricane.git/compare/v0.1.0...v0.1.1) (2024-10-12)

### Bug Fixes

- **test release:** test releash a new pre release version

<a name="v0.1.0"></a>

## v0.1.0 (2024-10-12)

### Bug Fixes

- **test release:** test releash a new pre release version

### Features

- inject configuration files, command line parameters into the application
- add server shutdown gracefully
- run server with running config for static options
- add static configuration options for hurricane
- **app:** command line parameters support to replace from \_ to -
- **log:** Wrapper logging library based on zap, with two create modes: withOptions and newOptions
- **make:** test and fix make build and make cover command
- **make:** test and fix make lint command
- **make:** test integrated make add-copyright and format command parameters
- **make:** test make gen command temporary for linux development machine
- **third_party:** add maxprocess as fork third party and adopt log
