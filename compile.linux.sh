#!/bin/bash

output_executable="gboard"

build_folder="build"

data_folder="$build_folder/data"
configs_folder="$data_folder/configs"
exports_folder="$data_folder/exports"

config_json="./.config.json"

base_config_json="./model.config.json"

go build -o "$build_folder/$output_executable"

mkdir -p "$data_folder"
mkdir -p "$configs_folder"
mkdir -p "$exports_folder"

cp "$config_json" "$build_folder/"
cp "$base_config_json" "$configs_folder/"

echo "gboard compiled."
