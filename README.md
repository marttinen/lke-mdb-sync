# Linode LKE MongoDB Sync

Sync the access control list of managed MongoDB clusters to the worker nodes of an LKE cluster.

This is a workaround for LKE cluster with autoscaling enabled, until managed MongoDB clusters support private IPs.

## Configuration

The current code is designed to be called via a cronjob (ie. via GitLab schedule).

| Env var | Required | Type | Purpose |
|---|---|---|---|
| DEBUG | no | bool | displayed some more information when calling Linode or K8s |
| LINODE_TOKEN | string | | Linde API Token with scope `databases:read_write` |
| INSTANCE_IDS | []int | | MongoDB instance IDs, comma separated eg. "1234,5678" |
| KUBECONFIG | string | | Absolute path to a K8s config file |

If `KUBECONFIG` is empty an in-cluster-authentication will be tried.

## Why?

The same can be achieved with:
```
addresses=$(kubectl get nodes -o jsonpath="{.items[*].status.addresses[?(@.type=='ExternalIP')].address}")
linode-cli databases mongodb-update --allow_list 192.168.128.0/17 --allow_list ${addresses// / --allow_list } $MDB_INSTANCE_ID
```

This will create an Linode auditlog entry everytime its called as the linode-cli doesn't support returning the allowList of a MongoDB cluster to avoid it.

## Todo

- code: structured logging
- code: tests
- code: run as long running process inside of K8s, watching node changes
- deploy: change to scratch image instead of alpine base
- deploy: provide Helm chart
- ci: provide release on GitHub
