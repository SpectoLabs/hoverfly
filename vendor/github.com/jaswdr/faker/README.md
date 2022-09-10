<p align="center">
  <img width="500" src="./cover.png">
</p>

Faker is a Go library that generates fake data for you. Whether you need to bootstrap your database, create good-looking XML documents, fill-in your persistence to stress test it, or anonymize data taken from a production service, Faker is for you.

Faker is heavily inspired by PHP"s [Faker](https://github.com/fzaninotto/Faker)

Faker requires Go >= 1.11

[![PkgGoDev](https://pkg.go.dev/badge/github.com/jaswdr/faker)](https://pkg.go.dev/github.com/jaswdr/faker)
[![Test](https://github.com/jaswdr/faker/actions/workflows/test.yml/badge.svg)](https://github.com/jaswdr/faker/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/jaswdr/faker/branch/master/graph/badge.svg)](https://codecov.io/gh/jaswdr/faker)
[![Go Report Card](https://goreportcard.com/badge/github.com/jaswdr/faker)](https://goreportcard.com/report/github.com/jaswdr/faker)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/ba14f84a3f824410be0a6f6670de012a)](https://app.codacy.com/gh/jaswdr/faker?utm_source=github.com&utm_medium=referral&utm_content=jaswdr/faker&utm_campaign=Badge_Grade)
[![CodeFactor](https://www.codefactor.io/repository/github/jaswdr/faker/badge)](https://www.codefactor.io/repository/github/jaswdr/faker)
[![Release](https://img.shields.io/github/release/jaswdr/faker.svg?style=flat-square)](https://github.com/jaswdr/faker/releases)
[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/jaswdr/faker)
[![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fjaswdr%2Ffaker&count_bg=%2379C83D&title_bg=%23555555&icon=&icon_color=%23E7E7E7&title=hits&edge_flat=false)](https://hits.seeyoufarm.com)

## Test it directly from your browser

Start at https://play.golang.org/p/JpTagDGBaHK

## Installation

Add this to your Go file

```go
import "github.com/jaswdr/faker"
```

And run `go get` or `dep ensure` to get the package.

## Basic Usage

Use `faker.New()` to create and initialize a faker generator, which can generate data by accessing properties named after the type of data you want.

```go
import "github.com/jaswdr/faker"

func main() {
    faker := faker.New()

    faker.Person().Name()
    // Lucy Cechtelar

    faker.Address().Address()
    // 426 Jordy Lodge

    faker.Lorem().Text(100)
    // Dolores sit sint laboriosam dolorem culpa et autem. Beatae nam sunt fugit
    // et sit et mollitia sed.
    // Fuga deserunt tempora facere magni omnis. Omnis quia temporibus laudantium
    // sit minima sint.
}
```

Even if this example shows a method access, each call to `faker.Name()` yields a different (random) result.

```go
p := faker.Person()

for i:=0; i < 10; i++ {
  fmt.Println(p.Name())
}
  // Adaline Reichel
  // Dr. Santa Prosacco DVM
  // Noemy Vandervort V
  // Lexi O"Conner
  // Gracie Weber
  // Roscoe Johns
  // Emmett Lebsack
  // Keegan Thiel
  // Wellington Koelpin II
  // Ms. Karley Kiehn V
```

You can also generate a profile image.

```go
image := p.Image()

fmt.Println(image.Name())
// /tmp/profil-picture-img-1064677774.jfif

fmt.Printf("%+v", image)
// &{file:0xc0002e4300}
```

Generate fake data using Structs

```go
type ExampleStruct struct {
	SimpleStringField string
	SimpleNumber int
	SimpleBool bool
	SomeFormatedString string `fake:"??? ###"`
	SomeStringArray [5]string `fake:"????"`
}

example := ExampleStruct{}
f.Struct().Fill(&example)
fmt.Printf("%+v", example)
//{SimpleStringField:87576a01c2a547b2bbf9b7c736d1db40 SimpleNumber:9223372036854775807 SimpleBool:false SomeFormatedString:cxo 321 SomeStringArray:[effr swxp ldnj obcs nvlg]}
```

Generate random placeholder images using [LoremFlickr](https://loremflickr.com/)

```go
// get a *os.File pointing to a file that is a random image
image := f.LoremFlickr().Image(100, 100, []string{}, "", false)

fmt.Println(image.Name())
// /tmp/loremflickr-img-4101493944.jpg
```

Generate profile images using [ThisPersonDoesNotExist](https://thispersondoesnotexist.com/)

```go
profileImage := f.ProfileImage().Image()

fmt.Println(profileImage.Name())
// /tmp/profil-picture-img-4022222298.jfif
```

See more formatters in [docs](https://pkg.go.dev/github.com/jaswdr/faker?tab=doc)

## Development

Create a fork and get the code.

```bash
$ go get github.com/jaswdr/faker
```

Do your changes, add tests, run the tests.

```bash
$ go test
PASS
ok      github.com/jaswdr/faker 0.010s
```

Push to your fork and send a new pull request from your fork to this repository.

## License

Faker is released under the MIT Licence. See the bundled LICENSE file for details.

## Maintainer

Created and maitained by Jonathan Schweder ([@jaswdr](https://github.com/jaswdr)) and [many others](https://github.com/jaswdr/faker/graphs/contributors)
