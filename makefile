all:
	@termline grey
	GOOS=linux GOARCH=amd64 go build -o build/mp3cat-linux-amd64/mp3cat
	GOOS=darwin GOARCH=amd64 go build -o build/mp3cat-mac-amd64/mp3cat
	GOOS=darwin GOARCH=arm64 go build -o build/mp3cat-mac-arm64/mp3cat
	GOOS=windows GOARCH=amd64 go build -o build/mp3cat-windows-amd64/mp3cat.exe
	@termline grey
	@tree build
	@termline grey
	@shasum -a 256 build/*/*
	@termline grey
	@mkdir -p zipped
	@cd build && zip -r ../zipped/mp3cat-linux-amd64.zip mp3cat-linux-amd64 > /dev/null
	@cd build && zip -r ../zipped/mp3cat-mac-amd64.zip mp3cat-mac-amd64 > /dev/null
	@cd build && zip -r ../zipped/mp3cat-mac-arm64.zip mp3cat-mac-arm64 > /dev/null
	@cd build && zip -r ../zipped/mp3cat-windows-amd64.zip mp3cat-windows-amd64 > /dev/null
	@tree zipped
	@termline grey

clean:
	rm -rf ./build
	rm -rf ./zipped
