package mapper

import (
	"fmt"
	"myBot/pkg/model"
	"strconv"
	"strings"
)

func MapGoodListToString(list []model.Good) string {
	fmt.Println(list)
	var res strings.Builder
	for i, good := range list {
		if good.Name == "" {
			continue
		}
		res.WriteString(strconv.Itoa(i+1) + ". " + good.Name + " " + strconv.Itoa(good.Price) + " рублей" + "\n")
	}
	return res.String()
}

func MapCartListToString() {

}
