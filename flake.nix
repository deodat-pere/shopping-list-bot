{
  description = "Bot to manage groceries";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    flake-utils.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
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
    )
    // (
      let
        pkgsCrossAarch64 = import nixpkgs {
          system = "x86_64-linux";
          crossSystem.config = "aarch64-unknown-linux-gnu";
        };
      in
      {
        packages.cross.listedecourse = pkgsCrossAarch64.buildGo123Module {
          pname = "listedecourse";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-0GfassNv0886Z1D6PB9wdMkmRXfaaXeGJKDDgSJ5KJM=";
        };
      }
    );
}
