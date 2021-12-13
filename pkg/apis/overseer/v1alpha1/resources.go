/*
Copyright 2021.

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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

// Git refers to the Git repository.
type Git struct {
	// Repository URL to clone from.
	Url string `json:"url,omitempty"`

	// Revision to checkout. (branch, tag, sha, ref, etc...)
	// +optional
	Revision *string `json:"revision,omitempty"`

	// Subpath A subpath within checked out source
	// where the source to build is located.
	// +optional
	Subpath *string `json:"subpath,omitempty"`

	// Credentials
	// +optional
	Credentials *corev1.LocalObjectReference `json:"credentials,omitempty"`
}

// Image refers to the docker registry.
type Image struct {
	// Docker image name.
	// +require
	Image string `json:"image"`

	// Credentials
	// +optional
	Credentials *corev1.LocalObjectReference `json:"credentials,omitempty"`
}

// Param is a key/value that populates a strategy parameter
type Param struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
