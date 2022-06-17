package nuwa

import (
	"testing"
)

func Test_mlImp_KMeans(t *testing.T) {
	ML().KMeans(5).
		AddTrainDataOnec([]float64{1.2, 2.4}).
		AddTrainDataOnec([]float64{}).
		AddTrainDataOnec([]float64{}).
		AddTrainDataOnec([]float64{}).
		AddTrainDataOnec([]float64{}).
		AddTrainDataOnec([]float64{}).Train()
}
