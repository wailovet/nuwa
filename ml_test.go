package nuwa

import (
	"testing"

	"golang.org/x/exp/errors/fmt"
)

func Test_mlImp_KMeans(t *testing.T) {
	k := ML().KMeans()
	k.AddTrainDataOnec([]float64{1.2, 2.4, 1.2, 2.4}).
		AddTrainDataOnec([]float64{7.2, 2.4, 1.2, 2.4}).
		AddTrainDataOnec([]float64{1.3, 2.4, 1.1, 2.4}).
		AddTrainDataOnec([]float64{4.2, 2.4, 3.2, 2.4}).
		AddTrainDataOnec([]float64{0.2, 2.4, 1.2, 7.4}).
		AddTrainDataOnec([]float64{3.2, 7.4, 7.2, 1.4}).Train(2)
	k.Save("123.json")

	fmt.Println(ML().KMeans().Load("123.json").Predict([]float64{3.4, 7.4, 7.2, 1.4}))

	fmt.Println(ML().KMeans().Load("123.json").Predict([]float64{7.2, 2.4, 1.2, 2.4}))
	fmt.Println(ML().KMeans().Load("123.json").Predict([]float64{0.2, 2.4, 1.2, 7.4}))
	fmt.Println(ML().KMeans().Load("123.json").Predict([]float64{1.3, 2.4, 1.1, 2.4}))

}
