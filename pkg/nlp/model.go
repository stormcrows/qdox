package nlp

import (
	"fmt"
	"sort"

	"github.com/james-bowman/nlp"
	"github.com/james-bowman/nlp/measures/pairwise"
	"gonum.org/v1/gonum/mat"
)

// Model holds processing pipeline and trained matrix, along with reference to the corpus
type Model struct {
	Pipeline *nlp.Pipeline
	Matrix   mat.Matrix
	Corpus   *Corpus
}

// QueryResult contains indexes of matched documents along with their similarities
type QueryResult struct {
	Query        string
	Matched      []int
	Similarities []float64
	Err          error
}

var stopWords = []string{"a", "about", "above", "above", "across", "after", "afterwards", "again", "against", "all", "almost", "alone", "along", "already", "also", "although", "always", "am", "among", "amongst", "amoungst", "amount", "an", "and", "another", "any", "anyhow", "anyone", "anything", "anyway", "anywhere", "are", "around", "as", "at", "back", "be", "became", "because", "become", "becomes", "becoming", "been", "before", "beforehand", "behind", "being", "below", "beside", "besides", "between", "beyond", "bill", "both", "bottom", "but", "by", "call", "can", "cannot", "cant", "co", "con", "could", "couldnt", "cry", "de", "describe", "detail", "do", "done", "down", "due", "during", "each", "eg", "eight", "either", "eleven", "else", "elsewhere", "empty", "enough", "etc", "even", "ever", "every", "everyone", "everything", "everywhere", "except", "few", "fifteen", "fify", "fill", "find", "fire", "first", "five", "for", "former", "formerly", "forty", "found", "four", "from", "front", "full", "further", "get", "give", "go", "had", "has", "hasnt", "have", "he", "hence", "her", "here", "hereafter", "hereby", "herein", "hereupon", "hers", "herself", "him", "himself", "his", "how", "however", "hundred", "ie", "if", "in", "inc", "indeed", "interest", "into", "is", "it", "its", "itself", "keep", "last", "latter", "latterly", "least", "less", "ltd", "made", "many", "may", "me", "meanwhile", "might", "mill", "mine", "more", "moreover", "most", "mostly", "move", "much", "must", "my", "myself", "name", "namely", "neither", "never", "nevertheless", "next", "nine", "no", "nobody", "none", "noone", "nor", "not", "nothing", "now", "nowhere", "of", "off", "often", "on", "once", "one", "only", "onto", "or", "other", "others", "otherwise", "our", "ours", "ourselves", "out", "over", "own", "part", "per", "perhaps", "please", "put", "rather", "re", "same", "see", "seem", "seemed", "seeming", "seems", "serious", "several", "she", "should", "show", "side", "since", "sincere", "six", "sixty", "so", "some", "somehow", "someone", "something", "sometime", "sometimes", "somewhere", "still", "such", "system", "take", "ten", "than", "that", "the", "their", "them", "themselves", "then", "thence", "there", "thereafter", "thereby", "therefore", "therein", "thereupon", "these", "they", "thickv", "thin", "third", "this", "those", "though", "three", "through", "throughout", "thru", "thus", "to", "together", "too", "top", "toward", "towards", "twelve", "twenty", "two", "un", "under", "until", "up", "upon", "us", "very", "via", "was", "we", "well", "were", "what", "whatever", "when", "whence", "whenever", "where", "whereafter", "whereas", "whereby", "wherein", "whereupon", "wherever", "whether", "which", "while", "whither", "who", "whoever", "whole", "whom", "whose", "why", "will", "with", "within", "without", "would", "yet", "you", "your", "yours", "yourself", "yourselves"}

// NewLSIModel initializes LSI pipeline
func NewLSIModel() *Model {
	vectoriser := nlp.NewCountVectoriser(stopWords...)
	transformer := nlp.NewTfidfTransformer()
	reducer := nlp.NewTruncatedSVD(4)
	pipeline := nlp.NewPipeline(vectoriser, transformer, reducer)

	return &Model{pipeline, nil, nil}
}

// Train fits the model to the given corpus, resulting in lsi matrix
func (m *Model) Train(c *Corpus) error {
	lsi, err := m.Pipeline.FitTransform(c.Contents()...)
	if err != nil {
		return fmt.Errorf("Failed to process documents: %q", err.Error())
	}
	m.Matrix = lsi
	c.Release()
	m.Corpus = c
	return nil
}

// Query returns document indexes matching given query
func (m *Model) Query(q string, n int, threshold float64) QueryResult {
	queryVector, err := m.Pipeline.Transform(q)
	if err != nil {
		return QueryResult{q, nil, nil, fmt.Errorf("Failed to process documents: %q", err.Error())}
	}
	_, docs := m.Matrix.Dims()
	matched := make([]int, 0)
	similarities := make([]float64, 0)
	maxSimilarity := -1.0
	for i := 0; i < docs; i++ {
		similarity := pairwise.CosineSimilarity(queryVector.(mat.ColViewer).ColView(0), m.Matrix.(mat.ColViewer).ColView(i))
		if similarity >= threshold && similarity >= maxSimilarity {
			maxSimilarity = similarity
			matched = append(matched, i)
			similarities = append(similarities, similarity)
		}
	}

	qr := QueryResult{q, matched, similarities, nil}
	sort.Sort(&qr)

	if len(qr.Matched) > n {
		qr.Matched = qr.Matched[0:n]
		qr.Similarities = qr.Similarities[0:n]
	}

	return qr
}

// Ordering of results

func (qr *QueryResult) Len() int {
	return len(qr.Matched)
}

func (qr *QueryResult) Swap(i, j int) {
	qr.Matched[i], qr.Matched[j] = qr.Matched[j], qr.Matched[i]
	qr.Similarities[i], qr.Similarities[j] = qr.Similarities[j], qr.Similarities[i]
}

func (qr *QueryResult) Less(i, j int) bool {
	return qr.Similarities[j] < qr.Similarities[i]
}
