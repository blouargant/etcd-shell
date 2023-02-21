all: dirs app completion

.PHONY: all clean app install uninstall

app:
	go build -o app/etcd-shell .
	chmod 774 app/etcd-shell

completion:
	go run main.go completion bash > app/completion.bash
	go run main.go completion zsh > app/completion.zsh
	go run main.go completion fish > app/completion.fish
	go run main.go completion powershell > app/completion.powershell

dirs:
	mkdir -p app

install:
	sudo cp app/etcd-shell /usr/local/bin/etcd-shell
	sudo chmod 755 /usr/local/bin/etcd-shell
	sudo cp ./app/completion.bash /usr/share/bash-completion/completions/etcd-shell
	sudo chmod 644 /usr/share/bash-completion/completions/etcd-shell

uninstall:
	sudo rm /usr/local/bin/etcd-shell
	sudo rm /usr/share/bash-completion/completions/etcd-shell

clean:
	rm -rf app
