// Package multi contains a metadata.Provider that merges the results from
// several underlying metadata.Provider instances.
package multi

import "github.com/rothskeller/photo-tools/metadata"

// Provider provides the merged results from its several underlying providers.
type Provider []metadata.Provider

var _ metadata.Provider = Provider{}

// ProviderName is the name for the provider, for debug purposes.
func (p Provider) ProviderName() string { return "Multi-Provider" }
