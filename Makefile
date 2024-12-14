goal: bin/mygame.exe bin/mygame-arm bin/mygame-amd

bin/mygame.exe:
	GOOS=windows GOARCH=amd64 go build -o bin/mygame.exe .

bin/mygame-arm:
	GOOS=darwin GOARCH=arm64 go build -o bin/mygame-arm .

bin/mygame-amd:
	GOOS=darwin GOARCH=amd64 go build -o bin/mygame-amd .

clean:
	rm bin/*

