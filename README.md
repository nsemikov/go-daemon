[![License][lic-img]][lic] [![FOSSA Status][fossa-img]][fossa] [![pkg.go.dev reference][go.dev-img]][go.dev] [![Build Status][ci-img]][ci] [![Go Report Card][report-img]][report] [![Release][release-img]][release]

# go-daemon

![Logo](https://github.com/stdatiks/go-daemon/blob/master/.github/images/go-daemon.1280x640.png?raw=true)

A daemon package for use with Go services without any dependencies (except for `golang.org/x/sys/windows`)


## Features

* Install and uninstall service file on Linux (SystemD, SystemV, UpStart), MacOS and FreeBSD
* More control: start, stop, restart, reload, pause and continue services
* Unified interface for all supported OS


[go.dev-img]: https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white
[go.dev]: https://pkg.go.dev/github.com/stdatiks/go-daemon
[doc-img]: https://img.shields.io/badge/go-documentation-blue.svg
[doc]: https://godoc.org/github.com/stdatiks/go-daemon
[ci-img]: https://img.shields.io/travis/com/stdatiks/go-daemon.svg
[ci]: https://travis-ci.com/stdatiks/go-daemon
[cov-img]: https://img.shields.io/codecov/c/github/stdatiks/go-daemon.svg
[cov]: https://codecov.io/gh/stdatiks/go-daemon
[report-img]: https://goreportcard.com/badge/github.com/stdatiks/go-daemon
[report]: https://goreportcard.com/report/stdatiks/go-daemon
[release-img]: https://img.shields.io/badge/release-v0.2.1-1eb0fc.svg
[release]: https://github.com/stdatiks/go-daemon/releases/tag/v0.2.1
[lic-img]: https://img.shields.io/badge/License-MIT-blue.svg
[lic]: https://opensource.org/licenses/MIT
[fossa-img]: https://app.fossa.com/api/projects/git%2Bgithub.com%2Fstdatiks%2Fgo-daemon.svg?type=shield
[fossa]: https://app.fossa.com/projects/git%2Bgithub.com%2Fstdatiks%2Fgo-daemon?ref=badge_shield
