package faker

import (
	"strings"
)

var (
	wordsList = []string{"a", "in", "et", "ut", "ut", "ad", "et", "at", "id", "et", "ut", "in", "ab", "ea", "ut", "et", "et", "et", "et", "et", "et", "ea", "id", "et", "et", "ut", "ut", "ex", "est", "sed", "qui", "est", "est", "aut", "eos", "qui", "cum", "nam", "non", "aut", "qui", "sed", "qui", "vel", "non", "sit", "rem", "eos", "qui", "qui", "sed", "est", "non", "est", "sit", "eum", "hic", "quo", "sit", "aut", "aut", "vel", "aut", "eum", "aut", "quo", "odio", "enim", "unde", "illo", "sunt", "quis", "sint", "sint", "quas", "fuga", "modi", "enim", "quos", "odit", "quia", "sunt", "eius", "quia", "quia", "nisi", "iste", "quam", "vero", "amet", "ipsa", "esse", "quis", "quae", "quia", "nemo", "iure", "quod", "illum", "ipsum", "dolor", "rerum", "velit", "culpa", "omnis", "nihil", "minus", "saepe", "iusto", "velit", "magni", "alias", "omnis", "porro", "autem", "nihil", "totam", "fugit", "dolor", "optio", "atque", "autem", "ipsam", "nobis", "nulla", "ullam", "rerum", "harum", "eaque", "error", "animi", "dicta", "vitae", "quasi", "natus", "earum", "rerum", "omnis", "neque", "sequi", "libero", "soluta", "cumque", "beatae", "maxime", "facere", "quidem", "labore", "dolore", "veniam", "minima", "fugiat", "itaque", "magnam", "dolorem", "laborum", "nostrum", "quaerat", "officia", "maiores", "facilis", "dolorem", "aliquam", "numquam", "aliquid", "dolorum", "aperiam", "tempore", "dolores", "eveniet", "dolores", "debitis", "commodi", "tempora", "ratione", "ducimus", "tenetur", "placeat", "impedit", "quisquam", "nesciunt", "adipisci", "pariatur", "deleniti", "voluptas", "incidunt", "repellat", "eligendi", "possimus", "corporis", "expedita", "sapiente", "delectus", "suscipit", "voluptas", "deserunt", "mollitia", "corrupti", "voluptas", "officiis", "accusamus", "similique", "doloribus", "provident", "occaecati", "quibusdam", "assumenda", "inventore", "veritatis", "explicabo", "voluptate", "molestiae", "molestias", "excepturi", "molestiae", "recusandae", "asperiores", "voluptatem", "reiciendis", "laudantium", "voluptatem", "temporibus", "voluptatum", "voluptatem", "laboriosam", "aspernatur", "voluptates", "voluptatem", "distinctio", "architecto", "cupiditate", "doloremque", "blanditiis", "dignissimos", "repellendus", "consequatur", "accusantium", "consectetur", "repudiandae", "consequatur", "praesentium", "perferendis", "consequatur", "voluptatibus", "perspiciatis", "consequuntur", "reprehenderit", "necessitatibus", "exercitationem"}
)

// Lorem is a faker struct for Lorem
type Lorem struct {
	Faker *Faker
}

// Word returns a fake word for Lorem
func (l Lorem) Word() string {
	index := l.Faker.IntBetween(0, len(wordsList)-1)
	return wordsList[index]
}

// Words returns fake words for Lorem
func (l Lorem) Words(nbWords int) (words []string) {
	for i := 0; i < nbWords; i++ {
		words = append(words, l.Word())
	}

	return
}

// Sentence returns a fake sentence for Lorem
func (l Lorem) Sentence(nbWords int) string {
	return strings.Join(l.Words(nbWords), " ") + "."
}

// Sentences returns fake sentences for Lorem
func (l Lorem) Sentences(nbSentences int) (sentences []string) {
	for i := 0; i < nbSentences; i++ {
		sentences = append(sentences, l.Sentence(l.Faker.RandomNumber(2)))
	}

	return
}

// Paragraph returns a fake paragraph for Lorem
func (l Lorem) Paragraph(nbSentences int) string {
	return strings.Join(l.Sentences(nbSentences), " ")
}

// Paragraphs returns fake paragraphs for Lorem
func (l Lorem) Paragraphs(nbParagraph int) (out []string) {
	for i := 0; i < nbParagraph; i++ {
		out = append(out, l.Paragraph(l.Faker.RandomNumber(2)))
	}

	return
}

// Text returns a fake text for Lorem
func (l Lorem) Text(maxNbChars int) (out string) {
	for _, w := range wordsList {
		if len(out)+len(w) > maxNbChars {
			break
		}

		out = out + w
	}

	return
}

// Bytes returns fake bytes for Lorem
func (l Lorem) Bytes(maxNbChars int) (out []byte) {
	return []byte(l.Text(maxNbChars))
}
