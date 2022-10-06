@echo off

@set wsl_path=

for /F %%i in ('wsl wslpath %1') do @set "wsl_path=%%i

start "" "wt.exe" "--maximized" "-d"  %1
start "" "neovide.exe" "--maximized" "--wsl" %wsl_path%
