[![Build status](https://github.com/aledegano/urlprober/workflows/ci/badge.svg)](https://github.com/aledegano/urlprober/actions) [![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

# UrlProber

Periodically GET a configurable URL.

## Why?

1. Because.

2. As an exercise in microservices, multi-arch containers, Helm charts.

3. Logging and metrics.


## How-to use

Urlprober is configured through environmental variables:

 * `URLPROBER_URL`: *Required* The base URL to periodically GET

 * `URLPROBER_INTERVAL`: *Required* The interval _in seconds_ of probing repetition

 * `URLPROBER_QUERY`: *Optional* An additional query that will be appeneded to the base URL. It will not appear in logs and metrics to prevent leaking secrets. Defaults to `""`.

 * `URLPROBER_REQUIRED_STATUS`: *Optional* A list of HTTP status codes (separated by comma `,`) that are to be considered valid responses. Defaults to `[200]`.
