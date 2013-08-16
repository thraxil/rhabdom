all: rhabdom

rhabdom: rhabdom.go post.go index.go
	go build .

run: rhabdom
	./rhabdom

clean:
	rm -f rhabdom

fmt:
	go fmt *.go

install_deps:
	go get github.com/mrb/riakpbc
	go get github.com/stvp/go-toml-config
