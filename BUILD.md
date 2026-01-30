
1. scripts/setup_submodules.sh
2. ./external/katran/build_bpf_modules_opensource.sh -s ./external/katran/
3.
```
  rm -rf _build/CMakeCache.txt _build/CMakeFiles
  cmake -B _build -DCMAKE_PREFIX_PATH="$(pwd)/_build/deps"
  make -C _build
```
4. copy_build_artifacts.sh
5. cd go/cmd/katran-server && go build
6. sudo ./go/cmd/katran-server/katran-server -static-dir ui/dist/ -bpf-prog-dir ./_build_go/ -config ./_build_go/config.yaml 
(-config is config exists; use go/config_example.yaml for config example)
