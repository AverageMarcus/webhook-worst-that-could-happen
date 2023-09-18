# different-content-type

In this scenario we're going to return a valid `Content-Type` header but one that doesn't match the content being returned.

![Preview of the different-content-type scenario in action](../../assets/different-content-type.gif)

## Impact

The api-server rejects any webhook responses that aren't JSON (or YAML) regardless of their `Content-Type` header.

## Solutions

None. The api-server doesn't seem to make use of the `Content-Type` header and instead just attempts to marshal the body string into a struct - supporting both JSON and YAML due to the library used in code.
