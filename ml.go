package nuwa

import (
	"io/ioutil"

	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
)

var mler = mlImp{}

func ML() *mlImp {
	return &mler
}

type mlImp struct {
}

func (m *mlImp) KMeans(groupNum int) *KMeans {
	return &KMeans{
		km:       kmeans.New(),
		GroupNum: groupNum,
	}
}

type KMeans struct {
	TrainData clusters.Observations `json:"train_data"`
	GroupNum  int                   `json:"group_num"`
	km        kmeans.Kmeans         `json:"-"`
	Result    clusters.Clusters     `json:"result"`
}

func (k *KMeans) AddTrainDataOnec(data []float64) *KMeans {
	k.TrainData = append(k.TrainData, clusters.Coordinates(data))
	return k
}

func (k *KMeans) AddTrainData(data clusters.Observations) *KMeans {
	k.TrainData = append(k.TrainData, data...)
	return k
}

func (k *KMeans) Train() error {
	clusters, err := k.km.Partition(k.TrainData, 4)
	if err != nil {
		return err
	}
	k.Result = clusters
	return nil
}

func (k *KMeans) Predict(data []float64) int {
	return k.Result.Nearest(clusters.Coordinates(data))
}

func (k *KMeans) Save(filename string) error {
	return ioutil.WriteFile(filename, []byte(Helper().JsonEncode(k)), 0644)
}

func (k *KMeans) Load(filename string) *KMeans {
	Helper().JsonByFile(filename, k)
	return k
}
