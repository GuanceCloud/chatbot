package difydataset

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAgent(t *testing.T) {
	apiserver, apikey := GetAPI()
	d := NewDataset(apiserver, apikey)
	v, err := d.ListDatasets(1, 20)
	if err != nil {
		t.Error(err)
	}
	t.Error(v)
}

func TestLoad(t *testing.T) {
	baseDir := "<path-to-repo>/dataflux-doc"

	v, err := os.ReadFile("./docs.txt")
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(string(v), "\n")
	docList := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "docs/zh/") && strings.HasSuffix(line, ".md") {
			docList = append(docList, line)
		}
	}

	var pipelineDoc []string
	var dqlDoc []string
	for _, line := range docList {
		if strings.Contains(line, "dql") {
			dqlDoc = append(dqlDoc, line)
		}
		if strings.Contains(line, "pipeline") {
			pipelineDoc = append(pipelineDoc, line)
		}
	}

	log, err := os.OpenFile("log.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	uploadDoc(baseDir, pipelineDoc, "fcbe0e92-66dd-4f50-8673-023506affa3d", log)
	uploadDoc(baseDir, dqlDoc, "a03f18b9-766e-4ec3-8513-96ccf94e5a7a", log)
}

func uploadDoc(baseDir string, docList []string, datasetID string, log *os.File) {
	apiserver, apikey := GetAPI()
	d := NewDataset(apiserver, apikey)
	for i := 0; i < len(docList); i++ {
		fp := filepath.Join(baseDir, docList[i])
		f, err := os.ReadFile(fp)
		if err != nil {
			log.WriteString(fmt.Sprintf("===file %s, err %s\n", fp, err.Error()))
			continue
		}
		v, err := d.CreateDocByText(datasetID, "economy", fp, string(f))
		if err != nil {
			log.WriteString(fmt.Sprintf("===file %s, err %s\n", fp, err.Error()))
			continue
		} else {
			if val, err := json.Marshal(v); err != nil {
				log.WriteString(fmt.Sprintf("===file %s, err %s\n", fp, err.Error()))
			} else {
				log.WriteString(fmt.Sprintf("===file %s, ret %s\n", fp, string(val)))
			}
		}
	}
}
