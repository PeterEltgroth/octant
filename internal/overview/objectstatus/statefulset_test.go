package objectstatus

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	cachefake "github.com/heptio/developer-dash/internal/cache/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func Test_statefulSet(t *testing.T) {
	cases := []struct {
		name     string
		init     func(*testing.T, *cachefake.MockCache) runtime.Object
		expected ObjectStatus
		isErr    bool
	}{
		{
			name: "in general",
			init: func(t *testing.T, c *cachefake.MockCache) runtime.Object {
				objectFile := "statefulset_ok.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusOK,
				Details:    component.TitleFromString("Stateful Set is OK"),
			},
		},
		{
			name: "not ready",
			init: func(t *testing.T, c *cachefake.MockCache) runtime.Object {
				objectFile := "statefulset_not_ready.yaml"
				return testutil.LoadObjectFromFile(t, objectFile)

			},
			expected: ObjectStatus{
				nodeStatus: component.NodeStatusWarning,
				Details:    component.TitleFromString("Stateful Set pods are not ready"),
			},
		},
		{
			name: "object is nil",
			init: func(t *testing.T, c *cachefake.MockCache) runtime.Object {
				return nil
			},
			isErr: true,
		},
		{
			name: "object is not a replication controller",
			init: func(t *testing.T, c *cachefake.MockCache) runtime.Object {
				return &unstructured.Unstructured{}
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			c := cachefake.NewMockCache(controller)

			object := tc.init(t, c)

			ctx := context.Background()
			status, err := statefulSet(ctx, object, c)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, status)
		})
	}
}
