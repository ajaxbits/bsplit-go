fmt:
    gofmt -w .

serve:
    go run ./cmd/.

update-modules:
    gomod2nix --outdir ./nix

clean:
    rm -fr ./expenses.db*
