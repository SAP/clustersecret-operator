---
title: "Webhook startup options"
linkTitle: "Webhook startup options"
weight: 20
type: "docs"
description: >
  Command line configuration options, environment variables
---

## Command line parameters

This repository ships two executables, controller and webhook.
The webhook accepts the following command line flags:

```bash
Usage of ./go/bin/webhook:
      --bind_address string              Bind address (default ":1080")
      --tls_enabled                      Enable TlS
      --tls_key_file string              Path to TLS key
      --tls_cert_file string             Path to TLS certificate
      --add_dir_header                   If true, adds the file directory to the header of the log messages
      --alsologtostderr                  log to standard error as well as files
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --log_file string                  If non-empty, use this log file
      --log_file_max_size uint           Defines the maximum size a log file can grow to. Unit is megabytes.
                                         If the value is 0, the maximum file size is unlimited. (default 1800)
      --logtostderr                      log to standard error instead of files (default true)
      --one_output                       If true, only write logs to their native severity level
                                         (vs also writing to each lower severity level)
      --skip_headers                     If true, avoid header prefixes in the log messages
      --skip_log_headers                 If true, avoid headers when opening log files
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          number for the log level verbosity
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

## Logging

The webhook uses [klog v2](https://github.com/kubernetes/klog) for logging.
Please check the according documentation for details about how to configure logging.
