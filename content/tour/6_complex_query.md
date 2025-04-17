---
title: 6. Complex Query
description: >
  This example shows how to build complex queries.
weight: 5
drivers: [modernc.org/sqlite]
scanners: [Scan, ScanStringSlice]
executors: [Exec, All, One]
configs: [ParseFiles, Lookup, Masterminds/sprig, Log, NoExpirationCache]
---


{{< code language="go-template" source="tour/complex_query/queries.go.tpl" >}}{{< /code >}}

{{< code language="go" source="tour/complex_query/repository.go" >}}{{< /code >}}

<div style="padding-top: 2em; text-align: right"><a href="/tour/7_multiple_databases/">>> 7. Multiple Databases</a></div>
