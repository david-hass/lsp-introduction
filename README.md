### FlowLang
A simple (proof of concept) language for defining data processing pipelines.

It defines where data comes from (source):
```hcl
source "raw_user_data" {
    provider: file
    path: "/data/users.csv"
}
```

what transformations (task) are applied to it:
```hcl
task "filter_active_users" {
    # Nimmt Output von "raw_user_data" als Input
    input: "raw_user_data"
    transformer: "filter_by_column 'status' == 'active'"
}
```


and where it flows to (sink):
```hcl
sink "active_user_report" {
    # Nimmt Output von "filter_active_users" als Input
    input: "filter_active_users"
    target: file
    path: "/reports/active_users.json"
}
```
