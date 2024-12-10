package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ESQuery struct {
	Query struct {
		Match struct {
			Field struct {
				Query    string `json:"query"`
				Analyzer string `json:"analyzer"`
			} `json:"field"`
		} `json:"match"`
	} `json:"query"`
}

func AutoCompleteHandler(esEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if query == "" {
			http.Error(w, "missing query parameter", http.StatusBadRequest)
			return
		}

		index := "autocomplete_en_index"
		analyzer := "ngram_analyzer_en"
		if containsChinese(query) {
			index = "autocomplete_cn_index"
			analyzer = "ngram_analyzer_cn"
		}

		log.Println(query)
		// Construct Elasticsearch Query
		esQuery := ESQuery{}
		esQuery.Query.Match.Field.Query = query
		esQuery.Query.Match.Field.Analyzer = analyzer

		queryBody, err := json.Marshal(esQuery)
		if err != nil {
			http.Error(w, "failed to create query body", http.StatusInternalServerError)
			return
		}

		// Send request to Elasticsearch
		resp, err := http.Post(fmt.Sprintf("%s/%s/_search", esEndpoint, index), "application/json", bytes.NewReader(queryBody))
		if err != nil {
			log.Printf("Elasticsearch query failed: %v", err)
			http.Error(w, "failed to query Elasticsearch", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Read response from Elasticsearch
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "failed to read response from Elasticsearch", http.StatusInternalServerError)
			return
		}

		// Parse Elasticsearch response
		var esResponse struct {
			Hits struct {
				Hits []struct {
					Source struct {
						Field string `json:"field"`
					} `json:"_source"`
				} `json:"hits"`
			} `json:"hits"`
		}
		err = json.Unmarshal(body, &esResponse)
		if err != nil {
			http.Error(w, "failed to parse response from Elasticsearch", http.StatusInternalServerError)
			return
		}

		// Collect field values
		results := make([]string, 0, len(esResponse.Hits.Hits))
		for _, hit := range esResponse.Hits.Hits {
			results = append(results, hit.Source.Field)
		}

		// Return results as a JSON array
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

func containsChinese(s string) bool {
	for _, r := range s {
		if r >= '\u4e00' && r <= '\u9fff' {
			return true
		}
	}
	return false
}
