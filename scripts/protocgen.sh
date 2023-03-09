
set -eo pipefail

protoc_install_proto_gen_doc() {
  echo "Installing protobuf protoc-gen-doc plugin"
  (go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest 2> /dev/null)
}

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find ./cosmwasm ./iov -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep "option go_package" $file &> /dev/null ; then
      buf generate --template buf.gen.gogo.yaml $file
    fi
  done
done

protoc_install_proto_gen_doc

echo "Generating proto docs"
buf generate --template buf.gen.doc.yaml

cd ..

# move proto files to the right places
<<<<<<< HEAD
=======
cp -r github.com/CosmWasm/wasmd/* ./
>>>>>>> tags/v0.11.6
cp -r github.com/iov-one/starnamed/* ./
rm -rf github.com

go mod tidy