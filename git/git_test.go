package git

import(
	"testing"
	"fmt"
	"os/exec"
	"path/filepath"
	"os"
	"strings"
	"errors"
)

func TestGitConle(T *testing.T){
	path,err:= getCurrentPath() 
	if err !=nil{
		fmt.Println(err)
	}
	Clone(path,"https://github.com/ssont/ich.example.git","master")
}

func getCurrentPath() (string, error) {
    file, err := exec.LookPath(os.Args[0])
    if err != nil {
        return "", err
    }
    path, err := filepath.Abs(file)
    if err != nil {
        return "", err
    }
    i := strings.LastIndex(path, "/")
    if i < 0 {
        i = strings.LastIndex(path, "\\")
    }
    if i < 0 {
        return "", errors.New(`error: "/" or "\".`)
    }
    return string(path[0 : i+1]), nil
}