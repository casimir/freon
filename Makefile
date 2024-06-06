VERSION := $(shell ./tools/extract_version.sh)

all:: ui server
	
clean:: ui-clean server-clean

ui::
	cd ui && \
	flutter build web --base-href /ui/
	
ui-clean::
	cd ui && \
	flutter clean

server:: server-headless server-full
	
server-clean::
	rm -rf ./server/build/
	rm -rf ./server/ui/statics/
	
server-headless::
	mkdir -p server/build
	cd server && \
	CGO_ENABLED=1 go build -o build/freon-headless -ldflags="-X github.com/casimir/freon/buildinfo.Version=${VERSION}"

server-full:: ui
	mkdir -p server/build
	cp -r ui/build/web server/ui/statics
	cd server && \
	CGO_ENABLED=1 go build -o build/freon -tags embed -ldflags="-X github.com/casimir/freon/buildinfo.Version=${VERSION}"
	