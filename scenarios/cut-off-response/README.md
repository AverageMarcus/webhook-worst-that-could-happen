# cut-off-response

In this scenario we're going to set the `Content-Length` response header to be longer than the actual response body.

![Preview of the cut-off-response scenario in action](../../assets/cut-off-response.gif)

## Impact

The api-server returns an error of unexpected EOF.

## Solutions

None needed, api-server handles the misconfiguration and returns an error when calling the webhooks.
