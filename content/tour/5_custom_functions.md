---
title: 5. Custom Functions
description: >
  How to use a custom Scanner.
weight: 4
drivers: [modernc.org/sqlite]
scanners: [Scan, Custom]
executors: [Exec, First, One, All]
configs: [ParseFiles, Lookup, Masterminds/sprig, Funcs]
---

{{% pageinfo color="info" %}}
An alternative to this approach is to implement the sql.Scanner and sql.Valuer interfaces.
{{% /pageinfo %}}

{{< code language="go-template" source="tour/custom_functions/queries.go.tpl" >}}{{< /code >}}

{{< code language="go" source="tour/custom_functions/repository.go" >}}{{< /code >}}

<div style="padding-top: 2em; text-align: right"><a href="/tour/6_complex_query/">>> 6. Complex Query</a></div>
