/*
Copyright 2021 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cr

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/tidwall/gjson"
)

//go:generate minimock -i Client -o cr_mock.go

type Client interface {
	Image(tag string) (v1.Image, error)
	Digest(tag string) (string, error)
}

type client struct {
	registryURL string
	authConfig  authn.AuthConfig
	options     *registryOptions
}

// NewClient creates container registry client using `repo` as prefix for tags passed to methods. If insecure flag is set to true, then no cert validation is performed.
// Repo example: "cr.example.com/ns/app"
func NewClient(repo, dockerCfgBase64 string, options ...Option) (Client, error) {
	opts := &registryOptions{}

	for _, opt := range options {
		opt(opts)
	}

	r := &client{
		registryURL: repo,
		options:     opts,
	}

	if !opts.withoutAuth {
		authConfig, err := readAuthConfig(dockerCfgBase64)
		if err != nil {
			return nil, err
		}
		r.authConfig = authConfig
	}

	return r, nil
}

func (r *client) Image(tag string) (v1.Image, error) {
	imageURL := r.registryURL + ":" + tag

	var nameOpts []name.Option
	if r.options.useHTTP {
		nameOpts = append(nameOpts, name.Insecure)
	}

	ref, err := name.ParseReference(imageURL, nameOpts...) // parse options available: weak validation, etc.
	if err != nil {
		return nil, err
	}

	imageOptions := make([]remote.Option, 0)
	if !r.options.withoutAuth {
		imageOptions = append(imageOptions, remote.WithAuth(authn.FromConfig(r.authConfig)))
	}
	if r.options.ca != "" {
		imageOptions = append(imageOptions, remote.WithTransport(GetHTTPTransport(r.options.ca)))
	}

	return remote.Image(
		ref,
		imageOptions...,
	)
}

func (r *client) Digest(tag string) (string, error) {
	image, err := r.Image(tag)
	if err != nil {
		return "", err
	}

	d, err := image.Digest()
	if err != nil {
		return "", err
	}

	return d.String(), nil
}

func readAuthConfig(dockerCfgBase64 string) (authn.AuthConfig, error) {
	dockerCfg, err := base64.StdEncoding.DecodeString(dockerCfgBase64)
	if err != nil {
		return authn.AuthConfig{}, err
	}
	auths := gjson.Get(string(dockerCfg), "auths").Map()
	authConfig := authn.AuthConfig{}

	// The config should have at least one .auths.* entry
	for _, a := range auths {
		err := json.Unmarshal([]byte(a.Raw), &authConfig)
		if err != nil {
			return authn.AuthConfig{}, err
		}
		return authConfig, nil
	}

	return authn.AuthConfig{}, fmt.Errorf("no auth data")
}

func GetHTTPTransport(ca string) (transport http.RoundTripper) {
	if ca == "" {
		return http.DefaultTransport
	}
	caPool, err := x509.SystemCertPool()
	if err != nil {
		panic(fmt.Errorf("cannot get system cert pool: %v", err))
	}

	caPool.AppendCertsFromPEM([]byte(ca))

	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{RootCAs: caPool},
		TLSNextProto:          make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
}

type registryOptions struct {
	ca          string
	useHTTP     bool
	withoutAuth bool
}

type Option func(options *registryOptions)

// WithCA use custom CA certificate
func WithCA(ca string) Option {
	return func(options *registryOptions) {
		options.ca = ca
	}
}

// WithInsecureSchema use http schema instead of https
func WithInsecureSchema(insecure bool) Option {
	return func(options *registryOptions) {
		options.useHTTP = insecure
	}
}

// WithDisabledAuth dont use authConfig
func WithDisabledAuth() Option {
	return func(options *registryOptions) {
		options.withoutAuth = true
	}
}
