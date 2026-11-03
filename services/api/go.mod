module github.com/Daddy-senpaii/Vajra/services/api

go 1.26.1

require github.com/Daddy-senpaii/Vajra/packages/blocking-engine v0.0.0

require (
	golang.org/x/net v0.62.0 // indirect
	golang.org/x/text v0.43.0 // indirect
)

replace github.com/Daddy-senpaii/Vajra/packages/blocking-engine => ../../packages/blocking-engine
