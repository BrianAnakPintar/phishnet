Fishnet is our DSL for processing the FilterChain.
In particular, it contains a parser as well as a template describing how it should work.

### About Fishnet
The following is how we describe the syntax of our DSL

```fn
<filter name>:[
    <map of params>
]
...
```

The parser will then read the config file (ends with .fn) and adds each line as a filter to the chain