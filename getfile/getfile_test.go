package getfile

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetFile(t *testing.T) {
	imageUrl1 := "https://static.wikia.nocookie.net/zelda_gamepedia_en/images/3/35/WW_Link_3.png/revision/latest?cb=20130913013026"
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	f := NewFileHandler(httpClient)
	size1, err := f.GetFile(imageUrl1)
	fmt.Printf("size of file: %d", size1)
	assert.Equal(t, nil, err, "get file error")
	assert.NotEqual(t, 0, size1, "get file: file's size")
	imageUrl2 := "https://sun9-14.userapi.com/impg/R_qVWzbHKY3YOyvGa3OMV_Um0Db__g8Xk34TWg/UUTRmqbhlOA.jpg?size=640x349&quality=96&proxy=1&sign=3e79dbc15855a03c4250b0a0fc0ddd12&type=album"
	size2, err := f.GetFile(imageUrl2)
	fmt.Printf("size of file: %d", size2)
	assert.Equal(t, nil, err, "get file error")
	assert.NotEqual(t, 0, size2, "get file: file's size")
}
