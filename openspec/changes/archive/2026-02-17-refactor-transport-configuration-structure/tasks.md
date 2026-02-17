## 1. Configuration Model Refactor

- [x] 1.1 Replace scalar transport field with a structured transport configuration object containing `type` and `streamable_http` settings
- [x] 1.2 Rename HTTP configuration type/fields to streamable HTTP naming and nest under transport in YAML/koanf tags
- [x] 1.3 Update default configuration values to populate `transport.type` and `transport.streamable_http.port`

## 2. Config Loading and Validation

- [x] 2.1 Update configuration parsing and validation to read and validate `transport.type`
- [x] 2.2 Ensure invalid transport type errors explicitly list valid values (`stdio`, `streamable_http`)
- [x] 2.3 Ensure legacy scalar `transport` / top-level `http` shapes are no longer supported by tests and parsing behavior

## 3. Runtime Transport Wiring

- [x] 3.1 Update server transport selection to use nested config (`config.Transport.Type`)
- [x] 3.2 Update streamable HTTP transport startup to use nested port (`config.Transport.StreamableHTTP.Port`)

## 4. Tests and Examples

- [x] 4.1 Refactor configuration unit tests to use nested transport YAML and struct expectations
- [x] 4.2 Refactor server tests to construct nested transport configuration
- [x] 4.3 Update example configuration and README snippets to the new transport object format
- [x] 4.4 Run `make test` and fix any failures caused by the configuration refactor
