package desultory

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"path"
)

func ReadZip(b *bytes.Buffer) (map[string][]byte, error) {
	fs := make(map[string][]byte)
	zr, err := zip.NewReader(bytes.NewReader(b.Bytes()), int64(b.Len()))
	if err != nil {
		return fs, err
	}
	for _, zf := range zr.File {
		f, err := zf.Open()
		if err != nil {
			return fs, err
		}
		defer f.Close()
		fb, err := ioutil.ReadAll(f)
		if err != nil {
			return fs, err
		}
		fs[zf.Name] = fb
	}
	return fs, nil
}

func ReadZipFromPath(file string) (map[string][]byte, error) {
	fs := make(map[string][]byte)
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return fs, err
	}
	return ReadZip(bytes.NewBuffer(b))
}

func ReadZipToPath(b *bytes.Buffer, directory string) error {
	fs, err := ReadZip(b)
	if err != nil {
		return err
	}
	WriteFilesToDirectory(fs, directory)
	return nil
}

func ReadZipFromPathToPath(file string, directory string) error {
	fs, err := ReadZipFromPath(file)
	if err != nil {
		return err
	}
	WriteFilesToDirectory(fs, directory)
	return nil
}

func WriteZip(files map[string][]byte) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	zw := zip.NewWriter(b)
	for fn, fs := range files {
		zf, err := zw.Create(fn)

		if err != nil {
			return b, err
		}
		_, err = zf.Write(fs)
		if err != nil {
			return b, err
		}
	}
	err := zw.Close()
	if err != nil {
		return b, err
	}
	return b, nil
}

func WriteZipToPath(files map[string][]byte, file string) error {
	b, err := WriteZip(files)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, b.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}

func WriteZipFromPath(directory string) (*bytes.Buffer, error) {
	fs, err := GetFilesFromDirectory(directory)
	if err != nil {
		return nil, err
	}
	b, err := WriteZip(fs)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func WriteZipFromPathByCommand(directory string) (*bytes.Buffer, error) {
	fn := "files.zip"
	fp := path.Join(directory, "/", fn)
	err := RunCommand("zip", []string{"-r", "-j", fp, directory}, nil, directory)
	if err != nil {
		return nil, err
	}
	fs, err := GetFilesFromDirectory(directory)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(fs[fn]), nil
}

func WriteZipFromPathToPath(directory string, file string) error {
	fs, err := GetFilesFromDirectory(directory)
	if err != nil {
		return err
	}
	return WriteZipToPath(fs, file)
}
