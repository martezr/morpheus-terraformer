module github.com/martezr/morpheus-terraformer

go 1.17

require (
	//	github.com/gomorpheus/morpheus-go-sdk v0.2.3
	github.com/hashicorp/hcl/v2 v2.11.1
	github.com/hashicorp/hcl2 v0.0.0-20191002203319-fb75b3253c80
	github.com/spf13/cobra v1.2.1
	github.com/zclconf/go-cty v1.10.0
)

replace github.com/gomorpheus/morpheus-go-sdk => ../morpheus-go-sdk

require github.com/gomorpheus/morpheus-go-sdk v0.0.0-00010101000000-000000000000

require (
	github.com/agext/levenshtein v1.2.1 // indirect
	github.com/apparentlymart/go-textseg v1.0.0 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/go-resty/resty/v2 v2.2.0 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mitchellh/go-wordwrap v0.0.0-20150314170334-ad45545899c7 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.0.0-20211111083644-e5c967477495 // indirect
	golang.org/x/text v0.3.6 // indirect
)
