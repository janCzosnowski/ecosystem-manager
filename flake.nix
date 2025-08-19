{
  description = "Interaction layer for my stuff";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }: {
    packages = builtins.genAttrs nixpkgs.lib.systems.flakeExposed (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };
      in {
        default = pkgs.buildGoModule {
          pname = "ecosystem-manager";
          version = "1.0.3";
          src = ./.;
          vendorHash = "sha256-m5mBubfbXXqXKsygF5j7cHEY+bXhAMcXUts5KBKoLzM=";
        };
      }
    );
  };
}

