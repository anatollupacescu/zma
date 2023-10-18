# Instructions

1. Start the server:

```sh
go run $(pwd)/cmd/server
```

Alternatively, in docker:

```sh
docker build . -t zma
docker run -p 8080:8080 zma
```

2. Upload image files to an arbitrary collection (cats):

```sh
for f in $(ls image/*)
do
  curl -X POST --data-binary @$f localhost:8080/cats
  echo
done
```

3. Store the last printed checksum in your env (it always prints the bmt root sum):

```sh
YOUR_ROOT_HASH=14e08396215c156ff4998eca119fc15173d544bc74fad48332b2012bcb226774
```

At this point it is safe to remove your local images.

You can download again the server stored file using:

```sh
wget localhost:8080/cats/2 -O my.cat
```

Now you can proceed to verify file contents integrity:

1. Request proofs for a given image (2):

```sh
curl -XGET -s localhost:8080/proof/cats/2 > proof.json
```

2. Validate file integrity aginst your root checksum and downloaded contents:

```sh
go run $(pwd)/cmd/checker -root $YOUR_ROOT_HASH -proofs $(<proof.json) my.cat
```

Checksums should match for this file but should not match for any other.
