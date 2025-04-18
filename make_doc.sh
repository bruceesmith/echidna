#!/bin/bash
echo '[![Go Reference][goreference_badge]][goreference_link]' > temp1
echo '[![Go Report Card][goreportcard_badge]][goreportcard_link]' >> temp1
echo " " >> temp1
echo '# <div align="center">Echidna Library for Go</div>' >> temp1
echo " " >> temp1
echo '![background image](echidna.png)' >> temp1
echo " " >> temp1
echo " " > temp2
echo '[goreference_badge]: https://pkg.go.dev/badge/github.com/bruceesmith/echidna/v3.svg' >> temp2
echo '[goreference_link]: https://pkg.go.dev/github.com/bruceesmith/echidna' >> temp2
echo '[goreportcard_badge]: https://goreportcard.com/badge/github.com/bruceesmith/echidna' >> temp2
echo '[goreportcard_link]: https://goreportcard.com/report/github.com/bruceesmith/echidna' >> temp2
go tool gomarkdoc ./... --output read
cat temp1 read temp2 > README.md
rm temp1 temp2 read
