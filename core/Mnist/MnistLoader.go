package Mnist

import (
	"github.com/goMLLibrary/core/Util"
	"github.com/petar/GoMNIST"
	"gonum.org/v1/gonum/mat"
	"image"
	"image/color"
	"os"
)

// LoadData : Mnistのデータセットを取得
func LoadData() (trainSet *MnistDataSet, testSet *MnistDataSet, err error) {
	// train : Mnistの学習用データ
	// test : Mnistのテスト用データ
	train, test, err := GoMNIST.Load(getDataSetPath())
	if err != nil {
		return nil, nil, err
	}

	trainDataSet := newMnistDataSet(train)
	testDataSet := newMnistDataSet(test)
	return trainDataSet, testDataSet, nil
}

type MnistDataSet struct {
	dataSet []MnistData
}

func newMnistDataSet(set *GoMNIST.Set) *MnistDataSet {
	dataSet := MnistDataSet{}
	dataSet.dataSet = make([]MnistData, 0, set.Count())
	for i, rawData := range set.Images {
		data := newMnistDataFromGoMNISTData(rawData, uint8(set.Labels[i]))
		dataSet.addData(data)
	}
	return &dataSet
}

// ExtractRandomDataSet : 指定したmnistのデータセットからランダムに指定サイズ分だけのデータを抽出する
func ExtractRandomDataSet(rawSet *MnistDataSet, count int) *MnistDataSet {
	dataSet := MnistDataSet{}
	dataSet.dataSet = make([]MnistData, 0, count)
	if rawSet.Count() < count {
		panic("count is not match!")
	}
	randomIndexs := Util.RandomIntArray(rawSet.Count(), count)
	for _, index := range randomIndexs {
		dataSet.dataSet = append(dataSet.dataSet, rawSet.GetData(index))
	}
	return &dataSet
}

func (set *MnistDataSet) addData(data MnistData) {
	set.dataSet = append(set.dataSet, data)
}

func (set *MnistDataSet) GetData(i int) MnistData {
	return set.dataSet[i]
}

func (set *MnistDataSet) GetDataSet() []MnistData {
	return set.dataSet
}

func (set *MnistDataSet) Count() int {
	return len(set.dataSet)
}

type MnistData struct {
	rawImage image.Image
	label    uint8
}

func newMnistDataFromGoMNISTData(src GoMNIST.RawImage, label uint8) MnistData {
	data := MnistData{src, label}
	return data
}

func (data *MnistData) GetImageVector() mat.Vector {
	if data.rawImage.ColorModel() != color.GrayModel {
		panic("mnist data is not gray model!")
	}
	if data.rawImage.Bounds().Min.X != 0 || data.rawImage.Bounds().Min.Y != 0 {
		panic("mnist data size is not match!")
	}

	r := data.rawImage.Bounds().Max.X * data.rawImage.Bounds().Max.Y

	vec := mat.NewVecDense(r, nil)
	for i := 0; i < data.rawImage.Bounds().Max.Y; r++ {
		for j := 0; j < data.rawImage.Bounds().Max.X; j++ {
			index := j + data.rawImage.Bounds().Max.X*i
			vec.SetVec(index, float64(data.rawImage.At(j, i).(color.Gray).Y))
		}
	}
	return vec
}

func (data *MnistData) GetLabelVector() mat.Vector {
	vec := mat.NewVecDense(10, nil)
	vec.SetVec(int(data.label), 1)
	return vec
}

func getDataSetPath() string {
	gopath := os.Getenv("GOPATH")
	return gopath + "/src/github.com/petar/GoMNIST/data"
}
