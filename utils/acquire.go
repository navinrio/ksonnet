package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	jsonnet "github.com/strickyak/jsonnet_cgo"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/pkg/runtime"
)

// Read fetches and decodes K8s objects by path.
// TODO: Replace this with something supporting more sophisticated
// content negotiation.
func Read(vm *jsonnet.VM, path string) ([]runtime.Object, error) {
	ext := filepath.Ext(path)
	if ext == ".json" {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return jsonReader(f)
	} else if ext == ".yaml" {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return yamlReader(f)
	} else if ext == ".jsonnet" {
		return jsonnetReader(vm, path)
	}

	return nil, fmt.Errorf("Unknown file extension: %s", path)
}

func jsonReader(r io.Reader) ([]runtime.Object, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	obj, _, err := runtime.UnstructuredJSONScheme.Decode(data, nil, nil)
	if err != nil {
		return nil, err
	}
	return []runtime.Object{obj}, nil
}

func yamlReader(r io.ReadCloser) ([]runtime.Object, error) {
	decoder := yaml.NewDocumentDecoder(r)
	ret := []runtime.Object{}
	buf := []byte{}
	for {
		_, err := decoder.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		jsondata, err := yaml.ToJSON(buf)
		if err != nil {
			return nil, err
		}
		obj, _, err := runtime.UnstructuredJSONScheme.Decode(jsondata, nil, nil)
		if err != nil {
			return nil, err
		}
		ret = append(ret, obj)
	}
	return ret, nil
}

func jsonnetReader(vm *jsonnet.VM, path string) ([]runtime.Object, error) {
	jsonstr, err := vm.EvaluateFile(path)
	if err != nil {
		return nil, err
	}

	glog.V(4).Infof("jsonnet result is: %s\n", jsonstr)

	return jsonReader(strings.NewReader(jsonstr))
}

// FlattenToV1 expands any List-type objects into their members, and
// cooerces everything to v1.Unstructured.  Panics if coercion
// encounters an unexpected object type.
func FlattenToV1(objs []runtime.Object) []*runtime.Unstructured {
	ret := make([]*runtime.Unstructured, 0, len(objs))
	for _, obj := range objs {
		switch o := obj.(type) {
		case *runtime.UnstructuredList:
			for _, item := range o.Items {
				ret = append(ret, item)
			}
		case *runtime.Unstructured:
			ret = append(ret, o)
		default:
			panic("Unexpected unstructured object type")
		}
	}
	return ret
}