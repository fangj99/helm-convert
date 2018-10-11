package transformers

import (
	"fmt"
	"reflect"
	"testing"

	"sigs.k8s.io/kustomize/pkg/gvk"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/resource"
	"sigs.k8s.io/kustomize/pkg/types"

	"github.com/davecgh/go-spew/spew"
)

type imageTagTransformerArgs struct {
	config    *types.Kustomization
	resources resmap.ResMap
}

func TestImageTagRun(t *testing.T) {
	var deploy = gvk.Gvk{Group: "apps", Version: "v1", Kind: "Deployment"}

	for _, test := range []struct {
		name     string
		input    *imageTagTransformerArgs
		expected *imageTagTransformerArgs
	}{
		{
			name: "it should retrieve images",
			input: &imageTagTransformerArgs{
				config: &types.Kustomization{},
				resources: resmap.ResMap{
					resource.NewResId(deploy, "deploy1"): resource.NewResourceFromMap(
						map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "Deployment",
							"metadata": map[string]interface{}{
								"name": "deploy1",
								"spec": map[string]interface{}{
									"template": map[string]interface{}{
										"spec": map[string]interface{}{
											"initContainers": []interface{}{
												map[string]interface{}{
													"name":  "busybox",
													"image": "busybox",
												},
											},
											"containers": []interface{}{
												map[string]interface{}{
													"name":  "nginx",
													"image": "nginx:1.7.9",
												},
												map[string]interface{}{
													"name":  "alpine",
													"image": "alpine@sha256:24a0c4b4a4c0eb97a1aabb8e29f18e917d05abfe1b7a7c07857230879ce7d3d3",
												},
											},
										},
									},
								},
							},
						}),
					resource.NewResId(deploy, "deploy2"): resource.NewResourceFromMap(
						map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "Deployment",
							"metadata": map[string]interface{}{
								"name": "deploy1",
								"spec": map[string]interface{}{
									"template": map[string]interface{}{
										"spec": map[string]interface{}{
											"containers": []interface{}{
												map[string]interface{}{
													"name":  "nginx",
													"image": "nginx:1.7.9",
												},
											},
										},
									},
								},
							},
						}),
				},
			},
			expected: &imageTagTransformerArgs{
				config: &types.Kustomization{
					ImageTags: []types.ImageTag{
						types.ImageTag{Name: "nginx", NewTag: "1.7.9"},
						types.ImageTag{Name: "alpine", Digest: "sha256:24a0c4b4a4c0eb97a1aabb8e29f18e917d05abfe1b7a7c07857230879ce7d3d3"},
						types.ImageTag{Name: "busybox"},
					},
				},
				resources: resmap.ResMap{
					resource.NewResId(deploy, "deploy1"): resource.NewResourceFromMap(
						map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "Deployment",
							"metadata": map[string]interface{}{
								"name": "deploy1",
								"spec": map[string]interface{}{
									"template": map[string]interface{}{
										"spec": map[string]interface{}{
											"initContainers": []interface{}{
												map[string]interface{}{
													"name":  "busybox",
													"image": "busybox",
												},
											},
											"containers": []interface{}{
												map[string]interface{}{
													"name":  "nginx",
													"image": "nginx:1.7.9",
												},
												map[string]interface{}{
													"name":  "alpine",
													"image": "alpine@sha256:24a0c4b4a4c0eb97a1aabb8e29f18e917d05abfe1b7a7c07857230879ce7d3d3",
												},
											},
										},
									},
								},
							},
						}),
					resource.NewResId(deploy, "deploy2"): resource.NewResourceFromMap(
						map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "Deployment",
							"metadata": map[string]interface{}{
								"name": "deploy1",
								"spec": map[string]interface{}{
									"template": map[string]interface{}{
										"spec": map[string]interface{}{
											"containers": []interface{}{
												map[string]interface{}{
													"name":  "nginx",
													"image": "nginx:1.7.9",
												},
											},
										},
									},
								},
							},
						}),
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("%s", test.name), func(t *testing.T) {
			lt := NewImageTagTransformer()
			err := lt.Transform(test.input.config, test.input.resources)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(test.input.config, test.expected.config) {
				t.Fatalf(
					"expected: \n %v\ngot:\n %v",
					spew.Sdump(test.expected.config.ImageTags),
					spew.Sdump(test.input.config.ImageTags),
				)
			}

			if !reflect.DeepEqual(test.input.resources, test.expected.resources) {
				t.Fatalf(
					"expected: \n %v\ngot:\n %v",
					spew.Sdump(test.expected.resources),
					spew.Sdump(test.input.resources),
				)
			}
		})
	}
}