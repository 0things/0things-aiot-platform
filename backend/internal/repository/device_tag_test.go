package repository

import (
	"context"
	"testing"

	"aiot-backend/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTagRepo(t *testing.T) *DeviceTagRepository {
	t.Helper()
	store := newRepositoryTestDB(t, &model.DeviceTag{})
	return NewDeviceTagRepository(store)
}

func findTag(tags []model.DeviceTag, key string) (model.DeviceTag, bool) {
	for _, tg := range tags {
		if tg.Key == key {
			return tg, true
		}
	}
	return model.DeviceTag{}, false
}

func TestDeviceTagRepository_ListEmpty(t *testing.T) {
	repo := newTagRepo(t)
	ctx := context.Background()

	tags, err := repo.ListTags(ctx, 11)
	require.NoError(t, err)
	assert.Empty(t, tags)
}

func TestDeviceTagRepository_SetIncremental(t *testing.T) {
	repo := newTagRepo(t)
	ctx := context.Background()

	require.NoError(t, repo.SetTags(ctx, 11, map[string]string{"a": "1", "b": "2"}, false))

	tags, err := repo.ListTags(ctx, 11)
	require.NoError(t, err)
	require.Len(t, tags, 2)
	a, ok := findTag(tags, "a")
	require.True(t, ok)
	assert.Equal(t, "manual", a.Source)

	// 增量添加：已存在的 key 更新 value，新 key 追加，旧 key 保留
	require.NoError(t, repo.SetTags(ctx, 11, map[string]string{"a": "new", "c": "3"}, false))

	tags, err = repo.ListTags(ctx, 11)
	require.NoError(t, err)
	require.Len(t, tags, 3)
	a, ok = findTag(tags, "a")
	require.True(t, ok)
	assert.Equal(t, "new", a.Value)
	assert.Equal(t, "manual", a.Source)
	_, ok = findTag(tags, "b")
	assert.True(t, ok)
	_, ok = findTag(tags, "c")
	assert.True(t, ok)
}

func TestDeviceTagRepository_SetReplace(t *testing.T) {
	repo := newTagRepo(t)
	ctx := context.Background()

	require.NoError(t, repo.SetTags(ctx, 11, map[string]string{"a": "1", "b": "2"}, false))
	// 全量替换：仅保留本次写入的 key
	require.NoError(t, repo.SetTags(ctx, 11, map[string]string{"x": "9"}, true))

	tags, err := repo.ListTags(ctx, 11)
	require.NoError(t, err)
	require.Len(t, tags, 1)
	_, ok := findTag(tags, "x")
	assert.True(t, ok)
	_, ok = findTag(tags, "a")
	assert.False(t, ok)
}

func TestDeviceTagRepository_Delete(t *testing.T) {
	repo := newTagRepo(t)
	ctx := context.Background()

	require.NoError(t, repo.SetTags(ctx, 11, map[string]string{"a": "1", "b": "2"}, false))

	require.NoError(t, repo.DeleteTags(ctx, 11, []string{"a"}))
	tags, err := repo.ListTags(ctx, 11)
	require.NoError(t, err)
	require.Len(t, tags, 1)
	_, ok := findTag(tags, "a")
	assert.False(t, ok)

	// 删除不存在的 key 不应报错
	require.NoError(t, repo.DeleteTags(ctx, 11, []string{"missing"}))
}
