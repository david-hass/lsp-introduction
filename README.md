### FlowLang
A simple (proof of concept) language for defining data processing pipelines.

It defines where data comes from (source):
<code>
source <font color="red">"raw_user_data"</font> {
    provider: file
    path: "/data/users.csv"
}
</code>


what transformations (task) are applied to it:
<code>
task <font color="red">"filter_active_users"</font> {
    # Nimmt Output von "raw_user_data" als Input
    input: "raw_user_data"
    transformer: "filter_by_column 'status' == 'active'"
}
</code>


and where it flows to (sink):
<code>
sink <font color="red">"active_user_report"</font> {
    # Nimmt den Output des Anonymisierungs-Tasks
    input: "anonymize_pii"
    target: file
    path: "/reports/active_users.json"
}
</code>
