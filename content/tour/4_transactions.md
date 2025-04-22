---
title: 4. Transactions
description: >
  This example shows how to use statements in transactions.
weight: 4
drivers: [modernc.org/sqlite]
scanners: [ScanInt, ScanString, ScanTime]
executors: [Exec, First, One, All]
configs: [ParseFiles, Lookup, Masterminds/sprig]
---

{{< code language="go-template" source="tour/transactions/queries.go.tpl" >}}{{< /code >}}

{{< code language="go" source="tour/transactions/repository.go" >}}{{< /code >}}

<div style="padding-top: 2em; text-align: center"><a href="/tour/5_custom_functions/">>> 5. Custom Functions</a></div>
