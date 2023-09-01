# open-connection

This scenario sets the response headers indicating the response will be returned in streaming chunks and then waits indefinitely until the client hangs up.

![Preview of the open-connection scenario in action](../../assets/open-connection.gif)

## Impact

Slows down, and eventually blocks, targetted API requests. The `timeoutSeconds` property on the webhook configuration determines how long until the api-server closes the connection and treats it as failed.

## Solutions

* Appropriately setting the `timeoutSeconds` value.
