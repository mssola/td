test :: unit_test checks
checks :: vet fmt lint errcheck climate

vet ::
	@echo "+ $@"
		@go vet $$(go list ./... | grep -v vendor)

fmt ::
	@echo "+ $@"
		@test -z "$$(gofmt -s -l . | grep -v vendor | tee /dev/stderr)"

lint ::
	@echo "+ $@"
		@test -z "$$(golint ./... | grep -v vendor | tee /dev/stderr)"

climate ::
	@echo "+ $@"
		@(./scripts/climate -o -a lib)

errcheck ::
	@echo "+ $@"
		@test -z "$$(errcheck github.com/mssola/td/lib | tee /dev/stderr)"

unit_test ::
	@echo "+ $@"
		@go test -v $$(go list ./... | grep -v vendor | grep -v integration)
