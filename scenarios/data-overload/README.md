# data-overload

In this scenario we're going to return as much data as possible in the response to the api-server. The achieve this we're piping random data from the `crypto/rand` package to the response writer. This scenario is configured to trigger from creation of Pods only.

![Preview of the data-overload scenario in action](../../assets/data-overload.gif)

## Impact

The api-sever completely locks up and stops responding to any further api calls.

## Solutions

* Restarting of api-server temporary fixes the issue until the webhook is next triggered.
* Reducing the `timeoutSeconds` seems to help the api-server handle the webhook calls but the default 10s causes it to lockup.
