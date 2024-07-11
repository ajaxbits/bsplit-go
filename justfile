fmt:
    gofmt -w .

serve:
    go run .

update-modules:
    gomod2nix --outdir ./nix
    direnv deny
    rm -fr .direnv
    direnv allow

clean:
    rm -fr ./expenses.db*

gen:
    sqlc generate