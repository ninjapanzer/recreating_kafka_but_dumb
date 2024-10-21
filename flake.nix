{
  description = "Go Event Streaming env";

  inputs.nixpkgs = {
    type = "github";
    owner = "ninjapanzer";
    repo = "nixpkgs";
    rev = "ee8a7accf156e3c1bde14f21f6acfc14ccec133f";
  };

  # inputs.nixpkgs.url = "/home/paulscoder/repos/nixpkgs";

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });
    in
    {
      devShells = forEachSupportedSystem ({ pkgs }: {
        default = pkgs.mkShell {
          packages = with pkgs; [
          ];
        };
      });
    };
}
