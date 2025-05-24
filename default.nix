{pkgs ? import <nixpkgs> {}}:
pkgs.buildGoModule {
  name = "matchbox";
  src = pkgs.lib.cleanSource ../matchbox;
  vendorHash = "sha256-sVC4xeQIcqAbKU4MOAtNicHcioYjdsleQwKWLstnjfk";
}
