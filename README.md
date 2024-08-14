# Bench Maker

> [!IMPORTANT]
>
> This is proof of concept for the **Bench Maker**. At present, it does not
> create a working bench, just verifies the ideas.
>
> See the [POC](#poc) section.

Bench Maker is a performant `bench` replacement. Specifically it is meant to replace the `bench init` and `bench get-app` commands.

## Why

The main purpose of this is to be used by _Frappe Cloud_ to speed up builds.
Majority of the time spent in a build is spent in the App Install stages which
involves running `bench get-app` sequentially for a selection of _Frappe Apps_.

## How

The `bench init` and `bench get-app` command when utilized for more than one app, consist of mutually independent steps that can be run concurrently. This is better explained by use of a graph:

```mermaid
---
title: Bench Maker Execution Flow
---
graph TD
  %% Nodes

  %% Common
  BM("Begin making bench")
  IB("Initialize directory")
  CO("Complete")

  %% Multiprocessing
  W(("Wait"))
  C(("Concurrent"))
  S(("Sequential"))

  %% Fetch apps
  F_Fr("<code>frappe</code>: fetch repo")
  F_A1("<code>app_1</code>: fetch repo")
  F_AN("<code>app_n</code>: fetch repo")

  %% Validate apps
  V_Fr("<code>frappe</code>: validate")
  V_A1("<code>app_1</code>: validate")
  V_AN("<code>app_n</code>: validate")

  %% Install JS Dependencies
  I_Fr("<code>frappe</code>: install JS")
  I_A1("<code>app_1</code>: install JS")
  I_AN("<code>app_n</code>: install JS")

  %% Build JS
  B_Fr("<code>frappe</code>: build JS")
  B_A1("<code>app_1</code>: build JS")
  B_AN("<code>app_n</code>: build JS")

  %% Install Python Dependencies
  P_Fr("<code>frappe</code>: install Py")
  P_A1("<code>app_1</code>: install Py")
  P_AN("<code>app_n</code>: install Py")


  %% Styling

  classDef fetch fill:#f8dee2,stroke:#e68d99
  class F_Fr,F_A1,F_AN fetch

  classDef validate fill:#ffead6,stroke:#ffc187
  class V_Fr,V_A1,V_AN validate

  classDef installjs fill:#fff6d4,stroke:#ffd224
  class I_Fr,I_A1,I_AN installjs

  classDef buildjs fill:#e6f3ed,stroke:#96CEB4
  class B_Fr,B_A1,B_AN buildjs

  classDef installpy fill:#daedee,stroke:#97cdcf
  class P_Fr,P_A1,P_AN installpy

  style IB fill:#eedae3,stroke:#cf97b1

  classDef prc fill:#eee,stroke:#888
  class W,S,C prc
  style W stroke-dasharray: 5 2

  classDef ends fill:#EEE,stroke:#000,stroke-width:2px
  class BM,CO ends


  %% Chart

  BM --> IB
  IB --> C

  C --> F_Fr
  C --> F_A1
  C -.-> F_AN

  F_Fr --> V_Fr
  F_A1 --> V_A1
  F_AN -.-> V_AN

  V_Fr --> I_Fr
  V_A1 --> I_A1
  V_AN -.-> I_AN

  I_Fr --> B_Fr
  I_A1 --> B_A1
  I_AN -.-> B_AN

  B_Fr --> W
  B_A1 --> W
  B_AN -.-> W

  W --> S

  S --> P_Fr
  P_Fr --> P_A1
  P_A1 -.-> P_AN

  P_AN --> CO

```

4 out of the 5 app stages involved in making a bench can be run concurrently.

> [!NOTE]
>
> Installing Python dependencies have to be run sequentially because all apps on
> a _Frappe Bench_ share the same python environment.

## POC

This is as of now a proof of concept. It may or may not be fleshed out. The
ideas and implementations I wanted to test out and have been verified are:

- Concurrent installation of _Frappe Apps_ being possible.
- Concurrent installation of _Frappe Apps_ taking much lesser time than sequential installation.
- Multiplexing of output from concurrent installs.
- Being able to cleanly stop execution if any app install fails.

Few things I have not yet tested out are:

- Building a working bench using bench maker
- Performance impact of multiple instances of bench maker running separately.
- Speed up from using alternative package managers than `yarn` or `pip`.
- Speed up from caching different stages. As of now only the fetch app stage is
  non optimally cached, other than that `yarn` and `pip` use their own caches.

### Run

> [!NOTE]
>
> To build this you'll need Go > 1.22.4.

```bash
# Build the Bench Maker binary `bm`
go build

# Example
./bm --apps erpnext hrms gameplan builder

# Example without cache
./bm --no-cache --apps erpnext hrms gameplan builder

# Example sequential install
./bm --seq --apps erpnext hrms gameplan builder
```

This will create a `temp` folder:

- `temp/bench`: dummy bench where the apps are installed.
- `temp/.cache`: where the repos, `yarn` and `pip` cache are.

## Time taken

1. Empty cache, successfully built.

```shell
# bm --apps erpnext hrms drive builder

Time Breakdown:
| org/repo         |     clone |  validate |    ins js |     build |    ins py |     total |
|------------------|-----------|-----------|-----------|-----------|-----------|-----------|
| frappe/drive     |    9.284s |    0.000s |   37.399s |   12.317s |    2.777s |   61.777s |
| frappe/hrms      |   11.015s |    0.000s |  435.528s |    6.674s |    2.732s |  455.950s |
| frappe/erpnext   |   12.659s |    0.000s |    0.908s |    0.000s |    2.825s |   16.392s |
| frappe/frappe    |   13.112s |    0.000s |   71.343s |    0.922s |    4.424s |   89.802s |
| frappe/builder   |    1.660s |    0.000s |  299.026s |    7.252s |    2.705s |  310.644s |

Totals:
Bench init            :    2.028s
Concurrent app stages :  919.101s
Sequential app stages :   15.464s
---------------------------------
Total app             :  934.565s
Total app + bench     :  936.593s
---------------------------------
Total wall time       :  470.710s
Time saved            :  465.883s
```

2. Non empty cache, successfully built, network issue caused slowdown:

```shell
# bm --apps erpnext hrms drive builder

Time Breakdown:
| org/repo         |     clone |  validate |    ins js |     build |    ins py |     total |
|------------------|-----------|-----------|-----------|-----------|-----------|-----------|
| frappe/builder   |    0.098s |    0.000s |   13.320s |    8.037s |    2.665s |   24.121s |
| frappe/drive     |    0.161s |    0.000s |   43.473s |   12.926s |    2.764s |   59.324s |
| frappe/hrms      |    0.405s |    0.000s |  300.492s |    6.574s |    2.661s |  310.131s |
| frappe/frappe    |    0.855s |    0.000s |   11.198s |    1.059s |    4.173s |   17.285s |
| frappe/erpnext   |    1.163s |    0.000s |    0.211s |    0.000s |    2.811s |    4.185s |

Totals:
Bench init            :    1.998s
Concurrent app stages :  399.970s
Sequential app stages :   15.075s
---------------------------------
Total app             :  415.045s
Total app + bench     :  417.044s
---------------------------------
Total wall time       :  324.544s
Time saved            :   92.500s
```

```shell
# When installing JS dependencies for frappe/hrms
info There appears to be trouble with your network connection. Retrying...
```

## Glossary

A glossary has been included cause due to asinine naming, the term "bench" is
terribly overloaded. In the context of FC, it refers to at least 4 different
things.

- **_Frappe Bench_**: A collection of _Frappe Apps_ managed by `bench`.
- **`bench`**: [Tool](https://github.com/frappe/bench) used to manage _Frappe Benches_.
- **_Frappe App_**: A web-app built using FF.
- **_Frappe Cloud_**: [Platform](https://frappecloud.com/) that hosts _Frappe Benches_.
- **BM**: Bench Maker.
- **FF**: [Frappe Framework](https://github.com/frappe/frappe).
- **FC**: Frappe Cloud.
