fmt:
    gofmt -w .

serve:
    go run .

update-modules:
    gomod2nix --outdir ./nix
