# j2xrp

A reverse proxy that converts a JSON request to an XML request

We host an instance of this application on https://j2xrp.herokuapp.com/.

## Installation

    $ go install github.com/flaccid/j2xrp

## Usage

### Example

Client request use with the hosted Heroku application:

```
curl -vvv \
  -H "Content-Type: application/json" \
  -X POST -d '{"username":"xyz","password":"xyz"}' \
    https://j2xrp.herokuapp.com/
```

### Build

    $ make

The resultant compiled binary is located at `bin/github.com/flaccid/j2xrp`.

#### Static Binary

To build a fully static binary (useful for docker):

    $ CGO_ENABLED=0 go build -a -installsuffix cgo -o bin/j2xrp .

#### Docker

    $ docker build -t flaccid/j2xrp .

### Run

See the usage options with `j2xrp help`. If you have installed the package, make sure `"$GOPATH/bin"` is within your `$PATH`.

A simple example:

    $ j2xrp --scheme https wstunnel10-1.rightscale.com

You can also just run from the main entrypoint locally without building:

    $ go run main.go


License and Authors
-------------------
- Author: Chris Fordham (<chris@fordham-nagy.id.au>)

```text
Copyright 2017, Chris Fordham

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
