# reinvocation

This scenario configures two mutating webhooks with a `reinvocationPolicy` set to `IfNeeded`. Both webhooks will mutate the object, causing the other webhook to be triggered again,

![Preview of the reinvocation scenario in action](../../assets/reinvocation.gif)

## Impact

Each webhook is triggered a couple times and then no more. The api-server keeps track of how many times it has called a specific webhook and avoid calling it endlessly.

From the Kubernetes documentation:

> if additional invocations result in further modifications to the object, webhooks are not guaranteed to be invoked again.

## Solutions

* None needed
