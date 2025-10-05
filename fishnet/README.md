Fishnet is our DSL for processing the FilterChain.
In particular, it contains a parser as well as a template describing how it should work.

### About Fishnet
The following is how we describe the syntax of our DSL

```
;; Comments are denoted with 2 semicolon (Like Racket!)

<name of filter 1>:[
    <map of params>
]

<name of filter 2>:[
    <map of params>
]
...
```

An example for the DSL can be found here:

```
;; PhishNet Bootstrap filters

PhishTankFilter:[
    bad1 = http://malicious.example/
    bad2 = https://phish.example/login
    bad3 = http://example.com/rickroll
    bad4 = https://badactor.test/steal
]

RegexFilter:[
    BlockYoutubeRegex = youtube
]

```

The parser will then read the config file (ends with .fn) and adds each line as a filter to the chain