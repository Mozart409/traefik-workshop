# https://just.systems

set dotenv-load
set unstable

default:
    just --list

clear:
    clear

flake-update: clear
    nix flake update
