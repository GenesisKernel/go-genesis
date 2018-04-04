package api

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/GenesisKernel/go-genesis/packages/consts"

	"github.com/GenesisKernel/go-genesis/packages/converter"
	"github.com/GenesisKernel/go-genesis/packages/model"
	log "github.com/sirupsen/logrus"
)

func getAvatar(w http.ResponseWriter, r *http.Request, data *apiData, logger *log.Entry) error {
	parMember := data.params["member"].(string)
	memberID := converter.StrToInt64(parMember)
	member := &model.Member{}
	member.SetTablePrefix(converter.Int64ToStr(data.ecosystemId))
	found, err := member.Get(memberID)
	if err != nil {
		log.WithFields(log.Fields{"type": consts.DBError, "error": err}).
			Errorf("getting member with ecosystem: %d member_id: %d", data.ecosystemId, memberID)
		return errorAPI(w, "E_SERVER", http.StatusInternalServerError)
	}

	if !found {
		return errorAPI(w, "E_SERVER", http.StatusNotFound)
	}

	if member.ImageID == nil {
		return errorAPI(w, "E_SERVER", http.StatusNotFound)
	}

	bin := &model.Binary{}
	bin.SetTablePrefix(converter.Int64ToStr(data.ecosystemId))
	found, err = bin.GetByID(*member.ImageID)
	if err != nil {
		log.WithFields(log.Fields{"type": consts.DBError, "error": err}).Errorf("on getting binary by id %d", *member.ImageID)
		return errorAPI(w, "E_SERVER", http.StatusInternalServerError)
	}

	if !found {
		return errorAPI(w, "E_SERVER", http.StatusNotFound)
	}

	// cut the prefix like a 'data:blah-blah;base64,'
	b64data := bin.Data[bytes.IndexByte(bin.Data, ',')+1:]
	buf, err := base64.StdEncoding.DecodeString(string(b64data))
	if err != nil {
		log.WithFields(log.Fields{"type": consts.ConversionError, "error": err}).Error("on decoding avatar")
		return errorAPI(w, "E_SERVER", http.StatusInternalServerError)
	}

	mime := http.DetectContentType(buf)
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Length", strconv.Itoa(len(buf)))
	if _, err := w.Write(buf); err != nil {
		log.WithFields(log.Fields{"type": consts.IOError, "error": err}).Error("unable to write image")
		return err
	}

	return nil
}