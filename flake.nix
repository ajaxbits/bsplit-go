{
  description = "bsplit: Split expenses, freely";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
    templ = {
      url = "github:a-h/templ";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.gomod2nix.follows = "gomod2nix";
    };
  };

  outputs = {
    self, # deadnix: skip
    nixpkgs,
    flake-utils,
    gomod2nix,
    templ,
  }:
    (
      flake-utils.lib.eachDefaultSystem
      (system: let
        pkgs = nixpkgs.legacyPackages.${system};
        templBin = templ.packages.${system}.templ;

        # The current default sdk for macOS fails to compile go projects, so we use a newer one for now.
        # This has no effect on other platforms.
        callPackage = pkgs.darwin.apple_sdk_11_0.callPackage or pkgs.callPackage;
      in {
        packages.default = callPackage ./nix/. {
          inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
          templ = templBin;
        };
        devShells.default = callPackage ./nix/shell.nix {
          inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
          templ = templBin;
        };
      })
    )
    // {
      nixConfig = {
        extra-substituters = [
          "https://cache.garnix.io"
        ];
        extra-trusted-public-keys = [
          "cache.garnix.io:CTFPyKSLcx5RMJKfLo5EEPUObbA78b0YQ2DTCJXqr9g="
        ];
      };
    };
}
