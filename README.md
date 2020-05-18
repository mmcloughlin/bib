# bib

BibTeX references for your Go source code.

* Organize your references in standard BibTeX format
* Cite references in source code with `[key]` syntax
* `bib` will insert a bibliography comment with details of the citations in each file
* BibTeX formatter
* Generate templated output: for example markdown biblography for documentation
* Link checker for bibliography URLs

## Install

```
go get -u github.com/mmcloughlin/bib
```

## Usage

Suppose you have some comments that need to refer to external sources, such
as the following [real example from the `crypto/ecdsa` standard library
package](https://github.com/golang/go/blob/go1.13.7/src/crypto/ecdsa/ecdsa.go).
Note the citations `[NSA]` and `[SECG]`:

[embedmd]:# (testdata/golden/ecdsa.in go /\/\/ hashToInt/ /orderBytes.+$/)
```go
// hashToInt converts a hash value to an integer. There is some disagreement
// about how this is done. [NSA] suggests that this is done in the obvious
// manner, but [SECG] truncates the hash to the bit-length of the curve order
// first. We follow [SECG] because that's what OpenSSL does. Additionally,
// OpenSSL right shifts excess bits from the number if the hash is too large
// and we mirror that too.
func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
	orderBits := c.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
```

Define these references in a [BibTeX](http://www.bibtex.org/) bibliography file.

[embedmd]:# (testdata/golden/ecdsa.bib)
```bib
@misc{NSA,
    title  = "Suite B Implementer’s Guide to FIPS 186-3 (ECDSA)",
    author = "NSA CSS",
    url    = "https://apps.nsa.gov/iaarchive/library/ia-guidance/ia-solutions-for-classified/algorithm-guidance/suite-b-implementers-guide-to-fips-186-3-ecdsa.cfm",
    year   = 2010,
}

@misc{SECG,
    title        = "SEC 1: Elliptic Curve Cryptography",
    author       = "Certicom Research",
    url          = "https://www.secg.org/sec1-v2.pdf",
    howpublished = "Standards for Efficient Cryptography 1",
    year         = 2009,
}
```

Include a comment starting `References:` in the source file where you would
like the bibliography to be inserted.

[embedmd]:# (testdata/golden/ecdsa.in go /\/\/ References:/ /References:$/)
```go
// References:
```

Now run `bib` on the file.

[embedmd]:# (testdata/scripts/basic.txt sh /bib process -w/ /source\.go/)
```sh
bib process -w -bib references.bib source.go
```

This will edit the file to insert a bibliography, as follows:

[embedmd]:# (testdata/golden/ecdsa.golden go /\/\/ References:/ /secg\.org.+$/)
```go
// References:
//
//	[NSA]   NSA CSS. Suite B Implementer’s Guide to FIPS 186-3 (ECDSA). 2010.
//	        https://apps.nsa.gov/iaarchive/library/ia-guidance/ia-solutions-for-classified/algorithm-guidance/suite-b-implementers-guide-to-fips-186-3-ecdsa.cfm
//	[SECG]  Certicom Research. SEC 1: Elliptic Curve Cryptography. Standards for Efficient
//	        Cryptography 1. 2009. https://www.secg.org/sec1-v2.pdf
```

## Additional Features

* Format BibTeX files with `bib fmt`
* Generate templated output with `bib generate`:
  - Markdown bibliography with `bib generate -type markdown`
  - Custom templates with `bib generate -tmpl <template>`
* Link check URLs in your bibliography with `bib linkcheck` command.

## License

`bib` is available under the [BSD 3-Clause License](LICENSE).
