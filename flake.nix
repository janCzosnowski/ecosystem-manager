{
  description = "Interaction layer for my stuff";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }: let
				system = "x86_64-linux";
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };
      in {
        defaultPackage = pkgs.buildGoModule {
          pname = "ecosystem-manager";
          version = "1.0.3";
          src = ./.;
          vendorHash = "sha256-m5mBubfbXXqXKsygF5j7cHEY+bXhAMcXUts5KBKoLzM=";
        };
      };
    
}

