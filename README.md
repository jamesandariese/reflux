# reflux

[![Build Status](https://travis-ci.org/jamesandariese/reflux.svg?branch=master)](https://travis-ci.org/jamesandariese/reflux)

Shove junk into influx...  more easily than using the pretty obnoxious library
from influx.

## Usage

If your plans include a long running application pushing metrics regularly, use
the NewClient facility.  Otherwise, use SendPointWithJsonTags to automatically
setup the client, point queue, and tear it down when done.

If you're going to use this from a command line tool, you can use `PrepareFlags`
to prep an influx-url and an influx-json-tags argument which will be used by
`SendPointUsingFlags` to make things that much easier.

See cmd/reflux/main.go for a sample.
