package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

type Mix struct {
	Title, Details, UrlOriginal, UrlDownload string
}

func (m *Mix) HashSum() string {
	hash := md5.Sum([]byte(m.UrlOriginal))
	return hex.EncodeToString(hash[:])
}

func (m *Mix) DetailsUrl() string {
	return fmt.Sprintf("https://www.mixcloud.com/player/details/?key=%s", m.Details)
}

func (m *Mix) FileName() string {
	name := strings.Trim(m.Details, "/")
	return strings.Replace(name, "/", "-", -1)
}

func (m *Mix) MixPrint ()  {
	fmt.Println("Mix Title:", m.Title)
	fmt.Println("Mix Details:", m.Details)
	fmt.Println("Mix HashSum:", m.HashSum())
}
