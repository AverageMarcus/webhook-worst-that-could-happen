# label-overwrite

This scenaio configures a mutating webhook that updates a pods labels with a new entry, but accidentally replaces the existing.

![Preview of the label-overwrite scenario in action](../../assets/label-overwrite.gif)

## Impact

Because all the pods labels are mistakenly replaced with whatever the webhook sets the `matchLabels` selectors will no longer match. This would lead to:
* ReplicaSets wont match so pods will be recreated infinitely until the cluster is at full capacity.
* Services wont match so routing will fail for new pods.
* The cluster will eventually fill up with pending pods and no new workloads can be created until the existing have been deleted.
* If autoscaling is enabled the cluster will scale up to max, incuring costs.

## Solutions

This is a tricky one. There's not much that can be done to block this action exactly. The only real solution here is to ensure webhooks are tested fully before deploying to a production environment.
