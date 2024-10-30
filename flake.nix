{
  description = "Bot to manage groceries";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    (flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        defaultPackage = pkgs.buildGo123Module {
          pname = "listedecourse";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-0GfassNv0886Z1D6PB9wdMkmRXfaaXeGJKDDgSJ5KJM=";
        };
      }
    ))
    // {
      checks = {
        listedecourse-aarch64-linux = self.defaultPackage.aarch64-linux;
      };
    };
}
