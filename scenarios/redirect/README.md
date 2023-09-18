# redirect

This scenario responds to all webhook requests with a redirect to a service that infinitely redirects the client.

![Preview of the redirect scenario in action](../../assets/redirect.gif)

## Impact

The api-server stops following the redirects after 10 redirects.

## Solutions

* None - the api-server handles this and returns an error
