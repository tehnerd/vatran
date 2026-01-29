


./external/katran/build_bpf_modules_opensource.sh -s ./external/katran/

  rm -rf _build/CMakeCache.txt _build/CMakeFiles
  cmake -B _build -DCMAKE_PREFIX_PATH="$(pwd)/_build/deps"
  make -C _build
