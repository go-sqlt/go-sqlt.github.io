---
title: 7. Multiple Databases
description: >
  This example shows how to use sqlt with multiple databases.
weight: 5
drivers: [modernc.org/sqlite, jackc/pgx]
scanners: [Scan, ScanStringSlice]
executors: [Exec, All, One]
configs: [Configure, ParseFiles, Lookup, Masterminds/sprig, Log, NoExpirationCache]
---

{{< code language="go-template" source="tour/multiple_databases/queries.go.tpl" >}}{{< /code >}}

{{< code language="go" source="tour/multiple_databases/repository.go" >}}{{< /code >}}
