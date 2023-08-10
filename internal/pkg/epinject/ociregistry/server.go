// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ociregistry

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/oklog/ulid/v2"
	"github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

type session struct {
	remote string
	buf    *bytes.Buffer
}

// Server is an http.Handler that provides an OCI registry proxy. The proxy can inject
// the waypoint entrypoint directly into the image as it passes through the proxy.
type Server struct {
	// Upstream is the partial url for the upstream, for instance, https://hub.docker.io
	Upstream string

	// Logger is the hclog Logger to log with.
	Logger hclog.Logger

	// AuthConfig is used to authenticate with the target repository
	AuthConfig authn.AuthConfig

	// Indicates that the entrypoint should not be pulled into the image.
	DisableEntrypoint bool

	// client is set by Negotiate to contain a client that will auth correctly.
	client *http.Client

	routerOnce sync.Once
	router     *mux.Router

	entrypointId     digest.Digest
	entrypointDiffId digest.Digest
	entrypointSize   int64
	entrypointRepo   string
	entrypointData   []byte

	mu             sync.Mutex
	uploadSessions map[string]*session
	jsonBlobs      map[string][]byte
}

var _ http.Handler = (*Server)(nil)

// BasicAuth generates the Authorization header for basic auth with the given
// username and password.
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

/*
Endpoints
ID	    Method	API Endpoint	          Success	Failure
end-1	  GET	    /v2/	                    200	404/401
end-2	  G/H	    /v2/<name>/blobs/<digest>	200	404
end-3	  G/H   	/v2/<name>/manifests/<reference>	200	404
end-4a	POST	  /v2/<name>/blobs/uploads/	202	404
end-4b	POST	  /v2/<name>/blobs/uploads/?digest=<digest>	201/202	404/400
end-5	  PATCH	  /v2/<name>/blobs/uploads/<reference>	202	404/416
end-6	  PUT	    /v2/<name>/blobs/uploads/<reference>?digest=<digest>	201	404/400
end-7	  PUT	    /v2/<name>/manifests/<reference>	201	404
end-8a	GET	    /v2/<name>/tags/list	200	404
end-8b	GET	    /v2/<name>/tags/list?n=<integer>&last=<integer>	200	404
end-9	  DELETE	/v2/<name>/manifests/<reference>	202	404/400/405
end-10	DELETE	/v2/<name>/blobs/<digest>	202	404/405
end-11	POST	  /v2/<name>/blobs/uploads/?mount=<digest>&from=<other_name>	201	404
*/

const (
	routePrefix          = "/v2"
	distAPIVersion       = "Docker-Distribution-API-Version"
	distContentDigestKey = "Docker-Content-Digest"
	blobUploadUUID       = "Blob-Upload-UUID"
	defaultMediaType     = "application/json"
	binaryMediaType      = "application/octet-stream"

	// The string put into the config history when the entrypoint is injected.
	// If a history entry already exists with this string, we presume the entrypoint
	// is already there.
	createdBy = "entrypoint injection"

	dockerRootFSMediaType = "application/vnd.docker.image.rootfs.diff.tar.gzip"
)

var (
	nameRegexp = regexp.MustCompile(`[a-z0-9]+(?:[._-][a-z0-9]+)*(?:/[a-z0-9]+(?:[._-][a-z0-9]+)*)*`)
	refRegexp  = regexp.MustCompile(`[a-zA-Z0-9_][a-zA-Z0-9._-]{0,127}`)
)

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if s.client == nil {
		http.Error(w, "server misconfigured, no client", http.StatusInternalServerError)
		return
	}

	s.routerOnce.Do(func() {
		r := mux.NewRouter()
		g := r.PathPrefix(routePrefix).Subrouter()
		g.HandleFunc(
			fmt.Sprintf("/{name:%s}/tags/list", nameRegexp.String()), s.listTags).Methods("GET")
		g.HandleFunc(
			fmt.Sprintf("/{name:%s}/manifests/{reference}", nameRegexp.String()), s.checkManifest).Methods("HEAD")
		g.HandleFunc(
			fmt.Sprintf("/{name:%s}/manifests/{reference}", nameRegexp.String()), s.getManifest).Methods("GET")
		g.HandleFunc(
			fmt.Sprintf("/{name:%s}/manifests/{reference}", nameRegexp.String()), s.updateManifest).Methods("PUT")
		// g.HandleFunc(fmt.Sprintf("/{name:%s}/manifests/{reference}", NameRegexp.String()),
		// s.deleteManifest).Methods("DELETE")
		g.HandleFunc(
			fmt.Sprintf("/{name:%s}/blobs/{digest}", nameRegexp.String()), s.checkBlob).Methods("HEAD")
		g.HandleFunc(
			fmt.Sprintf("/{name:%s}/blobs/{digest}", nameRegexp.String()), s.getBlob).Methods("GET")
		// g.HandleFunc(fmt.Sprintf("/{name:%s}/blobs/{digest}", NameRegexp.String()),
		// s.deleteBlob).Methods("DELETE")
		g.HandleFunc(fmt.Sprintf("/{name:%s}/blobs/uploads/", nameRegexp.String()),
			s.createBlobUpload).Methods("POST")
		// g.HandleFunc(fmt.Sprintf("/{name:%s}/blobs/uploads/{session_id}", NameRegexp.String()),
		// s.getBlobUpload).Methods("GET")
		g.HandleFunc(fmt.Sprintf("/{name:%s}/blobs/uploads/{session_id}", nameRegexp.String()),
			s.patchBlobUpload).Methods("PATCH")
		g.HandleFunc(fmt.Sprintf("/{name:%s}/blobs/uploads/{session_id}", nameRegexp.String()),
			s.updateBlobUpload).Methods("PUT")
		// g.HandleFunc(fmt.Sprintf("/{name:%s}/blobs/uploads/{session_id}", NameRegexp.String()),
		// s.deleteBlobUpload).Methods("DELETE")
		g.HandleFunc(
			"/_catalog", s.listRepositories).Methods("GET")
		g.HandleFunc(
			"/", s.checkVersionSupport).Methods("GET")

		s.router = r
	})

	s.Logger.Info("request",
		"method", req.Method,
		"url", req.URL.String(),
	)

	s.router.ServeHTTP(w, req)
}

func (s *Server) getClient(ctx context.Context, repoStr string) (*http.Client, error) {
	// So the big deal with using transport from go-containerregistry is that it handles basic
	// auth as well as the bearer token protocol used when the server returns WWW-Authenticate.
	// Plus it's well tested.

	repo, err := name.NewRepository(repoStr)
	if err != nil {
		return nil, err
	}

	auth := authn.FromConfig(s.AuthConfig)

	scopes := []string{repo.Scope(transport.PushScope), repo.Scope(transport.PullScope)}
	t, err := transport.NewWithContext(ctx, repo.Registry, auth, http.DefaultTransport, scopes)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Transport: t}

	return client, nil
}

func (s *Server) Negotiate(repo string) error {
	client, err := s.getClient(context.Background(), repo)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/", s.Upstream), nil)
	if err != nil {
		return err
	}

	s.Logger.Debug("sending ping request to /v2/ to test authentication")

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "attempting to contact registry: %s", s.Upstream)
	}

	defer resp.Body.Close()

	s.Logger.Debug("response from ping", "status", resp.Status)

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		s.client = client
		return nil
	}

	return errors.Errorf("Unexpected status: %s", resp.Status)
}

// SetupEntrypointLayer should be called before the server is in use. This
// will write the entrypoint as a layer into the upstream, so that it doesn't
// have to be written while attempting to update the manifest. The entrypoint
// could be large and writing it mid-manifest update could cause a timeout
// on the client side.
func (s *Server) SetupEntrypointLayer(name string, epData []byte) error {
	var buf bytes.Buffer

	diffid := digest.Canonical.Digester()
	layerid := digest.Canonical.Digester()

	dh := diffid.Hash()
	lh := layerid.Hash()

	gzw := gzip.NewWriter(io.MultiWriter(&buf, lh))
	tw := tar.NewWriter(io.MultiWriter(gzw, dh))

	tw.WriteHeader(&tar.Header{
		Name: "/waypoint-entrypoint",
		Size: int64(len(epData)),
		Mode: 0777,
	})

	tw.Write(epData)

	tw.Close()
	gzw.Close()

	s.entrypointRepo = name
	s.entrypointData = buf.Bytes()
	s.entrypointSize = int64(buf.Len())
	s.entrypointId = layerid.Digest()
	s.entrypointDiffId = diffid.Digest()

	s.Logger.Debug("entrypoint layer", "layer-id", s.entrypointId.String(), "size", s.entrypointSize)

	return s.uploadEntrypoint()
}

// FetchBlob will read a blob in the OCI digest format from the upstream repo.
func (s *Server) FetchBlob(name string, id string) ([]byte, error) {
	return s.retrieveBlob(digest.Digest(id), name)
}

func (s *Server) uploadEntrypoint() error {
	found, err := s.probeBlob(s.entrypointRepo, s.entrypointId)
	if err != nil {
		return err
	}

	if found {
		s.Logger.Debug("entrypoint layer already exists, not reuploading")
		return nil
	}
	s.Logger.Debug("uploading entrypoint layer")
	if _, err := s.writeBlob(s.entrypointRepo, s.entrypointData); err != nil {
		return err
	}

	s.Logger.Debug("uploaded entrypoint layer")
	return nil
}

// CheckVersionSupport godoc
// @Summary Check API support
// @Description Check if this API version is supported
// @Router 	/v2/ [get]
// @Accept  json
// @Produce json
// @Success 200 {string} string	"ok".
func (s *Server) checkVersionSupport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(distAPIVersion, "registry/2.0")
	w.Header().Set("Content-Type", defaultMediaType)
	w.Write(nil)
}

// ListTags godoc
// @Summary List image tags
// @Description List all image tags in a repository
// @Router 	/v2/{name}/tags/list [get]
// @Accept  json
// @Produce json
// @Param   name     path    string     true        "test"
// @Param 	n	 			 query 	 integer 		true				"limit entries for pagination"
// @Param 	last	 	 query 	 string 		true				"last tag value for pagination"
// @Success 200 {object} 	api.ImageTags
// @Failure 404 {string} 	string 				"not found"
// @Failure 400 {string} 	string 				"bad request".
func (s *Server) listTags(w http.ResponseWriter, r *http.Request) {
	s.proxy(w, r)
}

// CheckManifest godoc
// @Summary Check image manifest
// @Description Check an image's manifest given a reference or a digest
// @Router 	/v2/{name}/manifests/{reference} [head]
// @Accept  json
// @Produce json
// @Param   name     			path    string     true        "repository name"
// @Param   reference     path    string     true        "image reference or digest"
// @Success 200 {string} string	"ok"
// @Header  200 {object} api.DistContentDigestKey
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error".
func (s *Server) checkManifest(w http.ResponseWriter, r *http.Request) {
	s.proxy(w, r)
}

func (s *Server) outRequest(method, url string, r io.Reader) (*http.Request, error) {
	out, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Server) proxy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	path := r.URL.Path
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}

	var input io.Reader

	if strings.HasSuffix(r.Header.Get("Content-Type"), "+json") {
		data, _ := ioutil.ReadAll(r.Body)
		input = bytes.NewReader(data)
	} else {
		input = r.Body
	}

	out, err := s.outRequest(r.Method, s.Upstream+path, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.Logger.Info("proxy", "url", out.URL.String())

	for k, v := range r.Header {
		out.Header[k] = v
	}

	resp, err := s.client.Do(out)
	if err != nil {
		s.Logger.Error("do error", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// resp.Write(os.Stdout)

	s.Logger.Info("proxy response", "status", resp.Status)

	defer resp.Body.Close()

	hdr := w.Header()

	for k, v := range resp.Header {
		if k == "Location" {
			u2, err := url.Parse(v[0])
			if err == nil {
				u2.Scheme = ""
				u2.Host = ""

				v[0] = u2.String()
			}
		}

		hdr[k] = v
	}

	w.WriteHeader(resp.StatusCode)

	if strings.Contains(resp.Header.Get("Content-Type"), "json") {
		io.Copy(w, io.TeeReader(resp.Body, os.Stdout))
	} else {
		io.Copy(w, resp.Body)
	}
}

// GetManifest godoc
// @Summary Get image manifest
// @Description Get an image's manifest given a reference or a digest
// @Accept  json
// @Produce application/vnd.oci.image.manifest.v1+json
// @Param   name     			path    string     true        "repository name"
// @Param   reference     path    string     true        "image reference or digest"
// @Success 200 {object} 	api.ImageManifest
// @Header  200 {object} api.DistContentDigestKey
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /v2/{name}/manifests/{reference} [get].
func (s *Server) getManifest(w http.ResponseWriter, r *http.Request) {
	s.proxy(w, r)
}

func (s *Server) retrieveBlob(dig digest.Digest, name string) ([]byte, error) {
	req, err := s.outRequest("GET", fmt.Sprintf("%s/v2/%s/blobs/%s", s.Upstream, name, dig), nil)
	if err != nil {
		return nil, err
	}

	s.Logger.Info("retrieving blob", "method", req.Method, "url", req.URL.String())

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Server) probeBlob(name string, dig digest.Digest) (bool, error) {
	req, err := s.outRequest("HEAD", fmt.Sprintf("%s/v2/%s/blobs/%s", s.Upstream, name, dig), nil)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode <= 299, nil
}

func (s *Server) writeBlob(name string, data []byte) (digest.Digest, error) {
	dig := digest.Canonical.FromBytes(data)

	req, err := s.outRequest("POST", fmt.Sprintf("%s/v2/%s/blobs/uploads/", s.Upstream, name), nil)
	if err != nil {
		return "", err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("error starting blob upload: %s", resp.Status)
	}

	loc := resp.Header.Get("Location")

	if strings.Contains(loc, "?") {
		loc += "&digest=" + dig.String()
	} else {
		loc += "?digest=" + dig.String()
	}

	r2, err := s.outRequest("PUT", loc, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	resp, err = s.client.Do(r2)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return "", fmt.Errorf("error attempting to put entrypoint layer: %s", resp.Status)
	}

	return dig, nil
}

// UpdateManifest godoc
// @Summary Update image manifest
// @Description Update an image's manifest given a reference or a digest
// @Accept  json
// @Produce json
// @Param   name     			path    string     true        "repository name"
// @Param   reference     path    string     true        "image reference or digest"
// @Header  201 {object} api.DistContentDigestKey
// @Success 201 {string} string	"created"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /v2/{name}/manifests/{reference} [put].
func (s *Server) updateManifest(w http.ResponseWriter, r *http.Request) {
	if s.DisableEntrypoint {
		s.proxy(w, r)
		return
	}

	defer r.Body.Close()

	t := mux.Vars(r)
	name, ok := t["name"]
	if !ok {
		http.Error(w, "missing name", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading manifest body", http.StatusBadRequest)
		return
	}

	// We do this because v1.Manifest doesn't have the MediaType but the OG docker registry
	// requires it. So by doing this, we preserve it round trip as we manipulate it.
	var man struct {
		v1.Manifest
		MediaType string `json:"mediaType"`
	}

	err = json.Unmarshal(data, &man)
	if err != nil {
		s.Logger.Error("error decoding manifest", "error", err)
		http.Error(w, "bad manifest", http.StatusBadRequest)
		return
	}

	conKey := man.Config.Digest.String()

	configData, ok := s.jsonBlobs[conKey]
	if !ok {
		s.Logger.Error("missing config json in json blobs", "digest", conKey)
		http.Error(w, "bad config", http.StatusInternalServerError)
		return
	}

	var config v1.Image
	err = json.Unmarshal(configData, &config)
	if err != nil {
		s.Logger.Error("error unmarshaling config", "error", err)
		http.Error(w, "bad config", http.StatusBadRequest)
		return
	}

	// detect if the entrypoint injection already took place

	// injectEntrypoint is used to determine if the image has or should have
	// Entrypoint injected.
	injectEntrypoint := true
	for _, h := range config.History {
		if h.CreatedBy == createdBy {
			injectEntrypoint = false
		}
	}

	// Ok, now delete the config from the jsonBlobs and upload any remaining json
	// blobs so the upstream knows about them.
	delete(s.jsonBlobs, conKey)
	for id, data := range s.jsonBlobs {
		s.writeBlob(id, data)
	}

	now := time.Now()

	if injectEntrypoint {
		config.RootFS.DiffIDs = append(config.RootFS.DiffIDs, s.entrypointDiffId)
		config.History = append(config.History, v1.History{
			Created:   &now,
			CreatedBy: createdBy,
			Author:    "waypoint",
		})
	}

	// By default we want to prepend waypoint-entrypoint to the Entrypoint and
	// inject the CEB. If the Entrypoint is empty, or the first element is
	// already waypoint-entrypoint, we skip the injection
	if len(config.Config.Entrypoint) == 0 || config.Config.Entrypoint[0] != "/waypoint-entrypoint" {
		config.Config.Entrypoint = append([]string{"/waypoint-entrypoint"}, config.Config.Entrypoint...)
		s.Logger.Debug("injected entrypoint", "value", config.Config.Entrypoint)
	} else {
		s.Logger.Debug("entrypoint already included", "value", config.Config.Entrypoint)
	}

	newConfig, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		http.Error(w, "error updating config", http.StatusInternalServerError)
		return
	}

	id, err := s.writeBlob(name, newConfig)
	if err != nil {
		http.Error(w, "error writing new config", http.StatusInternalServerError)
		return
	}

	man.Config.Size = int64(len(newConfig))
	man.Config.Digest = id

	mediaType := v1.MediaTypeImageLayerGzip

	if strings.Contains(man.Layers[0].MediaType, "docker") {
		mediaType = dockerRootFSMediaType
	}

	if injectEntrypoint {
		// Upload Entrypoint again just to be sure the layer is still there
		// before we go and reference it. This protects against the repo
		// deleting the layer blob between server start and here.
		err := s.uploadEntrypoint()
		if err != nil {
			http.Error(w, "error uploading entrypoint", http.StatusInternalServerError)
			return
		}

		man.Layers = append(man.Layers, v1.Descriptor{
			MediaType: mediaType,
			Digest:    s.entrypointId,
			Size:      s.entrypointSize,
		})
	}

	newMan, err := json.MarshalIndent(man, "", "  ")
	if err != nil {
		http.Error(w, "error marshaling manifest", http.StatusInternalServerError)
		return
	}

	s.Logger.Info("updated manifest")

	r.Body = ioutil.NopCloser(bytes.NewReader(newMan))

	s.proxy(w, r)
}

// CheckBlob godoc
// @Summary Check image blob/layer
// @Description Check an image's blob/layer given a digest
// @Accept  json
// @Produce json
// @Param   name				path    string     true        "repository name"
// @Param   digest     	path    string     true        "blob/layer digest"
// @Success 200 {object} api.ImageManifest
// @Header  200 {object} api.DistContentDigestKey
// @Router /v2/{name}/blobs/{digest} [head].
func (s *Server) checkBlob(w http.ResponseWriter, r *http.Request) {
	if s.DisableEntrypoint {
		s.proxy(w, r)
		return
	}

	vars := mux.Vars(r)
	digest, ok := vars["digest"]

	if !ok || digest == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if data, ok := s.jsonBlobs[digest]; ok {
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.WriteHeader(200)
		return
	}

	s.proxy(w, r)
}

// GetBlob godoc
// @Summary Get image blob/layer
// @Description Get an image's blob/layer given a digest
// @Accept  json
// @Produce application/vnd.oci.image.layer.v1.tar+gzip
// @Param   name				path    string     true        "repository name"
// @Param   digest     	path    string     true        "blob/layer digest"
// @Header  200 {object} api.DistContentDigestKey
// @Success 200 {object} api.ImageManifest
// @Router /v2/{name}/blobs/{digest} [get].
func (s *Server) getBlob(w http.ResponseWriter, r *http.Request) {
	s.proxy(w, r)
}

// ListRepositories godoc
// @Summary List image repositories
// @Description List all image repositories
// @Accept  json
// @Produce json
// @Success 200 {object} 	api.RepositoryList
// @Failure 500 {string} string "internal server error"
// @Router /v2/_catalog [get].
func (s *Server) listRepositories(w http.ResponseWriter, r *http.Request) {
	out, err := s.outRequest("GET", s.Upstream+"/v2/_catalog", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := s.client.Do(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	hdr := w.Header()

	for k, v := range resp.Header {
		hdr[k] = v
	}

	io.Copy(w, resp.Body)
}

// CreateBlobUpload godoc
// @Summary Create image blob/layer upload
// @Description Create a new image blob/layer upload
// @Accept  json
// @Produce json
// @Param   name				path    string     true        "repository name"
// @Success 202 {string} string	"accepted"
// @Header  202 {string} Location "/v2/{name}/blobs/uploads/{session_id}"
// @Header  202 {string} Range "bytes=0-0"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /v2/{name}/blobs/uploads [post].
func (s *Server) createBlobUpload(w http.ResponseWriter, r *http.Request) {
	if s.DisableEntrypoint {
		s.proxy(w, r)
		return
	}

	vars := mux.Vars(r)
	name, ok := vars["name"]

	if !ok || name == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	id, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		http.Error(w, "error creating ulid", http.StatusInternalServerError)
		return
	}

	sid := id.String()

	s.mu.Lock()
	if s.uploadSessions == nil {
		s.uploadSessions = make(map[string]*session)
	}

	s.uploadSessions[sid] = &session{}
	s.mu.Unlock()

	s.Logger.Info("created new upload session", "id", sid)

	w.Header().Set("Location", fmt.Sprintf("/v2/%s/blobs/uploads/%s", name, sid))
	w.WriteHeader(202)
}

// PatchBlobUpload godoc
// @Summary Resume image blob/layer upload
// @Description Resume an image's blob/layer upload given an session_id
// @Accept  json
// @Produce json
// @Param   name     path    string     true        "repository name"
// @Param   session_id     path    string     true        "upload session_id"
// @Success 202 {string} string	"accepted"
// @Header  202 {string} Location "/v2/{name}/blobs/uploads/{session_id}"
// @Header  202 {string} Range "bytes=0-128"
// @Header  200 {object} api.BlobUploadUUID
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Failure 416 {string} string "range not satisfiable"
// @Failure 500 {string} string "internal server error"
// @Router /v2/{name}/blobs/uploads/{session_id} [patch].
func (s *Server) patchBlobUpload(w http.ResponseWriter, r *http.Request) {
	if s.DisableEntrypoint {
		s.proxy(w, r)
		return
	}

	defer r.Body.Close()

	vars := mux.Vars(r)
	name, ok := vars["name"]

	if !ok || name == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sid, ok := vars["session_id"]
	if !ok || sid == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	s.mu.Lock()
	sess := s.uploadSessions[sid]
	s.mu.Unlock()

	if sess == nil {
		s.Logger.Error("missing session", "id", sid)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var peekData []byte

	if r.Header.Get("Content-Length") == "" || r.Header.Get("Content-Range") == "" {
		peekData = make([]byte, 1)

		_, err := io.ReadFull(r.Body, peekData)
		if err != nil {
			http.Error(w, "unable to read body", http.StatusBadRequest)
			return
		}

		if peekData[0] == '{' {
			s.Logger.Debug("detected json blob")

			var buf bytes.Buffer
			buf.Write(peekData)

			io.Copy(&buf, r.Body)
			sess.buf = &buf

			wh := w.Header()

			wh.Set("Location", fmt.Sprintf("/v2/%s/blobs/uploads/%s", name, sid))
			wh.Set("Docker-Upload-Uuid", sid)
			wh.Set(distAPIVersion, "registry/2.0")
			wh.Set("Content-Length", "0")
			wh.Set("Range", fmt.Sprintf("0-%d", buf.Len()))
			wh.Set("X-Content-Type-Options", "nosniff")
			wh.Set("Date", time.Now().Format(time.RFC1123))

			w.WriteHeader(202)

			s.Logger.Debug("returned patch response for json blob")
			return
		}
	}

	// It's not json, so go ahead and just send it to the target registry. We start a new upload
	// session to get the remote location if we haven't already.
	if sess.remote == "" {
		s.Logger.Debug("getting upstream blob location url")

		req, err := s.outRequest("POST", fmt.Sprintf("%s/v2/%s/blobs/uploads/", s.Upstream, name), nil)
		if err != nil {
			s.Logger.Error("error creating upstream location", "error", err)
			http.Error(w, "error creating upstream request", http.StatusInternalServerError)
			return
		}

		resp, err := s.client.Do(req)
		if err != nil {
			s.Logger.Error("error creating upstream location", "error", err)
			http.Error(w, "error creating upstream request", http.StatusInternalServerError)
			return
		}

		resp.Body.Close()

		// Some error
		if resp.StatusCode >= 300 {
			s.Logger.Error("upstream failed to create upload session", "status", resp.Status, "url", req.URL.String())

			wh := w.Header()

			for k, v := range resp.Header {
				wh[k] = v
			}

			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
			return
		}

		sess.remote = resp.Header.Get("Location")
	}

	s.Logger.Debug("uploading blob", "target", sess.remote)

	out, err := s.outRequest(r.Method, sess.remote, io.MultiReader(bytes.NewReader(peekData), r.Body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.Logger.Info("proxy", "url", out.URL.String())

	for k, v := range r.Header {
		out.Header[k] = v
	}

	resp, err := s.client.Do(out)
	if err != nil {
		s.Logger.Error("do error", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sess.remote = resp.Header.Get("Location")

	// resp.Write(os.Stdout)

	s.Logger.Info("proxy response", "status", resp.Status)

	defer resp.Body.Close()

	hdr := w.Header()

	for k, v := range resp.Header {
		hdr[k] = v
	}

	hdr.Set("Location", fmt.Sprintf("/v2/%s/blobs/uploads/%s", name, sid))

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// UpdateBlobUpload godoc
// @Summary Update image blob/layer upload
// @Description Update and finish an image's blob/layer upload given a digest
// @Accept  json
// @Produce json
// @Param   name     path    string     true        "repository name"
// @Param   session_id     path    string     true        "upload session_id"
// @Param 	digest	 query 	 string 		true				"blob/layer digest"
// @Success 201 {string} string	"created"
// @Header  202 {string} Location "/v2/{name}/blobs/uploads/{digest}"
// @Header  200 {object} api.DistContentDigestKey
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /v2/{name}/blobs/uploads/{session_id} [put].
func (s *Server) updateBlobUpload(w http.ResponseWriter, r *http.Request) {
	if s.DisableEntrypoint {
		s.proxy(w, r)
		return
	}

	defer r.Body.Close()

	vars := mux.Vars(r)
	name, ok := vars["name"]

	if !ok || name == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sid, ok := vars["session_id"]
	if !ok || sid == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	s.mu.Lock()
	sess := s.uploadSessions[sid]
	delete(s.uploadSessions, sid)
	s.mu.Unlock()

	if sess == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	digest := r.URL.Query().Get("digest")
	if digest == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// We weren't buffering it, just send it up.
	if sess.remote != "" {
		url := sess.remote
		if strings.Contains(url, "?") {
			url += "&digest=" + digest
		} else {
			url += "?digest=" + digest
		}

		s.Logger.Debug("uploading and closing blob", "target", url)

		out, err := s.outRequest(r.Method, url, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.Logger.Info("proxy", "url", out.URL.String())

		for k, v := range r.Header {
			out.Header[k] = v
		}

		resp, err := s.client.Do(out)
		if err != nil {
			s.Logger.Error("do error", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// resp.Write(os.Stdout)

		s.Logger.Info("proxy response", "status", resp.Status)

		defer resp.Body.Close()
	} else {
		s.Logger.Debug("closing up json blob", "digest", digest)

		// If patch wasn't used, we might be seeing the only data here
		if sess.buf == nil {
			sess.buf = &bytes.Buffer{}
		}

		io.Copy(sess.buf, r.Body)

		s.mu.Lock()
		if s.jsonBlobs == nil {
			s.jsonBlobs = map[string][]byte{}
		}

		// at this point, we're holding a json blob in memory and have not sent
		// it to the remote repo yet.
		s.jsonBlobs[digest] = sess.buf.Bytes()

		s.mu.Unlock()
	}

	w.Header().Set("Location", fmt.Sprintf("/v2/%s/blobs/%s", name, digest))
	w.WriteHeader(201)
}
