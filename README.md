# Morpheus Terraformer

[![GoReportCard][report-badge]][report]
[![GitHub release](https://img.shields.io/github/release/martezr/morpheus-terraformer.svg)](https://github.com/martezr/morpheus-terraformer/releases/)
[![GoDoc](https://pkg.go.dev/badge/badge/github.com/martezr/morpheus-terraformer?utm_source=godoc)](https://godoc.org/github.com/martezr/morpheus-terraformer)

[report-badge]: https://goreportcard.com/badge/github.com/martezr/morpheus-terraformer
[report]: https://goreportcard.com/report/github.com/martezr/morpheus-terraformer

A CLI tool that generates `tf` files based on existing Morpheus configuration
(reverse Terraform).

*   Disclaimer: This is not an official Morpheus product

# Installation

Morpheus-Terraformer is built with Golang and compiled down to a single binary like Terraform.

## Linux

```bash
curl -LO https://github.com/martezr/morpheus-terraformer/releases/download/$(curl -s https://api.github.com/repos/martezr/morpheus-terraformer/releases/latest | grep tag_name | cut -d '"' -f 4)/morpheus-terraformer-linux-amd64
chmod +x morpheus-terraformer-linux-amd64
sudo mv morpheus-terraformer-linux-amd64 /usr/local/bin/morpheus-terraformer
```

## Mac OS X

```
curl -LO https://github.com/martezr/morpheus-terraformer/releases/download/$(curl -s https://api.github.com/repos/martezr/morpheus-terraformer/releases/latest | grep tag_name | cut -d '"' -f 4)/morpheus-terraformer-darwin-amd64
chmod +x morpheus-terraformer-darwin-amd64
sudo mv morpheus-terraformer-darwin-amd64 /usr/local/bin/morpheus-terraformer
```

## Windows



# Getting Started

```
export MORPHEUS_API_URL="https://morpheus.test.local"
export MORPHEUS_API_USERNAME="admin"
export MORPHEUS_API_PASSWORD="password123"
```

Generate the Terraform code using the `generate` command following by `-r` and the name of the resources to create or specify `'*'` to generate all supported resources.

```
morpheus-terraformer generate -r groups,environments
```

# Supported Resources

The following resources are supported:

|Name|Description|
|----|-----|
|environment|Morpheus environments|
|group|Morpheus groups|
|option list|Morpheus option lists (REST API, Manual, and Morpheus)|
|option type|Morpheus option types(select, text, number, password, typeahead)|
|tenant|Morpheus tenants|

# Contributing

If you have improvements or fixes, we would love to have your contributions. Please read CONTRIBUTING.md for more information on the process we would like contributors to follow.