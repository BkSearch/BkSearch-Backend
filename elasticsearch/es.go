package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"log"

	"github.com/BkSearch/BkSearch-Backend/common"
	handleerror "github.com/BkSearch/BkSearch-Backend/handle-error"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	_ "go.opentelemetry.io/otel/trace"
)

type StackOverflow struct {
	client *elasticsearch.Client
	index  string
}

type indexedDocument struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Type     int    `json:"type"`
	Accepted bool   `json:"accepted"`
}

func NewStackOverflow(addresses []string) *StackOverflow {

	fmt.Println(addresses)
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "Create Client")
	}
	return &StackOverflow{
		client: client,
		index:  "stackoverflows",
	}
}

func (s *StackOverflow) Index(document common.Document) error {
	body := indexedDocument{
		ID:      document.ID,
		Title:   document.Title,
		Content: document.Content,
		Type:    int(document.Type),
    Accepted: document.Accepted,
	}

	data, err := json.Marshal(body)

	if err != nil {
		fmt.Println(err)
		handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}

	req := esapi.IndexRequest{
		Index:      s.index,
		DocumentID: document.ID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	resp, err := req.Do(context.Background(), s.client)
	if err != nil {
		return handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "IndexRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return handleerror.NewErrorf(handleerror.ErrorCodeUnknown, "IndexRequest.Do %s", resp.StatusCode)
	}

	io.Copy(ioutil.Discard, resp.Body)

	return nil
}

func (s *StackOverflow) Delete(document common.Document, id string) error {
	req := esapi.DeleteRequest{
		Index:      s.index,
		DocumentID: id,
	}

	resp, err := req.Do(context.Background(), s.client)
	if err != nil {
		return handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "DeleteRequest.Do")
	}

	defer resp.Body.Close()

	if resp.IsError() {
		return handleerror.NewErrorf(handleerror.ErrorCodeUnknown, "DeleteRequest.Do %s", resp.StatusCode)
	}

	io.Copy(ioutil.Discard, resp.Body)

	return nil
}

func (s *StackOverflow) Search(content *string, page int) ([]DocumentElasticSearchView, error) {
	if content == nil {
		return nil, nil
	}
  if page == 0 {
    page = 1
  }

	var buf bytes.Buffer
  should := make([]interface{}, 0, 6)
  should = append(should, map[string]interface{}{
    "match_phrase": map[string]interface{}{
      "title": map[string]interface{}{
        "query": *content,
        "boost": 4,
      },
    },
  })
  should = append(should, map[string]interface{}{
    "match_phrase": map[string]interface{}{
      "content": map[string]interface{}{
        "query": *content,
        "boost": 2,
      },
    },
  })

  should = append(should, map[string]interface{}{
    "match": map[string]interface{}{
      "title": map[string]interface{}{
        "query": *content,
        "operator": "and",
        "boost": 3,
      },
    },
  })
  should = append(should, map[string]interface{}{
    "match": map[string]interface{}{
      "content": map[string]interface{}{
        "query": *content,
        "operator": "and",
        "boost": 1,
      },
    },
  })
  should = append(should, map[string]interface{}{
    "match": map[string]interface{}{
      "title": map[string]interface{}{
        "query": *content,
        "operator": "or",
        "boost": 2,
      },
    },
  })

  should = append(should, map[string]interface{}{
    "match": map[string]interface{}{
      "content": map[string]interface{}{
        "query": *content,
        "operator": "or",
        "boost": 1,
      },
    },
  })

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": should,
			},
		},
		"from": (page-1)*10,
		"size": 10,
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "json.SearchRequest")
	}

	resp, err := s.client.Search(
		s.client.Search.WithContext(context.Background()),
		s.client.Search.WithIndex(s.index),
		s.client.Search.WithBody(&buf),
		s.client.Search.WithTrackTotalHits(true),
		s.client.Search.WithPretty(),
	)
	if err != nil {
		fmt.Println(err)
		handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "json.SearchRequest.Do")
	}

	defer resp.Body.Close()

	if resp.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				resp.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var hits struct {
		Hits struct {
			Hits []struct {
				Score  float64         `json:"_score"`
				Source indexedDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return nil, handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	res := make([]DocumentElasticSearchView, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		res[i].ID = hit.Source.ID
    res[i].Title = hit.Source.Title
		res[i].Content = hit.Source.Content
		res[i].Type = hit.Source.Type
		res[i].Score = hit.Score
	}

	return res, nil
}

func (s *StackOverflow) SearchMatchPhrase(keyword *string, page int) ([]DocumentElasticSearchView, error) {
	if keyword == nil {
		return nil, nil
	}

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_phrase": map[string]interface{}{
				"content": *keyword,
			},
		},
		"from": (page-1)*10 + 1,
		"size": 10,
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "json.SearchRequest")
	}

	resp, err := s.client.Search(
		s.client.Search.WithContext(context.Background()),
		s.client.Search.WithIndex(s.index),
		s.client.Search.WithBody(&buf),
		s.client.Search.WithTrackTotalHits(true),
		s.client.Search.WithPretty(),
	)

	if err != nil {
		fmt.Println(err)
		handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "json.SearchRequest.Do")
	}

	defer resp.Body.Close()

	if resp.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				resp.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var hits struct {
		Hits struct {
			Hits []struct {
				Score  float64         `json:"_score"`
				Source indexedDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return nil, handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	res := make([]DocumentElasticSearchView, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		res[i].ID = hit.Source.ID
		res[i].Content = hit.Source.Content
		res[i].Type = hit.Source.Type
		res[i].Score = hit.Score
	}

	return res, nil
}

func (s *StackOverflow) SearchWithListString(matchQuery []*string, matchPhraseQuery []*string, page int) ([]DocumentElasticSearchView, error) {
	if len(matchQuery) == 0 && len(matchPhraseQuery) == 0 {
		return nil, nil
	}

	match := make([]interface{}, 0, len(matchQuery))
	matchPhrase := make([]interface{}, 0, len(matchPhraseQuery))
	for _, query := range matchQuery {
		match = append(match, map[string]interface{}{
			"match": map[string]interface{}{
				"content": *query,
			},
		})
	}

	for _, query := range matchPhraseQuery {
		matchPhrase = append(matchPhrase, map[string]interface{}{
			"match_phrase": map[string]interface{}{
				"content": *query,
			},
		})
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": match,
				"must":   matchPhrase,
			},
		},
		"from": (page-1)*10 + 1,
		"size": 10,
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "json.SearchRequest")
	}

	resp, err := s.client.Search(
		s.client.Search.WithContext(context.Background()),
		s.client.Search.WithIndex(s.index),
		s.client.Search.WithBody(&buf),
		s.client.Search.WithTrackTotalHits(true),
		s.client.Search.WithPretty(),
	)

	if err != nil {
		fmt.Println(err)
		handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "json.SearchRequest.Do")
	}

	defer resp.Body.Close()

	if resp.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				resp.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var hits struct {
		Hits struct {
			Hits []struct {
				Score  float64         `json:"_score"`
				Source indexedDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return nil, handleerror.WrapErrorf(err, handleerror.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	res := make([]DocumentElasticSearchView, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		res[i].ID = hit.Source.ID
		res[i].Content = hit.Source.Content
		res[i].Type = hit.Source.Type
		res[i].Score = hit.Score
	}

	return res, nil
}
