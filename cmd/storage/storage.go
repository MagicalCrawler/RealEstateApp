package storage
import (
	// "context"
	// "errors"
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/MagicalCrawler/RealEstateApp/cmd/err"
)
type storage interface{
	Save(p *Page)error
	PickRangom(UserName string)(*Page,error)
	Remove(p *Page)error
	Isexists(p *Page)(bool,error)
}
type Page struct{
	URL string
	UserName string
	// Created time.Time
}
func (p Page)Hash()(string,error){
	h:=sha1.New()

	if _,er:=io.WriteString(h,p.URL);er!=nil{
		return  "", err.Wrap("can't calculate hash",er)
	}
	if _,er:=io.WriteString(h,p.UserName);er!=nil{
		return "", err.Wrap("can't calculate hash",er)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}