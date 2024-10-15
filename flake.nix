{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:tweag/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs = {self, ...} @ inputs:
    inputs.flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import inputs.nixpkgs {
        inherit system;
        overlays = [inputs.gomod2nix.overlays.default];
      };
    in {
      packages = {
        default = pkgs.callPackage ./package.nix {};
      };
      apps = {
        default = {
          type = "app";
          program = "${self.packages."${system}".default}/bin/gothello";
        };
      };
      devShells = with pkgs; {
        default = mkShell {
          packages = [
            go
            gomod2nix
            gopls
          ];
          shellHook = ''
            ${pkgs.go}/bin/go mod tidy
          '';
        };
      };
    });
}
