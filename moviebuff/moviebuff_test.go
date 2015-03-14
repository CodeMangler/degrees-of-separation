package moviebuff

import (
	"../graph"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchEntityDecodesJSONCorrectly(t *testing.T) {
	json := `{"url":"a-movie","type":"Movie","name":"A Movie",
	"movies":[{"name":"Movie One","url":"movie-one","role":"Role One"},{"name":"Movie Two","url":"movie-two","role":"Role Two"}],
	"cast":[{"url":"cast-one","name":"Cast One","role":"Role Three"},{"url":"cast-two","name":"Cast Two","role":"Role Four"}]}`
	server := serve(json)
	defer server.Close()
	baseURL = server.URL

	e, _ := fetchEntity("a-node")
	assert.Equal(t, "a-movie", e.URL)
	assert.Equal(t, "Movie", e.Type)
	assert.Equal(t, "A Movie", e.Name)
	expectedMovies := []mbConnection{mbConnection{URL: "movie-one", Name: "Movie One", Role: "Role One"},
		mbConnection{URL: "movie-two", Name: "Movie Two", Role: "Role Two"}}
	assert.Equal(t, expectedMovies, e.Movies)
	expectedCast := []mbConnection{mbConnection{URL: "cast-one", Name: "Cast One", Role: "Role Three"},
		mbConnection{URL: "cast-two", Name: "Cast Two", Role: "Role Four"}}
	assert.Equal(t, expectedCast, e.Cast)
}

func TestFetchEntityReturnsNilOnHTTPError(t *testing.T) {
	server := serve("", errors.New("A server error"), 500)
	defer server.Close()
	baseURL = server.URL

	entity, err := fetchEntity("a-non-existent-node")
	assert.Nil(t, entity)
	assert.Equal(t, "server error: 500: A server error\n", err.Error())
}

func TestFetchMapsPersonConnectionsFromMovies(t *testing.T) {
	json := `{"url":"person-node","type":"Person","name":"An Actor",
	"movies":[{"name":"Movie One","url":"movie-one","role":"Role One"},{"name":"Movie Two","url":"movie-two","role":"Role Two"}]}`
	server := serve(json)
	defer server.Close()
	baseURL = server.URL

	node := graph.NewNode("person-node")
	Fetch(node)
	assert.True(t, node.IsNeighbour(&graph.Node{ID: "movie-one"}))
	assert.True(t, node.IsNeighbour(&graph.Node{ID: "movie-two"}))
}

func TestFetchMapsMovieConnectionsFromCast(t *testing.T) {
	json := `{"url":"movie-node","type":"Movie","name":"A Movie",
    "cast":[{"url":"cast-one","name":"Cast One","role":"Role Three"},{"url":"cast-two","name":"Cast Two","role":"Role Four"}]}`
	server := serve(json)
	defer server.Close()
	baseURL = server.URL

	node := graph.NewNode("movie-node")
	Fetch(node)
	assert.True(t, node.IsNeighbour(&graph.Node{ID: "cast-one"}))
	assert.True(t, node.IsNeighbour(&graph.Node{ID: "cast-two"}))
}

func serve(json string, args ...interface{}) *httptest.Server {
	var err error
	errorCode := 500

	argsCount := len(args)
	if argsCount > 0 {
		err = args[0].(error)
	}
	if argsCount > 1 {
		errorCode = args[1].(int)
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			http.Error(w, err.Error(), errorCode)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, json)
	}))
}
