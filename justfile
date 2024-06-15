fmt:
    gofmt -w .

serve:
    go run ./cmd/.

update-modules:
    gomod2nix --outdir ./nix
