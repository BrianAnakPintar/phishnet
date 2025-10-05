# Phishnet

Phishnet is a small, modular Go-based URL proxy framework focused on detecting and blocking phishing and other unsafe URLs. It provides a filter-oriented architecture so different detection strategies can be plugged in, configured, and chained together.

The project includes example filters (regex-based, third-party services, and generative-model checks), a minimal UI component, and a lightweight runtime to evaluate URLs against a configured filter chain.

### Key ideas

- Filters are small components that implement two methods: `Configure(map[string]string) error` and `Run(string) (FilterResult, error)`.
- Filters are registered by name using `RegisterFilter(name, factory)` and created with `CreateFilter(name)`.
- `FilterResult` carries the decision (allow/block) and a human-readable reason.

### Included components

- `filters/` — collection of filters and the filter factory/registry:
  - `filter.go` — core interfaces and registration helpers (`Filter`, `FilterResult`, `RegisterFilter`, `CreateFilter`).
  - `Regex.go` — a basic pattern-based filter that marks URLs matching configured regular expressions as unsafe.
  - `GoogleSafeBrowsing.go` — a filter that queries the Google Safe Browsing API (v4). Requires an API key; see notes below.
  - `PhishTank.go` — (included) a filter that can check PhishTank data (project contains a `data/verified_online.*` dataset used by some filters).
  - `GeminiFilter.go` — an example integration that asks a generative model (Google Gemini style) to classify a URL as phishing or not. This illustrates calling an external LLM-style service and interpreting short YES/NO replies.

- `fishnet/` — helper runtime and DSL for building and running filter chains. Contains a simple parser and a `filterchain.go` that wires multiple filters together.
- `ui/` — minimal UI pieces (popup and theme) used by the local UI.
- `data/` — supporting data files used by filters (e.g. `verified_online.csv` / `verified_online.gob` and small images used by the UI).
- `main.go` — entry point for running the program (local CLI / runtime).

### Build & run

You need Go installed (tested with recent Go 1.20+). From the repository root:

- Build: `go build ./...` or `go build -o phishnet ./`
- Run: `go run ./` (or `./phishnet` after building)

### Configuration

- Filters are configured at runtime by the surrounding runtime or DSL in `fishnet/`. Each filter exposes a `Configure(map[string]string)` method — for example, the Google Safe Browsing filter expects `map["API_KEY"] = "<YOUR_KEY>"`.
- API keys and secrets should never be committed to version control. Prefer environment variables, a local config file excluded from git, or a secrets manager.

### Notes on third-party integrations

- Google Safe Browsing: you must enable the Safe Browsing API in Google Cloud and create credentials. The filter expects a key (or token) provided to `Configure`.
- Gemini / Generative model: the example filter demonstrates how to call a generative model endpoint and parse a concise YES/NO reply. Actual deployment will require proper credentials and adherence to the provider's API/usage rules.

### Extending Phishnet

To add a new filter:

1. Create a type that implements the `Filter` interface in the `filters/` package.
2. Provide a constructor that returns `Filter` (e.g. `NewMyFilter()`).
3. Register the filter in an `init()` function with `RegisterFilter("MyFilterName", NewMyFilter)`.
4. Configure and wire it into your filter chain via the runtime or the DSL in `fishnet/`.

### Development notes

- Tests: No tests are added since this is for a hackathon.
- Envoy: Yes. This project was inspired by Envoy's filter chaining mechanism, it's simply very cool imo.
