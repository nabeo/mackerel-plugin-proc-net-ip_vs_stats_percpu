mackerel-plugin-proc-net-ip_vs_stats_percpu
===========================================

Per-CPU IPVS statistics for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-proc-net-ip_vs_stats_percpu [-cpus=<cpu cores>] [-tempfile=<tempfile>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.ipvs]
command = "mackerel-plugin-proc-net-ip_vs_stats_percpu"
```
