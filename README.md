Migration Warning
=================
This repository has been migrated to github.com/m3db/m3. It's contents can be found at github.com/m3db/m3/src/ctl. Follow along there for updates. This repository is marked archived, and will no longer receive any updates.

# m3ctl [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov]

Configuration controller for the M3DB ecosystem. Provides an http API to perform CRUD operatations on
the various configs for M3DB compontents.

### Run the R2 App

```bash
git clone --recursive https://github.com/m3db/m3ctl.git
cd m3ctl
glide install -v
make && ./bin/r2ctl -f config/base.yaml
open http://localhost:9000
```

**UI**
`/`

**API Server**
`/r2/v1`

**API Docs (via Swagger)**
`public/r2/v1/swagger`


<hr>

This project is released under the [Apache License, Version 2.0](LICENSE).

[doc-img]: https://godoc.org/github.com/m3db/m3ggregator?status.svg
[doc]: https://godoc.org/github.com/m3db/m3ctl
[ci-img]: https://travis-ci.org/m3db/m3ctl.svg?branch=master
[ci]: https://travis-ci.org/m3db/m3ctl
[cov-img]: https://coveralls.io/repos/m3db/m3ctl/badge.svg?branch=master&service=github
[cov]: https://coveralls.io/github/m3db/m3ctl?branch=master
