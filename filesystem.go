package desultory

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func DirectoryOrFileExists(directory string) (bool, error) {
	_, err := os.Stat(directory)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DeleteDirectoryContents(directory string) error {
	d, err := os.Open(directory)
	if err != nil {
		return err
	}
	defer d.Close()
	ns, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, n := range ns {
		fp := path.Join(directory, n)
		err = os.RemoveAll(fp)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteFile(file string, contents []byte, directory string) error {
	fp := path.Join(directory, file)
	fd := filepath.Dir(fp)
	de, err := DirectoryOrFileExists(fd)
	if err != nil {
		return err
	}
	if !de {
		err = os.MkdirAll(fd, 0777)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(fp, contents, 0666)
	if err != nil {
		return err
	}
	return nil
}

func WriteFilesToDirectory(files map[string][]byte, directory string) error {
	for f, c := range files {
		err := WriteFile(f, c, directory)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadFile(file string, directory string) ([]byte, error) {
	fp := path.Join(directory, file)
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GetFilesFromDirectory(directory string) (map[string][]byte, error) {
	fs := make(map[string][]byte)
	err := filepath.Walk(directory,
		func(f string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !fi.IsDir() {
				b, err := ioutil.ReadFile(f)
				if err != nil {
					return err
				}
				r, err := filepath.Rel(directory, f)
				if err != nil {
					return err
				}
				fs[r] = b
			}
			return nil
		})
	if err != nil {
		return fs, err
	}
	return fs, nil
}

func SerializeObject(o interface{}, file string, directory string) error {
	b, err := yaml.Marshal(o)
	if err != nil {
		return err
	}
	return WriteFile(file, b, directory)
}

func DeserializeObject(o interface{}, file string, directory string) error {
	b, err := ReadFile(file, directory)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, o)
	if err != nil {
		return err
	}
	return nil
}

func FindSubdirectory(name string, startDirectory string) (string, error) {
	p := ""
	err := filepath.Walk(startDirectory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && info.Name() == name {
				p = path
			}
			return nil
		})
	if err != nil {
		return "", err
	}
	a, err := filepath.Abs(p)
	if err != nil {
		return "", nil
	}
	return a, nil
}