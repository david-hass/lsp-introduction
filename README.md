[![en](https://img.shields.io/badge/lang-en-red.svg)](https://github.com/david-hass/lsp-introduction/blob/main/README.md)
[![de](https://img.shields.io/badge/lang-de-yellow.svg)](https://github.com/david-hass/lsp-introduction/blob/main/README.de.md)


### FlowLang
A simple (proof of concept) language for defining data processing pipelines.

It defines where data comes from (source):
```hcl
source "raw_user_data" {
    file: "/data/users.csv"
    encoding: "utf-8"
}
```

what transformations (task) are applied to it:
```hcl
task "filter_active_users" {
    # Nimmt Output von "raw_user_data" als Input
    input: raw_user_data
    transformer: "filter_by_column 'status' == 'active'"
}
```


and where it flows to (sink):
```hcl
sink "active_user_report" {
    # Nimmt den Output des Anonymisierungs-Tasks
    input: filter_active_users
    path: "/reports/active_users.json"
}
```

The grammar for an LR parser must be unambiguous, as is the case in the FlowLang example, or must be augmented by tie-breaking precedence rules.  
https://en.wikipedia.org/wiki/LR_parser#Conflicts_in_the_constructed_tables

