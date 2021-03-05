package generate

import (
	"log"
	"os"

	"github.com/w-haibara/tempura/config"
	"github.com/w-haibara/tempura/parts"
)

func Generate(c config.Config) error {
	f := open(c.OutFile)
	f.writeStr(parts.Header1())
	f.writeStr(parts.Style(c.Style))
	f.writeStr(parts.Header2())
	f.writeStr(parts.Commands(c.Commands))
	f.writeStr(parts.Footer(c.Greeting))

	return nil
}

type file struct {
	*os.File
}

func open(path string) file {
	if err := os.WriteFile(path, []byte{}, 0666); err != nil {
		log.Panic(err)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}

	return file{f}
}

func (f file) writeStr(s string) {
	f.writeBytes([]byte(s))
}

func (f file) writeBytes(b []byte) {
	if _, err := f.Write(b); err != nil {
		f.close()
		log.Panic(err)
	}
}

func (f file) close() {
	if err := f.Close(); err != nil {
		log.Panic(err)
	}
}
