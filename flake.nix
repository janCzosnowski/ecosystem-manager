{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs }: 
				let
								system = "x86_64-linux";
								pkgs = import nixpkgs {
												inherit system;
												config.allowUnfree = true;
								};
				in
				{
								packages.${system}.default = pkgs.buildGoModule {
												pname = "ecosystem-manager";
												version = "1.0.1";
												src = ./.;

												vendorHash = "sha256-m5mBubfbXXqXKsygF5j7cHEY+bXhAMcXUts5KBKoLzM=";
								};
				};
								
}
