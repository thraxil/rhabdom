all: rhabdom

rhabdom: *.go
	go build .

run: rhabdom
	./rhabdom

clean:
	rm -f rhabdom

fmt:
	go fmt *.go

install_deps:
	go get -u github.com/mrb/riakpbc
	go get -u github.com/stvp/go-toml-config
	go get -u github.com/nu7hatch/gouuid
	go get -u github.com/thraxil/paginate
	go get -u github.com/gorilla/feeds
