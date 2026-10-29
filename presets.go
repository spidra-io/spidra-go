package spidra

import "context"

// PresetsResource handles saved tool and template presets.
type PresetsResource struct {
	http *httpClient
}

// List returns saved presets with optional pagination parameters.
func (r *PresetsResource) List(ctx context.Context, queryParams map[string]string) (*ListPresetsResponse, error) {
	var resp ListPresetsResponse
	if err := r.http.get(ctx, "/presets", queryParams, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a saved preset by its ID.
func (r *PresetsResource) Get(ctx context.Context, presetID string) (*Preset, error) {
	var preset Preset
	if err := r.http.get(ctx, "/presets/"+presetID, nil, &preset); err != nil {
		return nil, err
	}
	return &preset, nil
}

// Create saves a new reusable preset configuration.
func (r *PresetsResource) Create(ctx context.Context, params PresetParams) (*Preset, error) {
	var preset Preset
	if err := r.http.post(ctx, "/presets", params, &preset); err != nil {
		return nil, err
	}
	return &preset, nil
}

// Update modifies an existing preset configuration.
func (r *PresetsResource) Update(ctx context.Context, presetID string, updates map[string]any) (*Preset, error) {
	var preset Preset
	if err := r.http.patch(ctx, "/presets/"+presetID, updates, &preset); err != nil {
		return nil, err
	}
	return &preset, nil
}

// Delete permanently removes a preset by its ID.
func (r *PresetsResource) Delete(ctx context.Context, presetID string) error {
	return r.http.delete(ctx, "/presets/"+presetID, nil)
}

// Run executes a saved preset with optional overrides and wait flag.
func (r *PresetsResource) Run(ctx context.Context, presetID string, overrides map[string]any, wait bool) (*RunPresetResponse, error) {
	body := map[string]any{
		"wait": wait,
	}
	if overrides != nil {
		body["overrides"] = overrides
	}
	var resp RunPresetResponse
	if err := r.http.post(ctx, "/presets/"+presetID+"/run", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

