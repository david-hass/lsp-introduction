### FlowLang
Eine einfache (Proof-of-Concept-)Sprache zur Definition von Datenverarbeitungs-Pipelines.

Sie definiert, woher die Daten stammen (Quelle):
```hcl
source „raw_user_data“ {
    provider: file
    path: „/data/users.csv“
}
```

welche Transformationen (Task) darauf angewendet werden:
```hcl
task „filter_active_users“ {
    # Nimmt Output von „raw_user_data“ als Input
    input: „raw_user_data“
    transformer: „filter_by_column ‚status‘ == ‚active‘“
}
```


und wohin sie fließen (Sink):
```hcl
sink „active_user_report“ {
    # Nimmt Output von „filter_active_users“ als Input
    input: „filter_active_users“
    target: file
    path: „/reports/active_users.json“
}
```

Die Grammatik für einen LR-Parser muss eindeutig sein, wie im FlowLang-Beispiel, oder durch Regeln zur Auflösung von Prioritätskonflikten ergänzt werden.  
https://en.wikipedia.org/wiki/LR_parser#Conflicts_in_the_constructed_tables

