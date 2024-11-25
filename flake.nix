{
  description = "Go Event Streaming env";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    giopkg.url = "sourcehut://git.sr.ht/~eliasnaur/gio";
  };

  outputs = { self, nixpkgs, giopkg }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
        giopkgDevShell = giopkg.outputs.devShells.${system}.default;
      });
    in
    {
      devShells = forEachSupportedSystem ({ pkgs, giopkgDevShell }: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            gnumake
          ];

          inputsFrom = [
            giopkgDevShell
          ];
        };
      });
    };
}
