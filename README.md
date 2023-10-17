# Instructions

1. Start the server:

```sh
go run $(pwd)/cmd/server
```

2. Upload image files:

```sh
for f in $(ls image/*)
do
  curl -X POST --data-binary @$f localhost:8080/cats
  echo
done
```

3. Store the last printed checksum in your env (it always prints the bmt root sum):

```sh
YOUR_ROOT_HASH=c1ea98387203ae52b4f70d9ee1b9ac0b02c9da5247ea293d8786c7234bf53b20
```

4. Request proofs for a given image (2):

```sh
curl -XGET -s localhost:8080/proof/cats/2 > proof.json
```

4. Validate file integrity:

```sh
go run $(pwd)/cmd/checker -root $YOUR_ROOT_HASH -proofs $(<proof.json) image/cat471.jpg
```

Checksums should match for this file, and should not match for any other.
