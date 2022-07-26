# Morpheus Terraformer

[![GoReportCard][report-badge]][report]
[![GitHub release](https://img.shields.io/github/release/martezr/morpheus-terraformer.svg)](https://github.com/martezr/morpheus-terraformer/releases/)
[![GoDoc](https://pkg.go.dev/badge/badge/github.com/martezr/morpheus-terraformer?utm_source=godoc)](https://godoc.org/github.com/martezr/morpheus-terraformer)
![GitHub](https://img.shields.io/github/license/martezr/morpheus-terraformer)


[📖 Documentation][docs]

[docs]: https://martezr.github.io/morpheus-terraformer
[report-badge]: https://goreportcard.com/badge/github.com/martezr/morpheus-terraformer
[report]: https://goreportcard.com/report/github.com/martezr/morpheus-terraformer

A CLI tool that generates `tf` files based on existing Morpheus configuration
(reverse Terraform).

*   Disclaimer: This is not an official Morpheus product

# Installation

Morpheus-Terraformer is built with Golang and compiled down to a single binary like Terraform.

## Linux

```bash
curl -LO https://github.com/martezr/morpheus-terraformer/releases/download/$(curl -s https://api.github.com/repos/martezr/morpheus-terraformer/releases/latest | grep tag_name | cut -d '"' -f 4)/morpheus-terraformer-linux-amd64.tar.gz
tar -xzf morpheus-terraformer-linux-amd64.tar.gz
cd morpheus-terraformer-linux-amd64/
chmod +x morpheus-terraformer-linux-amd64
sudo mv morpheus-terraformer-linux-amd64 /usr/local/bin/morpheus-terraformer
```

## Mac OS X

```
curl -LO https://github.com/martezr/morpheus-terraformer/releases/download/$(curl -s https://api.github.com/repos/martezr/morpheus-terraformer/releases/latest | grep tag_name | cut -d '"' -f 4)/morpheus-terraformer-darwin-amd64.tar.gz
tar -xzf morpheus-terraformer-darwin-amd64.tar.gz
cd morpheus-terraformer-darwin-amd64/
chmod +x morpheus-terraformer-darwin-amd64
sudo mv morpheus-terraformer-darwin-amd64 /usr/local/bin/morpheus-terraformer
```

## Windows

1. Download the .exe file from here - https://github.com/martezr/morpheus-terraformer/releases
2. Add the exe file path to path variable

# Getting Started

To get started with the Morpheus Terraformer the connection information for the Morpheus instance need to be set as environment variables.

**Username and Password**

```
export MORPHEUS_API_URL="https://morpheus.test.local"
export MORPHEUS_API_USERNAME="admin"
export MORPHEUS_API_PASSWORD="password123"
```

**Access Token**

```
export MORPHEUS_API_URL="https://morpheus.test.local"
export MORPHEUS_API_TOKEN="d3a4c6fa-fb54-44af"
```

Generate the Terraform code using the `generate` command following by `-r` and the name of the resources to create or specify `'*'` to generate all supported resources.

```
morpheus-terraformer generate -r "*"
```

```
morpheus-terraformer generate -r group,environment
```

# Supported Resources

The following resources are supported:

|Name|Description|
|----|-----|
|contact|Morpheus monitoring contact|
|environment|Morpheus environments|
|group|Morpheus groups|
|optionlist|Morpheus option lists (REST API, Manual, and Morpheus)|
|optiontype|Morpheus option types(select, text, number, password, typeahead)|
|policy|Morpheus policy|
|scripttemplate|Morpheus script template|
|spectemplate|Morpheus spec template|
|task|Morpheus automation task|
|tenant|Morpheus tenants|
|wiki|Morpheus wiki pages|
|workflow|Morpheus automation workflows|

# Contributing

If you have improvements or fixes, we would love to have your contributions. Please read CONTRIBUTING.md for more information on the process we would like contributors to follow.