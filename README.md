[![Build
Status](https://travis-ci.org/5Sigma/spyder.svg?branch=master)](https://travis-ci.org/5Sigma/spyder)

# Spyder
API Testing and Request Framework

## Installation

### OSX

On OSX, spyder can be install with brew:

```
brew install 5sigma/tap/spyder
```

### Linux 

Download the linux package from for the latest release:

https://github.com/5Sigma/spyder/releases/latest

### Windows

Windows binaries can be found in the release:

https://github.com/5Sigma/spyder/releases/latest

## API Testing and Requests

Spyder provides an easy interface to make and test API endpoints from within the terminal.
It uses simple JSON configuration files which are meant to be versioned. It also has built in
scripting support for dynamically modifying endpoint configuration and performing certain
types of tasks, such as appending authentication, persisting tokens, etc.


## Start a project

Spyder expects the configuration for a set of endpoints to live in its own folder structure.
This folder can then be versioned to allow for easy maintain and team synchronization. 
The `init` command will generate the folder structure for you. Just run it in an empty folder.

```
spyder init
```

## Configuring an endpoint

Endpoints are defined in JSON files inside the endpoints folder. For example to
create an product list endpoint we might place a file at 
`project_root/endpoints/products/list.json`. A simple endpoint looks like this:

```js
  {
    "url": "http://example.com/api/products"
    "method": "GET"
  }
```

## Handling parameters

Request parameters can be passed using the "data" node in the configuration.
For GET requests these are encoded and added to the url when the request is
made.

For POST requests the node is submitted as stringified JSON in the post body.

### Example: A simple GET request

```js
  {
    "url": "http://example.com/api/products"
    "method": "GET",
    "data": {
      "pageSize": 10,
      "page": 1
    }
  }
```


### Example: A simple POST request

```js
  {
    "url": "http://example.com/api/products"
    "method": "POST",
    "data": {
      "data": {
        "attributes": {
           "name": "My Product",
           "price": 399.99
        }
      }
    }
  }
```

## Handling dynamic data

The easiest way of handling dynamic data is by using variables directly inside
the configuration. There are two configuration files: 

- config.json - This holds the base configuration for the project and can store
    variables that may change over time.
- config.local.json - This file is meant for the specific user and should not be
    checked in to versioning.  It allows overriding the default global variables
    as well as adding ones specific to the person. Such as test account
    credentials, tokens,etc.

If your config.local.json looked something like this:

```js
  {
    "variables": {
      "token": "ababab123121"
    }
  }
```

You could then make a endpoint configuration that uses it. Variables are
specified with a preceding `$` and are expanded into their values before the
request is made.

```js
  {
    "url": "http://exmaple.com/api/auth",
    "method": "POST",
    "headers": {
      "Content-Type": "application/json"
    },
    "data": {
      "token": "$token"
    }
  }
```

## Advanced requests using scripting

Spyder has a built in JavaScript interpreter. It supports two types of hooks.

- **On Complete** - When the request is made and a response is received on
    complete scripts can be fired. These are useful for setting variables to the
    result of a request.

- **Transform** - Transform scripts are given the request before it is sent and
    given the option to transform it in some way. This allows more dynamic
    control over the request. Things like injecting authorization tokens,
    or conditional data.


### Example: Handling an endpoint that requires an HAMC signature 

Consider having a set of endpoints where you must first authenticate to an auth
endpoint and receive a session token. Then use this token to sign future
requests using HMAC.

The authorization endpoint might look like this:

```js

{
  "url": "http://example.com/auth",
  "method": "POST",
  "headers": {
    "Content-Type": "application/json"
  },
  "onComplete": ["storeAuthSession"],
  "data": {
    "token_id": "aaabab12312",
    "token_secret": "aaabab12312",
  }
}
```

This utilizes an onComplete script to save out the token. A
`scripts/storeAuthSession.js` file might look like:

```js
data = JSON.parse($response.body);
$variables.set('session_token_id', data.body.data.session.session_id);
$variables.set('session_token_secret', data.body.data.session.session_secret);
```

A standard request to the API then might look like: 

```js
  {
    "url": "http://exmaple.com/api/products",
    "method": "get",
    "trasnform": ["signRequest"]
  }
```

This request uses a transform script located at `scripts/signRequest.js`. That
could look like:

```js
signature = $hmac($variables.get('session_token_secret'), $request.body);
$request.headers.set('Authorization', $variables.get('session_token_id') + ':' + signature)
```

For more information on scripting see the [Scripting Reference](https://github.com/5Sigma/spyder/wiki/Script-Reference)
