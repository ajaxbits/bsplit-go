fmt:
    gofmt -w .

serve:
    go run .

update-modules:
    gomod2nix --outdir ./nix

clean:
    rm -fr ./expenses.db*
