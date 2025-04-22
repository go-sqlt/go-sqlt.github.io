---
title: 3. Bulk insert
description: >
  This example shows how to create statements for bulk inserts.
weight: 3
drivers: [modernc.org/sqlite]
scanners: [ScanInt, ScanString]
executors: [Exec, First, All]
configs: [ParseFiles]
---


{{< code language="go-template" source="tour/bulk_insert/queries.go.tpl" >}}{{< /code >}}

{{< code language="go" source="tour/bulk_insert/repository.go" >}}{{< /code >}}

<div style="padding-top: 2em; text-align: center"><a href="/tour/4_transactions/">>> 4. Transactions</a></div>
