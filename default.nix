{pkgs ? import <nixpkgs> {}}: let
  lib = pkgs.lib;
  matchbox = pkgs.buildGoModule {
  name = "matchbox";
  src = lib.cleanSource ../matchbox;
  vendorHash = "sha256-sVC4xeQIcqAbKU4MOAtNicHcioYjdsleQwKWLstnjfk";
};
in matchbox

