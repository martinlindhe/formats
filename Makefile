status:
	echo "| Format | Status |" > STATUS.md
	echo "| ------ | ------:|" >> STATUS.md
	grep -r STATUS parse/* | sed 's/parse\//| /' | sed 's/\/\/ STATUS:/ |/' | sed 's/%/% |/' | sed 's/://' | sort >> STATUS.md

install:
	go get -u golang.org/x/tools/cmd/stringer
	go generate ./...
	go get ./cmd/...
