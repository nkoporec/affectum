version := 0.0.0-SNAPSHOT # `make version=0.1.0 tag`
author  := Nejc Koporec
app     := Affectum
id      := com.nkoporec.affectum

install:
	go get -u ./... && go mod tidy

format:
	gofmt -s -w .

generate:
	cd icon/ && ./gen.sh
	make format

build-win:
	GOOS=windows GOARCH=amd64 go build -v -ldflags "-s -w -X main.version=$(version) -H=windowsgui" -o bin/win/affectum.exe ./

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -v -ldflags "-s -w -X main.version=$(version)" -o bin/darwin/affectum ./

build-linux:
	GOOS=linux GOARCH=amd64 go build -v -ldflags "-s -w -X main.version=$(version)" -o bin/linux/affectum ./

build:
	go build -v -o bin/affectum ./

clean:
	rm -rf bin/ dist/

bundle-darwin: build-darwin
	# Package solution to .app folder
	cd bin/darwin/ && \
		appify \
			-author "$(author)" \
			-id $(id) \
			-version $(version) \
			-name "$(app)" \
			-icon ../../icon/assets/base.png \
			./affectum
	/usr/libexec/PlistBuddy -c 'Add :LSUIElement bool true' 'bin/darwin/$(app).app/Contents/Info.plist'
	rm 'bin/darwin/$(app).app/Contents/README'
	# Package solution to .dmg image
	cd bin/darwin/ && \
		create-dmg --volname '$(app)' --app-drop-link 600 185 '$(app) $(version).dmg' ./'$(app).app' \
			|| true # ignore Error 2
	# Rename .dmg appropriately
	mv 'bin/darwin/$(app) $(version).dmg' bin/darwin/affectum_$(version)_darwin_x86_64.dmg
	# Remove temp files
	rm -rf 'bin/darwin/$(app).app'

bundle-linux:
	docker build . -t affectum
	docker run -v $$PWD/bin:/build/bin -t affectum make version=$(version) build-linux

tag:
	git tag -a v$(version) -m "Version $(version)"
	git push origin --tags

release-snapshot:
	goreleaser --rm-dist --skip-publish --snapshot
	cd dist && ls *.dmg | xargs shasum -a256 >> $$(ls *_checksums.txt)

release:
	goreleaser --rm-dist --skip-publish
	cd dist && ls *.dmg | xargs shasum -a256 >> $$(ls *_checksums.txt)

start: run # alias for run
run:
	pkill affectum || true
	nohup go run ./ >/dev/null 2>&1 &
