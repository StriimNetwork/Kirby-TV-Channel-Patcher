package wad
import (
	"errors"
	"fmt"
	"github.com/wii-tools/wadlib"
	"io/ioutil"
	"os"
	"path/filepath"
)


// Pack is based off of the Pack function found here https://github.com/wii-tools/archon/blob/master/cmd/wad/pack.go
func Pack() error {
	// Ensure the inpath is sane.
	in := filepath.Join("patching-dir")
	stat, err := os.Stat(in)
	if err != nil {
		// This isn't an error we know how to cope with.
		return err
	} else if !stat.IsDir() {
		return errors.New(fmt.Sprintf("%s is not a directory", in))
	}

	err = os.Mkdir("WAD", 0777)
	if err != nil {
		if os.IsExist(err) {

		} else {
			return err
		}
	}

	// Ensure the outpath is sane.
	out := filepath.Join("WAD/Kirby-TV-Channel(Striim Network)")
	stat, err = os.Stat(out)
	if err != nil {
		if os.IsExist(err) || os.IsNotExist(err) {
			// We'll overwrite this file if necessary.
			// It's okay that it exists.
		} else {
			// This isn't an error we know how to cope with.
			return err
		}
	} else if stat.IsDir() {
		return errors.New(fmt.Sprintf("%s is not a file", in))
	}

	dir := directory{
		dir: in,
		// We don't know this yet.
		titleId: "",
	}

	// We'll create an empty WAD.
	wad := wadlib.WAD{}

	// Load the ticket first.
	ticket, err := dir.readFileWithSuffix("tik")
	if err != nil {
		return err
	}
	if err = wad.LoadTicket(ticket); err != nil {
		return err
	}

	// Now we know this!
	dir.titleId = fmt.Sprintf("%016x", wad.Ticket.TitleID)

	// Next, the TMD.
	tmd, err := dir.readSection("tmd")
	if err != nil {
		return err
	}
	if err = wad.LoadTMD(tmd); err != nil {
		return err
	}

	// Finally, certificates.
	certs, err := dir.readSection("certs")
	if err != nil {
		return err
	}
	wad.CertificateChain = certs

	// The meta and CRL sections may (or may not) exist.
	meta, err := dir.readSection("meta")
	if err != nil {
		// We don't mind this not existing.
		if !os.IsNotExist(err) {
			return err
		}
	} else if len(meta) != 0 {
		wad.Meta = meta
	}
	crl, err := dir.readSection("crl")
	if err != nil {
		// We don't mind this not existing.
		if !os.IsNotExist(err) {
			return err
		}
	} else if len(crl) != 0 {
		wad.CertificateRevocationList = crl
	}

	// Next up: data contents.
	// These should be exactly the same as what is listed within the TMD.
	wadfiles := make([]wadlib.WADFile, len(wad.TMD.Contents))
	for index, content := range wad.TMD.Contents {
		// Read the data file listed.
		data, err := dir.readFile(fmt.Sprintf("%08x.app", content.Index))
		if err != nil {
			return err
		}

		wadfiles[index] = wadlib.WADFile{
			ContentRecord: content,
			RawData:       data,
		}

		err = wadfiles[index].EncryptData(wad.Ticket.TitleKey)
		if err != nil {
			return err
		}
	}
	wad.Data = wadfiles

	// We're going to assume this value.
	// TODO: perhaps this type should be set elsewhere?
	// Perhaps a flag would do.
	wadContents, err := wad.GetWAD(wadlib.WADTypeCommon)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(out, wadContents, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

