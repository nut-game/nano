package cluster

import (
	"testing"

	"github.com/nut-game/nano/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInfoRetrieverRegion(t *testing.T) {
	t.Parallel()

	c := viper.New()
	c.Set("nano.cluster.info.region", "us")
	conf := config.NewConfig(c)

	infoRetriever := NewInfoRetriever(config.NewNanoConfig(conf).Cluster.Info)

	assert.Equal(t, "us", infoRetriever.Region())
}
