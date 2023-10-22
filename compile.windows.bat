@echo off

set output_executable=gboard.exe
set build_folder=build

set data_folder=%build_folder%\data
set configs_folder=%data_folder%\configs
set exports_folder=%data_folder%\exports

set config_json=.\.config.json
set base_config_json=.\model.config.json

go build -o %build_folder%\%output_executable%

if not exist %data_folder% mkdir %data_folder%
if not exist %configs_folder% mkdir %configs_folder%
if not exist %exports_folder% mkdir %exports_folder%

copy %config_json% %build_folder%
copy %base_config_json% %configs_folder%

echo gboard compiled.